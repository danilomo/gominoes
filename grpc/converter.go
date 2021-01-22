package grpc

import (
	g "github.com/danilomo/gominoes/src"
)

// IsValid returns true if the message is a proper message or false
// if all fields are set to false and nill
func (m *Message) IsValid() bool {
	return true
}

// IsSkip Returns true if the message represents a skip signal
func (m *Message) IsSkip() bool {
	switch m.Content.(type) {
	case *Message_Skip:
		return true
	default:
		return false
	}
}

// IsJoin Returns true if the message represents a join signal
func (m *Message) IsJoin() bool {
	switch m.Content.(type) {
	case *Message_Join:
		return true
	default:
		return false
	}
}

// AsMove cast the message to a Move structure, if the message is a valid move message
func (m *Message) AsMove() (*g.Move, bool) {
	switch m.Content.(type) {
	case *Message_Move:
		contents := m.GetMove()
		move := &g.Move{
			Player:       int(contents.Player),
			HandPosition: int(contents.HandPosition),
			Side:         g.BoardSide(contents.Side),
		}

		return move, true
	default:
		return nil, false
	}
}

// AsUpdate cast the message to a Update structure, if the message is a valid update message
func (m *Message) AsUpdate() (*g.Update, bool) {
	switch m.Content.(type) {
	case *Message_Update:
		contents := m.GetUpdate()
		move := &g.Move{
			Player:       int(contents.Move.Player),
			HandPosition: int(contents.Move.HandPosition),
			Side:         g.BoardSide(contents.Move.Side),
		}
		update := &g.Update{
			Move: move,
			Turn: int(contents.Turn),
		}

		return update, true
	default:
		return nil, false
	}
}

// AsResponse cast the message to a Response structure, if the message is a valid response message
func (m *Message) AsResponse() (*g.Response, bool) {
	switch m.Content.(type) {
	case *Message_Update:
		contents := m.GetResponse()

		response := &g.Response{
			Ok:       contents.Ok,
			ErrorMsg: contents.Error,
		}

		return response, true
	default:
		return nil, false
	}
}

// ToProtobuf converts a gominoe's Message into the protobuf format
func ToProtobuf(msg g.Message) *Message {
	if msg.IsJoin() {
		return &Message{Content: &Message_Join{}}
	}

	if msg.IsSkip() {
		return &Message{Content: &Message_Skip{}}
	}

	if move, ok := msg.AsMove(); ok {
		return moveMessageToProtobuf(move)
	}

	if update, ok := msg.AsUpdate(); ok {
		return updateMessageToProtobuf(update)
	}

	return &Message{}
}

func moveMessageToProtobuf(move *g.Move) *Message {
	return &Message{Content: &Message_Move{
		Move: moveToProtobuf(move),
	}}
}

func updateMessageToProtobuf(update *g.Update) *Message {
	return &Message{Content: &Message_Update{
		Update: &Update{
			Turn: int32(update.Turn),
			Move: moveToProtobuf(update.Move),
		},
	}}
}

func moveToProtobuf(move *g.Move) *Move {
	if move == nil {
		return &Move{}
	}

	return &Move{
		HandPosition: int32(move.HandPosition),
		Player:       int32(move.Player),
		Side:         int32(move.Side),
	}
}
