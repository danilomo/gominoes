package grpc

import (
	reflect "reflect"
	"testing"

	g "github.com/danilomo/gominoes/src"
)

func TestConvertProtobufToMessage(t *testing.T) {
	joinMessage := &Message{Content: &Message_Join{}}

	if !joinMessage.IsJoin() {
		t.Error("Join message was not successfully converted")
		return
	}

	moveMessageProto := &Message{Content: &Message_Move{
		Move: &Move{
			Player:       1,
			HandPosition: 1,
			Side:         1,
		},
	}}
	move := &g.Move{
		Player:       1,
		HandPosition: 1,
		Side:         1,
	}

	moveProto, _ := moveMessageProto.AsMove()

	if !reflect.DeepEqual(moveProto, move) {
		t.Error("Unable to convert protobuf message to Message interface")
		return
	}

}
