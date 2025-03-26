package chat

import (
	"context"
	"gf-chat/api"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"sync"
	"unicode/utf8"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

type iWsConn interface {
	readMsg()
	sendMsg()
	close()
	run()
	deliver(action *api.ChatAction)
	getUserId() uint
	getUser() iChatUser
	getUuid() string
	getPlatform() string
	getCustomerId() uint
	getLastActive() *gtime.Time
	createTime() *gtime.Time
}

// client represents a chat client connection
type client struct {
	conn        *websocket.Conn
	closeSignal chan interface{}     // 连接断开后的广播通道，用于中断readMsg,sendMsg goroutine
	send        chan *api.ChatAction // 发送的消息chan
	sync.Once
	manager    connManager
	user       iChatUser
	uuid       string
	created    *gtime.Time
	limiter    *rate.Limiter
	lastActive *gtime.Time
	platform   string
}

// newClient creates a new client instance
func newClient(conn *websocket.Conn, user iChatUser, platform string) *client {
	return &client{
		conn:        conn,
		closeSignal: make(chan interface{}),
		send:        make(chan *api.ChatAction, 100),
		Once:        sync.Once{},
		user:        user,
		uuid:        guid.S(),
		created:     gtime.Now(),
		limiter:     rate.NewLimiter(5, 10),
		platform:    platform,
		lastActive:  gtime.Now(),
	}
}

// GetLastActive returns the last active time of the client
func (c *client) getLastActive() *gtime.Time {
	return c.lastActive
}

// GetCreatedTime returns the creation time of the client
func (c *client) createTime() *gtime.Time {
	return c.created
}

// GetCustomerId returns the customer ID of the client
func (c *client) getCustomerId() uint {
	return c.user.getCustomerId()
}

// GetUuid 每个连接的unique id
func (c *client) getUuid() string {
	return c.uuid
}

// GetPlatform returns the platform of the client
func (c *client) getPlatform() string {
	return c.platform
}

// GetUser returns the user of the client
func (c *client) getUser() iChatUser {
	return c.user
}

// GetUserId returns the user ID of the client
func (c *client) getUserId() uint {
	return c.user.getPrimaryKey()
}

// runs the client's read and send message goroutines
func (c *client) run() {
	go c.readMsg()
	go c.sendMsg()
}

// closes the connection and unregister the client from the manager
func (c *client) close() {
	c.Once.Do(func() {
		close(c.closeSignal)
		_ = c.conn.Close()
		c.manager.unregister(c)
	})
}

// validates a chat message before sending it:
// - Checks rate limiting to prevent spam
// - Validates content length (non-empty and under 512 chars)
// - Validates request ID exists and is valid length (1-20 chars)
// - Validates message type is allowed
// Returns error if validation fails
func (c *client) validate(data map[string]interface{}) error {
	if !c.limiter.Allow() {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "发送过于频繁，请慢一些")
	}
	content, exist := data["content"]
	if exist {
		s, ok := content.(string)
		if !ok {
			return gerror.NewCode(gcode.CodeBusinessValidationFailed, "消息不合法")
		}
		length := utf8.RuneCountInString(s)
		if length == 0 {
			return gerror.NewCode(gcode.CodeBusinessValidationFailed, "请勿发送空内容")
		}
		if length > 512 {
			return gerror.NewCode(gcode.CodeBusinessValidationFailed, "内容长度必须小于512个字符")
		}
	}
	reqId, exist := data["req_id"]
	if !exist {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "消息不合法")
	}
	reqIdStr, ok := reqId.(string)
	if !ok {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "消息不合法")
	}
	length := len(reqIdStr)
	if length <= 0 || length > 20 {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "消息不合法")
	}
	types, exist := data["type"]
	if !exist {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "消息不合法")
	}
	typeStr, ok := types.(string)
	if !ok {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "消息不合法")
	}
	if !c.isTypeValid(typeStr) {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "消息不合法")
	}
	return nil
}

