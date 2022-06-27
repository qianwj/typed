package grpc

import (
	"fmt"
	"github.com/qianwj/typed/fx/options"
	"google.golang.org/grpc"
	"testing"
)

func TestApplication_RegisterService(t *testing.T) {
	app := NewApp(options.Address(":8081"))
	app.RegisterService(newTestService)
	app.Run()
}

type testService struct{}

func newTestService() Service {
	return &testService{}
}

func (s *testService) Register(srv *grpc.Server) {
	fmt.Println("mock grpc service register")
}
