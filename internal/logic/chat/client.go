package chat

import (
	"gf-chat/api/v1"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/guid"
	"sync"
	"unicode/utf8"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

type iWsConn interface {
	ReadMsg()
	SendMsg()
	Close()
	Run()
	Deliver(action *v1.ChatAction)
	GetUserId() uint
	GetUser() IChatUser
	GetUuid() string
	GetPlatform() string
	GetCustomerId() uint
	CreateTime() *gtime.Time
}

type client struct {
	conn        *websocket.Conn
	closeSignal chan interface{}    // 连接断开后的广播通道，用于中断readMsg,sendMsg goroutine
	send        chan *v1.ChatAction // 发送的消息chan
	sync.Once
	manager  connManager
	user     IChatUser
	uuid     string
	created  *gtime.Time
	limiter  *rate.Limiter
	platform string
}

func newClient(conn *websocket.Conn, user IChatUser, platform string) *client {
	return &client{
		conn:        conn,
		closeSignal: make(chan interface{}),
		send:        make(chan *v1.ChatAction, 100),
		Once:        sync.Once{},
		user:        user,
		uuid:        guid.S(),
		created:     gtime.Now(),
		limiter:     rate.NewLimiter(5, 10),
		platform:    platform,
	}
}

func (c *client) CreateTime() *gtime.Time {
	return c.created
}

func (c *client) GetCustomerId() uint {
	return c.user.GetCustomerId()
}

// GetUuid 每个连接的unique id
func (c *client) GetUuid() string {
	return c.uuid
}

func (c *client) GetPlatform() string {
	return c.platform
}

func (c *client) GetUser() IChatUser {
	return c.user
}

func (c *client) GetUserId() uint {
	return c.user.GetPrimaryKey()
}

func (c *client) Run() {
	go c.ReadMsg()
	go c.SendMsg()
}

// Close 幂等close方法 关闭连接，相关清理
func (c *client) Close() {
	c.Once.Do(func() {
		close(c.closeSignal)
		_ = c.conn.Close()
		c.manager.Unregister(c)
	})
}

// 发送消息验证
func (c *client) validate(data map[string]interface{}) error {
	if !c.limiter.Allow() {
		return gerror.New("发送过于频繁，请慢一些")
	}
	content, exist := data["content"]
	if exist {
		s, ok := content.(string)
		if ok {
			length := utf8.RuneCountInString(s)
			if length == 0 {
				return gerror.New("请勿发送空内容")
			}
			if length > 512 {
				return gerror.New("内容长度必须小于512个字符")
			}
		}
	}
	reqId, exist := data["req_id"]
	if !exist {
		return gerror.New("消息不合法")
	}
	reqIdStr, ok := reqId.(string)
	if !ok {
		return gerror.New("消息不合法")

	}
	length := len(reqIdStr)
	if length <= 0 || length > 20 {
		return gerror.New("消息不合法")
	}
	types, exist := data["type"]
	if !exist {
		return gerror.New("消息不合法")
	}
	typeStr, ok := types.(string)
	if !ok {
		return gerror.New("消息不合法")
	}
	if !c.isTypeValid(typeStr) {
		return gerror.New("消息不合法")
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
func (c *client) ReadMsg() {
	var msg = make(chan []byte, 50)
	for {
		go func() {
			_, message, err := c.conn.ReadMessage()
			// 读消息失败说明连接异常，调用close方法
			if err != nil {
				c.Close()
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
			g.Log("ws").Debug(ctx, msgStr)
			if err != nil {
				g.Log().Errorf(ctx, "%+v", err)
				break
			}
			data, ok := act.Data.(map[string]interface{})
			if !ok {
				break
			}
			err = c.validate(data)
			if err != nil {
				c.Deliver(action.newErrorMessage(err.Error()))
			} else {
				switch act.Action {
				case consts.ActionSendMessage:
					msg, err := action.getMessage(act)
					if err != nil {
						g.Log().Errorf(ctx, "%+v", err)
						break
					}
					iu := c.GetUser()
					switch iu.(type) {
					case *admin:
						u := iu.(*admin)
						msg.Admin = u.Entity
					case *user:
						u := iu.(*user)
						msg.User = u.Entity
					}
					msg.CustomerId = c.GetCustomerId()
					msg.ReceivedAt = gtime.Now()
					c.manager.receiveMessage(&chatConnMessage{
						Msg:  msg,
						Conn: c,
					})
				}
			}

		}
	}
}

// Deliver 投递消息
func (c *client) Deliver(act *v1.ChatAction) {
	c.send <- act
}

// SendMsg 发消息
func (c *client) SendMsg() {
	for {
		select {
		case act := <-c.send:
			ctx := gctx.New()
			msgByte, err := action.marshal(ctx, *act)
			if err != nil {
				g.Log().Errorf(ctx, "%+v", err)
				break
			}
			err = c.conn.WriteMessage(websocket.TextMessage, msgByte)
			if err != nil {
				g.Log().Errorf(ctx, "%+v", err)
				c.Close()
				return
			}
			switch act.Action {
			case consts.ActionMoreThanOne:
				c.Close()
			case consts.ActionOtherLogin:
				c.Close()
			case consts.ActionReceiveMessage:
				msg, ok := act.Data.(*model.CustomerChatMessage)
				if !ok {
					break
				}
				if msg.SendAt == nil {
					_, err := service.ChatMessage().UpdatePri(ctx, msg.Id, do.CustomerChatMessages{
						SendAt: gtime.Now(),
					})
					if err != nil {
						g.Log().Errorf(ctx, "%+v", err)
					}
				}

			default:
			}
		case <-c.closeSignal:
			return
		}
	}
}
