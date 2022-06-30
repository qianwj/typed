package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"testing"
)

func TestApplication_RegisterService(t *testing.T) {
	app := NewServer(Address(":8081"))
	app.RegisterService(newTestService)
	app.Provide()
}

type testService struct{}

func newTestService() Service {
	return &testService{}
}

func (s *testService) Register(srv *grpc.Server) {
	fmt.Println("mock grpc service register")
}
