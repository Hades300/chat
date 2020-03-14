package common

import "time"

type Message struct {
	CreateAt    time.Time
	Source      string // 应该是个主键  如qq号，但是目前不考虑使用数据库，先使用强制使用不可重复的username
	Target      string
	Content     string // 内容 目前是纯文本
	MessageType MessageType
}

type MessageType string

const UserMessage MessageType = "userMessage"
const RoomMessage MessageType = "roomMessage"
const EnterMessage MessageType = "enterMessage"
const CountMessage MessageType = "countMessage"
const InfoMessage MessageType = "infoMessage"
const AuthMessage MessageType = "authMessage"
