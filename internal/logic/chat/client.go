package chat

import (
	"gf-chat/api"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
	"sync"
	"unicode/utf8"

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
func (c *client) getLastActive() *gtime.Time {
	return c.lastActive
}
func (c *client) createTime() *gtime.Time {
	return c.created
}

func (c *client) getCustomerId() uint {
	return c.user.getCustomerId()
}

// GetUuid 每个连接的unique id
func (c *client) getUuid() string {
	return c.uuid
}

func (c *client) getPlatform() string {
	return c.platform
}

func (c *client) getUser() iChatUser {
	return c.user
}

func (c *client) getUserId() uint {
	return c.user.getPrimaryKey()
}

func (c *client) run() {
	go c.readMsg()
	go c.sendMsg()
}

// Close 幂等close方法 关闭连接，相关清理
func (c *client) close() {
	c.Once.Do(func() {
		close(c.closeSignal)
		_ = c.conn.Close()
		c.manager.unregister(c)
	})
}

// 发送消息验证
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

// ReadMsg 从websocket读消息
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
						switch iu.(type) {
						case *admin:
							u := iu.(*admin)
							msg.Admin = u.Entity
						case *user:
							u := iu.(*user)
							msg.User = u.Entity
						}
						msg.CustomerId = c.getCustomerId()
						msg.ReceivedAt = gtime.Now()
						c.manager.receiveMessage(&chatConnMessage{
							Msg:  msg,
							Conn: c,
						})
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

// SendMsg 发消息
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
