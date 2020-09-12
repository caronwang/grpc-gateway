package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	pb "grpc-gateway/proto/hello_http" // 引入proto包
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

func main() {
	var conn *grpc.ClientConn
	var err error

	// 普通链接
	conn, err = grpc.Dial(Address, grpc.WithInsecure())

	if err != nil {
		grpclog.Fatalln(err)
	}
	defer conn.Close()

	// 初始化客户端
	c := pb.NewHelloHTTPClient(conn)

	// 调用方法
	req := &pb.HelloHTTPRequest{Name: "gRPC"}
	res, err := c.SayHello(context.Background(), req)

	if err != nil {
		grpclog.Fatalln(err)
	}

	fmt.Println(res.Message)
}
