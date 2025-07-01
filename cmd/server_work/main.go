package main

import (
	"context"
	"gRPCqueue/messagepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"os"
	"time"

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
func GetMessageByOffset(messageClient messagepb.MsgServiceClient, Offset uint32, Limit uint32) ([]string, error) {
	out, err := messageClient.GetMessageByOffset(context.Background(), &messagepb.GetMessageByOffsetRequest{
		Offset: Offset,
		Limit:  Limit,
	})
	if err != nil {
		return nil, err
	}
	return ([]string)(out.Messages), nil
}
func getMessageByID(messageClient messagepb.MsgServiceClient, Id string) (string, error) {
	out, err := getMessageByID(messageClient, Id)
	if err != nil {
		return "", err
	}
	return (string)(out), nil
}

func main() {
	newClient, err := grpc.NewClient("localhost:1488", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	idx_of_msg := 0
	messageClient := messagepb.NewMsgServiceClient(newClient)
	file, err := os.Create("cmd/server_work/log/" + "message_" + strconv.Itoa(idx_of_msg) + "-" + strconv.Itoa(idx_of_msg+100) + ".txt")
	if err != nil {
		panic(err)
	}
	var content string
	var idempotencyKey string

	for {
		cnt_for_new_file := 0
		content = "number: " + strconv.Itoa(rand.Intn(100000))
		idempotencyKey = "test" + strconv.Itoa(rand.Intn(100000))
		out, err := WriteMessage(messageClient, content, idempotencyKey)
		if err != nil {
			panic(err)
		}
		file.WriteString(out + " " + content + " " + idempotencyKey + " idx: " + strconv.Itoa(idx_of_msg) + "\n")
		cnt_for_new_file++
		idx_of_msg += 1
		if cnt_for_new_file == 100 {
			file.Close()
			break
		}

		time.Sleep(2 * time.Second)
	}
}
