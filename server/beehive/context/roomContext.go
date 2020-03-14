package context

import (
	"chat/common"
	"errors"
	"fmt"
	"sync"
	"time"
)

type RoomContext struct {
	rooms       map[string]map[string]chan common.Message
	roomCount   map[string]uint32
	roomChannel chan common.Message
	lock        sync.RWMutex
}

func NewRoomContext() *RoomContext {
	rContext := new(RoomContext)
	rContext.rooms = make(map[string]map[string]chan common.Message)
	rContext.lock = sync.RWMutex{}
	rContext.roomCount = make(map[string]uint32)
	rContext.roomChannel = make(chan common.Message, 50)
	return rContext
}

/**
将用户加入房间 广播人数消息
*/

func (rc *RoomContext) AddUserToRoom(ch *Hcon, roomName string) error {
	rc.lock.Lock()
	defer rc.lock.Unlock()
	if rc.rooms[roomName] == nil {
		rc.rooms[roomName] = make(map[string]chan common.Message)
	}
	rc.rooms[roomName][ch.UserName] = ch.MessageQueue
	rc.roomCount[roomName] += 1
	//debug
	fmt.Printf("发送了Count Message 当前人数%d\n", rc.roomCount[roomName])
	rc.RoomBroadCast(roomName, common.Message{
		CreateAt:    time.Time{},
		Source:      roomName,
		Target:      roomName,
		Content:     fmt.Sprintf("%d", rc.roomCount[roomName]),
		MessageType: common.CountMessage,
	})
	return nil
}

/**
删除某房间中的某用户
*/

func (rc *RoomContext) DeleteUserFromRoom(ch *Hcon, roomName string) error {
	rc.lock.Lock()
	defer rc.lock.Unlock()
	if ok := rc.UserInRoom(ch.UserName, roomName); !ok {
		return nil
	}
	delete(rc.rooms[roomName], ch.UserName)
	rc.roomCount[roomName]--
	countMessage := common.Message{
		CreateAt:    time.Now(),
		Source:      roomName,
		Target:      roomName,
		Content:     fmt.Sprintf("%d", rc.roomCount[roomName]),
		MessageType: common.CountMessage,
	}
	rc.RoomBroadCast(roomName, countMessage)
	return nil
}

func (rc *RoomContext) DeleteUserFromAllRoom(ch *Hcon) error {
	for roomName, _ := range rc.rooms {
		rc.DeleteUserFromRoom(ch, roomName)
	}
	return nil
}

func (rc *RoomContext) UserInRoom(userName string, roomName string) bool {
	if rc.rooms[roomName][userName] != nil {
		return true
	}
	return false
}

/**
房间内广播
*/

func (rc *RoomContext) RoomBroadCast(roomName string, message common.Message) error {
	if ok := rc.RoomEst(roomName); !ok {
		return errors.New("Room " + roomName + " not exists When broadCasting")
	}
	for _, userMessageQueue := range rc.rooms[roomName] {
		userMessageQueue <- message
	}
	return nil
}

/**
懒狗选择这个名字 来当 RoomExists的缩写
*/

func (rc *RoomContext) RoomEst(roomName string) bool {
	_, ok := rc.rooms[roomName]
	return ok
}

/**
监听房间信道 阻塞
*/

func (rc *RoomContext) RoomSupervisor() {
	// TODO: 删掉这个DeBUG用的输出
	for {
		select {
		case message := <-rc.roomChannel:
			if message.MessageType == common.CountMessage {
				fmt.Printf("有用户 进入房间，目前人数%s\n", message.Content)
			}
			if message.MessageType == common.RoomMessage {
				fmt.Printf("用户%s:%s", message.Source, message.Content)
			}
			rc.RoomBroadCast(message.Target, message)
		}
	}
}
