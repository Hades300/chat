package common

type LoginForm struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

type RegisterForm struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

type LoginResult struct {
	Code      int    `json:"code"`
	AuthToken string `json:"token"`
}

type RegisterResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
