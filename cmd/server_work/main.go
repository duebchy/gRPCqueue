package main

import (
	"context"
	"fmt"
	"gRPCqueue/messagepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func Pow10(n int, x int) int {
	if n == 1 {
		return x
	} else {
		return Pow10(n-1, x*x)
	}
}
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
func end(idx int) {
	file, err := os.Create("cmd/server_work/log/lastIndex.txt")
	if err != nil {
		panic(err)
	}
	file.WriteString(strconv.Itoa(idx))
	file.Close()
}
func GetLastIndex() int {
	file, err := os.Open("cmd/server_work/log/lastIndex.txt")
	if err != nil {
		panic(err)
	}
	data := make([]byte, 64)

	for {
		n, err := file.Read(data)
		if err == io.EOF { // если конец файла
			break
		}
		data = data[:n]
	}

	file.Close()
	file, err = os.OpenFile("cmd/server_work/log/lastindex.txt", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	file.Close()
	idx, err := strconv.Atoi(string(data))
	if err != nil {
		panic(err)
	}
	return idx
}
func main() {
	//init client
	newClient, err := grpc.NewClient("localhost:1488", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	messageClient := messagepb.NewMsgServiceClient(newClient)
	//init variables
	idx_of_msg := GetLastIndex()
	endIndex := idx_of_msg + 100
	cnt_for_new_file := idx_of_msg % 100
	var content string
	var idempotencyKey string

	fmt.Println(idx_of_msg)
	//init log files
	var fileName string
	if idx_of_msg < 100 {
		fileName = "cmd/server_work/log/0-100.txt"
	} else {
		fileName = fmt.Sprintf("cmd/server_work/log/%d00-%d00.txt", idx_of_msg/Pow10((len(strconv.Itoa(idx_of_msg))-1), 10), endIndex/Pow10((len(strconv.Itoa(endIndex))-1), 10))
	}
	fmt.Println(idx_of_msg, endIndex)
	fmt.Println(fileName)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		panic(err)
	}
	//file, err = os.Create("cmd/server_work/log/" + strconv.Itoa(idx_of_msg) + "-" + strconv.Itoa(idx_of_msg+100) + ".txt")
	//if err != nil {
	//	panic(err)
	//}

	for i := 0; i < 67; i++ {

		content = "number: " + strconv.Itoa(rand.Intn(100000))
		idempotencyKey = "test" + strconv.Itoa(rand.Intn(1000000000000000000))

		out, err := WriteMessage(messageClient, content, idempotencyKey)
		if err != nil {
			panic(err)
		}

		_, err = file.WriteString(out + " " + content + " " + idempotencyKey + " idx: " + strconv.Itoa(idx_of_msg) + "\n")
		if err != nil {
			panic(err)
		}
		cnt_for_new_file += 1
		idx_of_msg += 1
		if cnt_for_new_file == 100 {
			cnt_for_new_file = 0
			file.Close()
			file, err = os.Create("cmd/server_work/log/" + strconv.Itoa(idx_of_msg) + "-" + strconv.Itoa(idx_of_msg+100) + ".txt")

		}
		time.Sleep(50 * time.Millisecond)
	}
	end(idx_of_msg)
}
