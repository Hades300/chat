package context

import (
	"chat/common"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Hcon struct {
	conn         *websocket.Conn
	MessageQueue chan common.Message // Message类型多样，使用空接口队列
	UserName     string
	roomName     string
	hotelName    string
}

const (
	DefaultRoomName = "hades300"

	DefaultHotelName = "heyao"

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var beehiveContext = NewBeehiveContext()

/**
新建默认的连接对象
roomName 类似于群id
hotelName 是个更大的集群概念了，先留个接口
*/
func NewDefaultHcon(conn *websocket.Conn) *Hcon {
	h := new(Hcon)
	h.roomName = DefaultRoomName
	h.hotelName = DefaultHotelName
	h.MessageQueue = make(chan common.Message, 10)
	h.conn = conn
	return h
}

// TODO: 之后考虑更换logrus，但是目前使用自带的，先把主体逻辑写好

/**
Start 是每个用户都要执行的部分，负责
read distribute

在Read前判断连接是否断开。
*/
func (h *Hcon) Start() {
	// 每个连接注册一次
	beehiveContext.Register(h)
	defer func() {
		beehiveContext.Leave(h)
		h.conn.Close()
	}()

	// 读超时检查
	h.conn.SetReadDeadline(time.Now().Add(pongWait))
	h.conn.SetPongHandler(func(data string) error {
		h.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var data common.Message
		if err := h.conn.ReadJSON(&data); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break // 产生了错误 大多数情况是用户直接关闭网页 直接break
		}
		// 在消息分发时，先处理登录消息，验证登录后，才处理其他消息，在消息处理之前，验证有无水平越权
		beehiveContext.MessageDistribute(h, data)
	}
}

/**
从用户的专属conn中读出Message 并发送到客户端


*/
func (h *Hcon) ConsumeMessageQueue() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			h.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := h.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case message, ok := <-h.MessageQueue:
			h.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				h.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			h.conn.WriteJSON(message)
		}
	}
}

/**
消息合法验证
*/

func (h *Hcon) Vertify(message common.Message) bool {
	if h.UserName == message.Source {
		return true
	} else {
		return false
	}
}
