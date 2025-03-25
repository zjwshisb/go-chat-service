package chat

import (
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
		if ok {
			length := utf8.RuneCountInString(s)
			if length == 0 {
				return gerror.NewCode(gcode.CodeBusinessValidationFailed, "请勿发送空内容")
			}
			if length > 512 {
				return gerror.NewCode(gcode.CodeBusinessValidationFailed, "内容长度必须小于512个字符")
			}
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

// readMsg continuously reads messages from the websocket connection.
// It runs in a goroutine and handles incoming messages by:
// 1. Reading raw messages from the websocket connection
// 2. Unmarshaling them into ChatAction objects
// 3. Validating the message data
// 4. Processing messages based on their action type (e.g. sending messages)
// 5. Updating the last active timestamp
// The method will exit when the connection is closed via the closeSignal channel.
func (c *client) readMsg() {
	var msg = make(chan []byte, 50)
	for {
		go func() {
			_, message, err := c.conn.ReadMessage()
			// 读消息失败说明连接异常，调用close方法
			if err != nil {
				c.close()
			} else {
				msg <- message
			}
		}()
		select {
		case <-c.closeSignal:
			return
		case msgStr := <-msg:
			ctx := gctx.New()
			act, err := action.unMarshal(msgStr)
			log.Debug(ctx, "read", msgStr)
			if err != nil {
				log.Errorf(ctx, "%+v", err)
				break
			}
			data, ok := act.Data.(g.Map)
			if !ok {
				break
			}
			err = c.validate(data)
			if err != nil {
				c.deliver(action.newErrorMessage(err.Error()))
			} else {
				switch act.Action {
				case consts.ActionSendMessage:
					msg, err := action.getMessage(act)
					if err != nil {
						log.Errorf(ctx, "%+v", err)
					} else {
						iu := c.getUser()
						switch u := iu.(type) {
						case *admin:
							msg.Admin = u.Entity
						case *user:
							msg.User = u.Entity
						}
						msg.CustomerId = c.getCustomerId()
						msg.ReceivedAt = gtime.Now()
						go func() {
							c.manager.handleMessage(ctx, c, msg)
						}()
						c.lastActive = gtime.Now()
					}
				}
			}

		}
	}
}

// Deliver 投递消息
func (c *client) deliver(act *api.ChatAction) {
	c.send <- act
}

// sendMsg handles sending messages to the websocket connection. It runs in a separate goroutine and continuously processes messages until closed.
// It handles:
// - Receiving ChatAction messages from the client's send channel
// - Marshaling messages to JSON format before sending
// - Sending messages over the websocket connection using WriteMessage
// - Special action handling:
//   - ActionMoreThanOne: Closes connection when user has multiple active sessions
//   - ActionOtherLogin: Closes connection when user logs in elsewhere
//   - ActionReceiveMessage: Updates message send timestamp in database after successful delivery
//
// - Graceful shutdown via closeSignal channel
// - Error handling for marshal/send failures
// The method blocks until either:
// 1. An unrecoverable error occurs during message sending
// 2. The connection is explicitly closed
// 3. The closeSignal channel receives a shutdown signal
func (c *client) sendMsg() {
	for {
		select {
		case act := <-c.send:
			ctx := gctx.New()
			msgByte, err := action.marshal(ctx, *act)
			log.Debug(ctx, "send", string(msgByte))
			if err != nil {
				log.Errorf(ctx, "%+v", err)
				break
			}
			err = c.conn.WriteMessage(websocket.TextMessage, msgByte)
			if err != nil {
				log.Errorf(ctx, "%+v", err)
				c.close()
				return
			}
			switch act.Action {
			case consts.ActionMoreThanOne:
				c.close()
			case consts.ActionOtherLogin:
				c.close()
			case consts.ActionReceiveMessage:
				msg, ok := act.Data.(*model.CustomerChatMessage)
				if !ok {
					err = gerror.NewCode(gcode.CodeValidationFailed, "action.data is not a message model")
					log.Errorf(ctx, "%+v", err)
				} else {
					if msg.SendAt == nil {
						_, err := service.ChatMessage().UpdatePri(ctx, msg.Id, do.CustomerChatMessages{
							SendAt: gtime.Now(),
						})
						if err != nil {
							log.Errorf(ctx, "%+v", err)
						}
					}
				}
			default:
			}
		case <-c.closeSignal:
			return
		}
	}
}
