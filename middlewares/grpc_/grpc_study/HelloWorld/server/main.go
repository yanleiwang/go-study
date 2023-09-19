package main

import (
	"context"
	studyapi "geektime-go-study/study/middleware/grpc_study/api/gen"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	studyapi.UnimplementedSimpleServiceServer // 必须组合 UnimplementedXXX
}

// SayHello 实现服务的接口 在proto中定义的所有服务都是接口
func (s *server) SayHello(ctx context.Context, req *studyapi.HelloReq) (*studyapi.HelloResp, error) {
	log.Println(req)
	return &studyapi.HelloResp{
		Message: "Hello, " + req.Name,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 创建grpc服务
	s := grpc.NewServer()

	// 注册 GreeterServer 到 grpc 服务上
	studyapi.RegisterSimpleServiceServer(s, &server{})

	// 让grpc 监听 连接
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
