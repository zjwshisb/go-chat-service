package chat

import (
	"errors"
	"gf-chat/internal/consts"
	"gf-chat/internal/contract"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/chat"
	"gf-chat/internal/model/relation"
	"gf-chat/internal/service"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

type client struct {
	Conn        *ghttp.WebSocket
	CloseSignal chan interface{}  // 连接断开后的广播通道，用于中断readMsg,sendMsg goroutine
	Send        chan *chat.Action // 发送的消息chan
	sync.Once
	Manager  connManager
	User     contract.IChatUser
	Uuid     string
	Created  int64
	Limiter  *rate.Limiter
	Platform string
}

func (c *client) GetCreateTime() int64 {
	return c.Created
}

func (c *client) GetCustomerId() uint {
	return c.User.GetCustomerId()
}

// GetUuid 每个连接的unique id
func (c *client) GetUuid() string {
	return c.Uuid
}

func (c *client) GetPlatform() string {
	return c.Platform
}

func (c *client) GetUser() contract.IChatUser {
	return c.User
}

func (c *client) GetUserId() uint {
	return c.User.GetPrimaryKey()
}

func (c *client) Run() {
	go c.ReadMsg()
	go c.SendMsg()
}

// Close 幂等close方法 关闭连接，相关清理
func (c *client) Close() {
	c.Once.Do(func() {
		close(c.CloseSignal)
		_ = c.Conn.Close()
		c.Manager.Unregister(c)
	})
}

// 发送消息验证
func (c *client) validate(data map[string]interface{}) error {
	if !c.Limiter.Allow() {
		return errors.New("发送过于频繁，请慢一些")
	}
	content, exist := data["content"]
	if exist {
		s, ok := content.(string)
		if ok {
			length := utf8.RuneCountInString(s)
			if length == 0 {
				return errors.New("请勿发送空内容")
			}
			if length > 512 {
				return errors.New("内容长度必须小于512个字符")
			}
		}
	}
	reqId, exist := data["req_id"]
	if !exist {
		return errors.New("消息不合法")
	} else {
		idStr, ok := reqId.(string)
		if ok {
			length := len(idStr)
			if length <= 0 && length > 20 {
				return errors.New("消息不合法")
			}
		} else {
			return errors.New("消息不合法")
		}
	}
	types, exist := data["type"]
	if !exist {
		return errors.New("消息不合法")
	} else {
		typeStr, ok := types.(string)
		if ok {
			if !c.isTypeValid(typeStr) {
				return errors.New("消息不合法")
			}
		} else {
			return errors.New("消息不合法")
		}
	}
	return nil
}

func (c *client) isTypeValid(t string) bool {
	if t != consts.MessageTypeText &&
		t != consts.MessageTypeImage &&
		t != consts.MessageTypeNavigate &&
		t != consts.MessageTypeRate {
		return false
	}
	return true
}

// ReadMsg 从websocket读消息
func (c *client) ReadMsg() {
	var msg = make(chan []byte, 50)
	for {
		go func() {
			_, message, err := c.Conn.ReadMessage()
			// 读消息失败说明连接异常，调用close方法
			if err != nil {
				c.Close()
			} else {
				msg <- message
			}
		}()
		select {
		case <-c.CloseSignal:
			return
		case msgStr := <-msg:
			act, err := service.Action().UnMarshalAction(msgStr)
			if err != nil {
				g.Log().Error(gctx.New(), err)
			} else {
				data, ok := act.Data.(map[string]interface{})
				if ok {
					err = c.validate(data)
					if err != nil {
						c.Deliver(service.Action().NewErrorMessage(err.Error()))
					} else {
						switch act.Action {
						case consts.ActionSendMessage:
							msg, err := service.Action().GetMessage(act)
							if err == nil {
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
								msg.ReceivedAt = gtime.New()
								c.Manager.ReceiveMessage(&chatConnMessage{
									Msg:  msg,
									Conn: c,
								})
							}

						}
					}
				}
			}

		}
	}
}

// Deliver 投递消息
func (c *client) Deliver(act *chat.Action) {
	c.Send <- act
}

// SendMsg 发消息
func (c *client) SendMsg() {
	for {
		select {
		case act := <-c.Send:
			msgByte, err := service.Action().MarshalAction(act)
			if err != nil {
				g.Log().Error(gctx.New(), err)
			} else {
				err := c.Conn.WriteMessage(websocket.TextMessage, msgByte)
				if err == nil {
					switch act.Action {
					case consts.ActionMoreThanOne:
						c.Close()
					case consts.ActionOtherLogin:
						c.Close()
					case consts.ActionReceiveMessage:
						msg, ok := act.Data.(*relation.CustomerChatMessages)
						if ok {
							dao.CustomerChatMessages.Ctx(gctx.New()).WherePri(msg.Id).
								Data("send_at", time.Now().Unix()).Update()
						}
					default:
					}
				} else {
					c.Close()
					return
				}
			}
		case <-c.CloseSignal:
			return
		}
	}
}
