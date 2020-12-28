package gominoes

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	messages := map[string]Message{
		"skip":       SkipMsg(),
		"join":       JoinMsg(),
		"skip     ":  SkipMsg(),
		"join   ":    JoinMsg(),
		"move,2,R":   MoveMsg(Move{0, 2, Right}),
		"move,1,L":   MoveMsg(Move{0, 1, Left}),
		"move,1,L\n": MoveMsg(Move{0, 1, Left}),
	}

	for k, v := range messages {
		converted, ok := Convert(k)
		move, isMove := v.AsMove()

		if !ok {
			t.Error("Unable to convert message: [" + k + "]")
			continue
		}

		if cMove, cIsMove := converted.AsMove(); isMove && cIsMove && !reflect.DeepEqual(cMove, move) {
			fmt.Println("Key: " + k)
			fmt.Println(v)
			fmt.Println(converted)
			t.Error("Unable to play a successful move with a empty board")
		}
	}
}