// isTypeValid checks if the message type is valid
func (c *client) isTypeValid(t string) bool {
	allowTypes := []string{
		consts.MessageTypeText,
		consts.MessageTypeImage,
		consts.MessageTypeAudio,
		consts.MessageTypeVideo,
		consts.MessageTypePdf,
		consts.MessageTypeNavigate,
	}
	return slice.Contain(allowTypes, t)
}

func (c *client) readMsg() {
	msgChan := make(chan []byte, 50)
	ctx := gctx.New()

	// Start message reader goroutine
	go c.readMessageLoop(msgChan)

	for {
		select {
		case <-c.closeSignal:
			return
		case msgBytes := <-msgChan:
			if err := c.handleMessage(ctx, msgBytes); err != nil {
				log.Errorf(ctx, "%+v", err)
			}
		}
	}
}

// readMessageLoop continuously reads messages from websocket connection
func (c *client) readMessageLoop(msgChan chan<- []byte) {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			c.close()
			return
		}
		select {
		case <-c.closeSignal:
			return
		case msgChan <- message:
		}
	}
}

// handleMessage processes a single message
func (c *client) handleMessage(ctx context.Context, msgBytes []byte) error {
	log.Debug(ctx, "read", msgBytes)
	act, err := action.unMarshal(msgBytes)
	if err != nil {
		return err
	}
	data, ok := act.Data.(g.Map)
	if !ok {
		return gerror.New("invalid action data type")
	}
	if err = c.validate(data); err != nil {
		c.deliver(action.newErrorMessage(err.Error()))
		return nil
	}
	if act.Action != consts.ActionSendMessage {
		return gerror.New("invalid action type")
	}
	return c.processMessage(ctx, act)
}

// processMessage handles message type actions
func (c *client) processMessage(ctx context.Context, act *api.ChatAction) error {
	msg, err := action.getMessage(act)
	if err != nil {
		return err
	}
	// Set message metadata
	msg.CustomerId = c.getCustomerId()
	msg.ReceivedAt = gtime.Now()

	// Set user information based on type
	switch u := c.getUser().(type) {
	case *admin:
		msg.Admin = u.Entity
	case *user:
		msg.User = u.Entity
	}

	// Process message asynchronously
	go c.manager.handleMessage(ctx, c, msg)
	c.lastActive = gtime.Now()
	return nil
}

// Deliver 投递消息
func (c *client) deliver(act *api.ChatAction) {
	c.send <- act
}

// sendMsg handles sending messages to the websocket connection
func (c *client) sendMsg() {
	ctx := gctx.New()
	for {
		select {
		case act := <-c.send:
			if err := c.processOutgoingMessage(ctx, act); err != nil {
				log.Errorf(ctx, "%+v", err)
			}
		case <-c.closeSignal:
			return
		}
	}
}

// processOutgoingMessage handles the processing and sending of a single message
func (c *client) processOutgoingMessage(ctx context.Context, act *api.ChatAction) error {
	// Marshal message
	msgBytes, err := action.marshal(ctx, *act)
	if err != nil {
		return err
	}
	log.Debug(ctx, "send", string(msgBytes))

	// Send message
	if err = c.conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		c.close()
		return gerror.New("connection closed")
	}

	// Handle special actions
	if err = c.handleSpecialActions(ctx, act); err != nil {
		return err
	}

	return nil
}

// handleSpecialActions processes special action types that require additional handling
func (c *client) handleSpecialActions(ctx context.Context, act *api.ChatAction) error {
	switch act.Action {
	case consts.ActionMoreThanOne, consts.ActionOtherLogin:
		c.close()
		return nil

	case consts.ActionReceiveMessage:
		return c.handleReceiveMessage(ctx, act)
	}

	return nil
}

// handleReceiveMessage processes received messages and updates their send timestamp
func (c *client) handleReceiveMessage(ctx context.Context, act *api.ChatAction) error {
	msg, ok := act.Data.(*model.CustomerChatMessage)
	if !ok {
		return gerror.NewCode(gcode.CodeValidationFailed, "action.data is not a message model")
	}
	if msg.SendAt != nil {
		return nil // Message already marked as sent
	}
	// Update message send timestamp
	_, err := service.ChatMessage().UpdatePri(ctx, msg.Id, do.CustomerChatMessages{
		SendAt: gtime.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}
