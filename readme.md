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
├── gen.sh		//proto文件编译命令
├── hello_http	
│   ├── client
│   │   └── client.go	//客户端程序
│   ├── server_http
│   │   └── main.go		//http服务端程序
│   └── server_rpc
│       └── server.go	//rpc服务端程序
├── proto
   └── hello_http			
       ├── hello_http.pb.go		//protobuf生成的go文件
       ├── hello_http.pb.gw.go	//protobuf生成的gateway文件
       └── hello_http.proto	//protobuf文件
```



## 使用说明

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

Step 2. 编译proto，运行gen.sh文件
```shell script
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:. \
  proto/hello_http/*.proto

protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:. \
  proto/hello_http/*.proto
```
编译后生成两个go文件

       ├── hello_http.pb.go		//protobuf生成的go文件
       ├── hello_http.pb.gw.go	//protobuf生成的gateway文件
Step 3. 运行服务端程序

```shell
#运行RPC服务端
go run  hello_http/server_rpc/server.go

#运行HTTP服务端
go run  hello_http/server_http/main.go
```

Step 4. 测试

测试RPC客户端

```shell
#运行客户端
go run  hello_http/client/main.go

#得到返回
Hello gRPC.
```

测试HTTP响应

```
curl -H "Content-Type: application/json" -X POST -d '{"name": "123"}'  http://127.0.0.1:8080/example/echo

#得到返回
{"message":"Hello 123."}
```





参考：http://www.topgoer.com/%E5%BE%AE%E6%9C%8D%E5%8A%A1/gRPC/HTTP%E7%BD%91%E5%85%B3.html