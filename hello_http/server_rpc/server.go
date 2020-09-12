package main

import (
	"fmt"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	. "grpc-gateway/proto/hello_http"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

// 定义helloService并实现约定的接口
type helloService struct{}

// HelloService Hello服务
var HelloService = helloService{}

// SayHello 实现Hello服务接口
func (h helloService) SayHello(ctx context.Context, in *HelloHTTPRequest) (*HelloHTTPResponse, error) {
	resp := new(HelloHTTPResponse)
	resp.Message = fmt.Sprintf("Hello %s.", in.Name)

	return resp, nil
}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}

	var s *grpc.Server

	//无认证
	s = grpc.NewServer()

	// 注册HelloService
	RegisterHelloHTTPServer(s, HelloService)

	fmt.Println("Listen on " + Address)
	s.Serve(listen)
}
