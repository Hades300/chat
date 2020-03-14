package context

import (
	"chat/common"
	"chat/server/beehive/rpc"
	"log"
	"time"
)

type BeehiveContext struct {
	uContext *UserContext
	rContext *RoomContext
}

func NewBeehiveContext() *BeehiveContext {
	beehiveContext := new(BeehiveContext)
	beehiveContext.uContext = NewUserContext()
	beehiveContext.rContext = NewRoomContext()
	return beehiveContext
}

/**
在 uContext && rContext 写原子操作
在BeehiveContext中封装一下同名函数，加入逻辑控制
*/

/**
注册到uContext的users中
*/

func (this *BeehiveContext) Register(ch *Hcon) error {
	return this.uContext.AddUser(ch)

}

/**
注销用户的连接 alias for deleteuser
*/

func (this *BeehiveContext) Leave(ch *Hcon) error {
	return this.DeleteUser(ch)
}

func (this *BeehiveContext) AddUserToRoom(ch *Hcon, roomName string) error {
	return this.rContext.AddUserToRoom(ch, roomName)
}

/**
删除操作不考虑返回错误，不管有无，都返回了。但是为了一致性，保留error返回值
*/

func (this *BeehiveContext) DeleteUser(ch *Hcon) error {
	this.rContext.DeleteUserFromAllRoom(ch)
	this.uContext.DeleteUser(ch)
	return nil
}

/**
凑凑行数，这里可能不写也可以
*/
func (this *BeehiveContext) DeleteUserFromRoom(ch *Hcon, roomName string) error {
	return this.rContext.DeleteUserFromRoom(ch, roomName)
}

/**
一旦RoomChannel中读到信息，通过广播按房间名单进行分发
可能的问题：某一用户阻塞的话，会影响到其他用户
*/

func (this *BeehiveContext) RoomBroadCast(roomName string, message common.Message) error {
	return this.rContext.RoomBroadCast(roomName, message)
}

func (this *BeehiveContext) RoomSuperviser() {
	this.rContext.RoomSupervisor()
	return
}

func (this *BeehiveContext) Send2User(userName string, message common.Message) error {
	return this.uContext.Send2User(userName, message)
}

func (this *BeehiveContext) MessageDistribute(h *Hcon, message common.Message) {
	// 首先处理登录Message
	if message.MessageType == common.AuthMessage {
		if !this.UserAuth(h, message) {
			infoMessage := common.Message{
				CreateAt:    time.Time{},
				Source:      "System",
				Target:      "anonymous",
				Content:     "登录失败 无效Token",
				MessageType: common.InfoMessage,
			}
			h.MessageQueue <- infoMessage
			return
		}
	}
	if h.UserName == "" {
		infoMessage := common.Message{
			CreateAt:    time.Time{},
			Source:      "System",
			Target:      "anonymous",
			Content:     "请先登录",
			MessageType: common.InfoMessage,
		}
		h.MessageQueue <- infoMessage
		return
	}
	if !h.Vertify(message) {
		infoMessage := common.Message{
			CreateAt:    time.Time{},
			Source:      "System",
			Target:      "anonymous",
			Content:     "越权",
			MessageType: common.InfoMessage,
		}
		h.MessageQueue <- infoMessage
		return
	} else {
		switch message.MessageType {
		case common.RoomMessage:
			if err := this.RoomBroadCast(message.Target, message); err != nil {
				log.Printf("Message from %s fail to send to room  %s because %s", message.Source, message.Target, err)
				return
			}

		case common.UserMessage:
			h.MessageQueue <- message
			if err := this.Send2User(message.Target, message); err != nil {
				log.Printf("Message from %s fail to send to User %s because %s", message.Source, message.Target, err)
				return

			}
		}
	}
}

func (this *BeehiveContext) UserAuth(h *Hcon, loginmessage common.Message) bool {
	uService := rpc.UserServiceRpcClient{}
	userName := uService.CheckAuthByToken(loginmessage.Content)
	if userName == "" {
		return false
	} else {
		h.UserName = userName
		return true
	}
}

/**
两个查
*/

func (this *BeehiveContext) UserEst(userName string) bool {
	return this.uContext.UserEst(userName)
}

func (this *BeehiveContext) RoomEst(roomName string) bool {
	return this.rContext.RoomEst(roomName)
}
