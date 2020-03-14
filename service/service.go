package service

type service struct{}

func NewService() *service {
	return new(service)
}

func (s *service) Run() {
	go initUserServiceRpcServer()
}
