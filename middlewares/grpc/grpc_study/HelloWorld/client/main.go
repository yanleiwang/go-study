package main

import (
	"context"
	studyapi "geektime-go-study/study/middleware/grpc_study/api/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {
	// step1 创建 grpc 连接
	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// step2 创建 grpc 客户端
	c := studyapi.NewSimpleServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// step3 调用rpc方法
	r, err := c.SayHello(ctx, &studyapi.HelloReq{
		Name: "wangyanlei",
		Age:  18,
		Phones: []*studyapi.HelloReq_PhoneNumber{
			{
				Number: "123",
				Type:   studyapi.HelloReq_PHONE_TYPE_HOME,
			},
			{
				Number: "567",
			},
		},
	})
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("resp: %s", r.Message)
}
