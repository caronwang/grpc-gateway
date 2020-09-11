# HTTP网关
etcd3 API全面升级为gRPC后，同时要提供REST API服务，维护两个版本的服务显然不太合理，所以grpc-gateway诞生了。
通过protobuf的自定义option实现了一个网关，服务端同时开启gRPC和HTTP服务，HTTP服务接收客户端请求后转换为grpc
请求数据，获取响应后转为json数据返回给客户端。

## 安装grpc-gateway
```shell script
$ go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
```

## 目录结构
```shell script
|—— hello_http/
    |—— client/
        |—— main.go   // 客户端
    |—— server/
        |—— main.go   // GRPC服务端
    |—— server_http/
        |—— main.go   // HTTP服务端
|—— proto/
    |—— google       // googleApi http-proto定义
        |—— api
            |—— annotations.proto
            |—— annotations.pb.go
            |—— http.proto
            |—— http.pb.go
    |—— hello_http/
        |—— hello_http.proto   // proto描述文件
        |—— hello_http.pb.go   // proto编译后文件
        |—— hello_http_pb.gw.go // gateway编译后文件
```

## 代码详解
Step 1. 编写proto描述文件：proto/hello_http.proto
```prototext
syntax = "proto3";

package hello_http;
option go_package = "hello_http";

import "google/api/annotations.proto";

// 定义Hello服务
service HelloHTTP {
    // 定义SayHello方法
    rpc SayHello(HelloHTTPRequest) returns (HelloHTTPResponse) {
        // http option
        option (google.api.http) = {
            post: "/example/echo"
            body: "*"
        };
    }
}

// HelloRequest 请求结构
message HelloHTTPRequest {
    string name = 1;
}

// HelloResponse 响应结构
message HelloHTTPResponse {
    string message = 1;
}
```

Step 2. 编译proto
```shell script
$ cd proto

# 编译google.api
$ protoc -I . --go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:. google/api/*.proto

# 编译hello_http.proto
$ protoc -I . --go_out=plugins=grpc,Mgoogle/api/annotations.proto=github.com/jergoo/go-grpc-example/proto/google/api:. hello_http/*.proto

# 编译hello_http.proto gateway
$ protoc --grpc-gateway_out=logtostderr=true:. hello_http/hello_http.proto
```

注意这里需要编译google/api中的两个proto文件，同时在编译hello_http.proto时使用M参数指定引入包名，最后使用
grpc-gateway编译生成hello_http_pb.gw.go文件，这个文件就是用来做协议转换的，查看文件可以看到里面生成的
http handler，处理proto文件中定义的路由"example/echo"接收POST参数，调用HelloHTTP服务的客户端请求grpc
服务并响应结果。


参考：http://www.topgoer.com/%E5%BE%AE%E6%9C%8D%E5%8A%A1/gRPC/HTTP%E7%BD%91%E5%85%B3.html