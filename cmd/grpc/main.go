package main

import (
	"gRPCqueue/internal/grpc/messages"
	"gRPCqueue/messagepb"
	"google.golang.org/grpc"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":1488")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	messagepb.RegisterMsgServiceServer(grpcServer, messages.Service{})
	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}

}
