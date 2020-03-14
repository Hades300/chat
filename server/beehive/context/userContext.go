package context

import (
	"chat/common"
	"errors"
	"sync"
)

type UserContext struct {
	users map[string]chan common.Message
	lock  sync.RWMutex
}

/**
userContext 没有什么需要配置的地方。因此使用
unexport的函数
*/

func NewUserContext() *UserContext {
	uContext := new(UserContext)
	uContext.users = make(map[string]chan common.Message)
	uContext.lock = sync.RWMutex{}
	return uContext
}

/**
加入Users 已存在则返回error
*/
func (uc *UserContext) AddUser(ch *Hcon) error {
	uc.lock.Lock()
	defer uc.lock.Unlock()
	if _, ok := uc.users[ch.UserName]; ok {
		return errors.New("UserName 已存在")
	}
	uc.users[ch.UserName] = ch.MessageQueue
	return nil
}

/**
存在即删除
*/

func (uc *UserContext) DeleteUser(ch *Hcon) error {
	uc.lock.Lock()
	defer uc.lock.Unlock()
	delete(uc.users, ch.UserName)
	return nil
}

/*
	检查user是否存在
*/
func (uc *UserContext) UserEst(userName string) bool {
	_, ok := uc.users[userName]
	return ok
}

/**
投递到某个用户
*/

func (uc *UserContext) Send2User(userName string, message common.Message) error {
	if ok := uc.UserEst(userName); !ok {
		return errors.New("Error in Send2User : No Such User")
	}
	uc.users[userName] <- message
	return nil
}
