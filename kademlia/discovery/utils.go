package discovery

import (
	"context"
	"github.com/gogo/protobuf/proto"
	"github.com/perlin-network/noise/internal/protobuf"
	"github.com/perlin-network/noise/protocol"
	"github.com/pkg/errors"
)

const (
	ServiceID            = 5
	OpCodePing           = 1
	OpCodePong           = 2
	OpCodeLookupRequest  = 3
	OpCodeLookupResponse = 4
)

type SendHandler interface {
	Request(ctx context.Context, target []byte, body *protocol.MessageBody) (*protocol.MessageBody, error)
	Broadcast(body *protocol.MessageBody) error
}

func ToMessageBody(serviceID int, opcode int, content proto.Message) (*protocol.MessageBody, error) {
	raw, err := proto.Marshal(content)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to marshal content")
	}
	msg := &protobuf.Message{
		Message: raw,
		Opcode:  uint32(opcode),
	}
	msgBytes, err := msg.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to marshal message")
	}
	body := &protocol.MessageBody{
		Service: uint16(serviceID),
		Payload: msgBytes,
	}
	return body, nil
}

func ParseMessageBody(body *protocol.MessageBody) (*protobuf.Message, error) {
	if body == nil || len(body.Payload) == 0 {
		return nil, errors.New("body is empty")
	}
	var msg protobuf.Message
	if err := proto.Unmarshal(body.Payload, &msg); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal payload")
	}
	return &msg, nil
}
