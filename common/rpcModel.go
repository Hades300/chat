package common

type LoginRequest struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

type LoginReply struct {
	Code      int    `json:"Code"`
	AuthToken string `json:"autotoken"`
}

type RegisterRequest struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

type RegisterReply struct {
	Code    int    `json:"Code"`
	Message string `json:"message"`
}

type CheckAuthRequest struct {
	Token string `json:"token"`
}

type CheckAuthReply struct {
	UserName string `json:"username"`
	Status   int    `json:"status"`
}
