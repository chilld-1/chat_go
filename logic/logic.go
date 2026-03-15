package logic

type Logic struct{}

func New() *Logic {
	return &Logic{}
}

func (l *Logic) Run() {
	// 启动 gRPC 服务器
	err := StartGrpcServer()
	if err != nil {
		panic(err)
	}
}
