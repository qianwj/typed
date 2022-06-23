package tests

import (
	"fmt"
	"go.uber.org/fx"
	"testing"
)

type Option interface {
	GetOption() any
}

type GrpcOption struct {
	addr        string
	healthCheck bool
}

func (g *GrpcOption) GetOption() any {
	return g
}

func Addr() *GrpcOption {
	return &GrpcOption{addr: "123"}
}

func EnableHealthCheck() *GrpcOption {
	return &GrpcOption{healthCheck: true}
}

type Options struct {
	opts []Option `group:"options"`
}

type Service interface {
	Register()
}

type appleService struct{}

func newAppleService() Service {
	return &appleService{}
}

func (a *appleService) Register() {
	println("apple")
}

type bananaService struct{}

func newBananaService() Service {
	return &bananaService{}
}

func (a *bananaService) Register() {
	println("banana")
}

type Configuration struct {
	Addr string
}

type serviceModule struct {
	fx.In
	Services []Service `group:"svc"`
}

type Server struct{}

func NewServer(c *Configuration, svc serviceModule) *Server {
	fmt.Printf("addr: %s\r\n", c.Addr)
	for _, service := range svc.Services {
		service.Register()
	}
	return &Server{}
}

func TestDI(t *testing.T) {
	module := fx.Module("service",
		fx.Provide(fx.Annotated{Target: newAppleService, Group: "svc"}),
		fx.Provide(fx.Annotated{Target: newBananaService, Group: "svc"}),
	)
	app := fx.New(module, fx.Invoke(NewServer), fx.Supply(&Configuration{Addr: "abc"}))
	app.Run()
	t.Log("test complete")
}

func TestMap(t *testing.T) {
	m := make(map[string]int, 1)
	m["aaa"] = 1
	m["aa"] = 2
	for s, i := range m {
		t.Logf("%s:%d", s, i)
	}
	a := make([]int, 0, 10)
	a = append(a, 1, 1)
	t.Log(a[1])
}
