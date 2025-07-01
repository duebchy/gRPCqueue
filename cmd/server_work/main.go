package main

import (
	"context"
	"fmt"
	"gRPCqueue/messagepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"math/rand"
	"strconv"
)

func WriteMessage(messageClient messagepb.MsgServiceClient, content string, IdempotencyKey string) (string, error) {
	out, err := messageClient.CreateMessage(context.Background(), &messagepb.CreateMessageRequest{
		Content:        content,
		IdempotencyKey: IdempotencyKey,
	})
	if err != nil {
		panic(err)
	}

	return out.ID, err
}
func GetMessageByOffser(messageClient messagepb.MsgServiceClient, Offset uint32, Limit uint32) ([]string, error) {
	out, err := messageClient.GetMessageByOffset(context.Background(), &messagepb.GetMessageByOffsetRequest{
		Offset: Offset,
		Limit:  Limit,
	})
	if err != nil {
		return nil, err
	}
	return ([]string)(out.Messages), nil
}
func main() {
	newClient, err := grpc.NewClient("localhost:1488", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	messageClient := messagepb.NewMsgServiceClient(newClient)
	for i := 0; i < 8; i++ {
		ID, err := WriteMessage(messageClient, "yo"+strconv.Itoa(i), "2007"+strconv.Itoa(rand.Intn(100000000000)))
		if err != nil {
			panic(err)
		}
		fmt.Println("ID:", ID)
	}
	out, err := GetMessageByOffser(messageClient, 0, 5)
	fmt.Println("offset out: ", out)
}
