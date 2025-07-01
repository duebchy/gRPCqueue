package messages

import (
	"context"
	"errors"
	"gRPCqueue/messagepb"
	"github.com/google/uuid"
)

var Data = make([]message, 0, 100)

type Service struct {
	messagepb.UnimplementedMsgServiceServer
}
type message struct {
	ID             string
	content        string
	IdempotencyKey string
}

func (s Service) CreateMessage(ctx context.Context, request *messagepb.CreateMessageRequest) (*messagepb.CreateMessageResponse, error) {
	for _, val := range Data {
		if val.IdempotencyKey == request.IdempotencyKey {
			return nil, errors.New("Already exists!")
		}
	}
	msg := message{
		content:        request.Content,
		IdempotencyKey: request.IdempotencyKey,
		ID:             uuid.New().String(),
	}

	Data = append(Data, msg)
	return &messagepb.CreateMessageResponse{
		ID: msg.ID,
	}, nil

}
func (s Service) findMessageById(ID string) (message, error) {
	for _, m := range Data {
		if m.ID == ID {
			return m, nil
		}
	}

	return message{}, errors.New("user not found")
}
func (s Service) GetMessage(ctx context.Context, request *messagepb.GetMessageByIDRequest) (*messagepb.GetMessageByIDResponse, error) {
	messageInfo, err := s.findMessageById(request.ID)
	if err != nil {
		return nil, err
	}
	return &messagepb.GetMessageByIDResponse{Messages: ([]string{
		messageInfo.content,
	})}, nil
}
func (s Service) GetMessageByOffset(ctx context.Context, request *messagepb.GetMessageByOffsetRequest) (*messagepb.GetMessageByOffsetResponse, error) {
	Limit := request.Limit
	Offset := request.Offset
	if Limit > uint32(len(Data))-1 || Offset >= Limit {
		return &messagepb.GetMessageByOffsetResponse{Messages: ([]string{})}, errors.New("limit or Offset out of range")
	}
	temp := Data[Offset : Limit+Offset]
	newData := make([]string, 0, len(temp)*2)
	for _, m := range temp {
		newData = append(newData, m.content)
		newData = append(newData, m.ID)
	}

	return &messagepb.GetMessageByOffsetResponse{Messages: newData}, nil
}
func New() *message {
	return &message{}
}
