package gominoes

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Message is a union type for a message exchanged between the game server and the players.
// It can be a signal with a boolean field set to true, or a data message with a single field set
// to a non-nil pointer.
type Message struct {
	move     *Move
	update   *Update
	response *Response
	skip     bool
	join     bool
}

// Update describes a message sent to all players after some players make a valid move.
// It also signals it is your turn (if Turn == your player id)
type Update struct {
	Move *Move
	Turn int
}

// Response describes the outcome of a move. If it is not okay, the move is not valid and
// the game's turn does not increase until the player make a valid move
type Response struct {
	ok       bool
	errorMsg string
}

// IsValid returns true if the message is a proper message or false
// if all fields are set to false and nill
func (msg Message) IsValid() bool {
	return !(msg.move == nil && msg.update == nil && !msg.skip && !msg.join)
}

// IsSkip Returns true if the message represents a skip signal
func (msg Message) IsSkip() bool {
	return msg.skip
}

// IsJoin Returns true if the message represents a join signal
func (msg Message) IsJoin() bool {
	return msg.skip
}

// AsMove cast the message to a Move structure, if the message is a valid move message
func (msg Message) AsMove() (*Move, bool) {
	if msg.move != nil {
		return msg.move, true
	}

	return nil, false
}

// AsUpdate cast the message to a Update structure, if the message is a valid update message
func (msg Message) AsUpdate() (*Update, bool) {
	if msg.update != nil {
		return msg.update, true
	}

	return nil, false
}

// AsResponse cast the message to a Response structure, if the message is a valid response message
func (msg Message) AsResponse() (*Response, bool) {
	if msg.response != nil {
		return msg.response, true
	}

	return nil, false
}

// SkipMsg creates a new skip message
func SkipMsg() Message {
	return Message{
		skip: true,
	}
}

// JoinMsg creates a new join message
func JoinMsg() Message {
	return Message{
		join: true,
	}
}

// MoveMsg creates a new move message
func MoveMsg(move Move) Message {
	return Message{
		move: &move,
	}
}

// UpdateMsg creates a new update message
func UpdateMsg(update Update) Message {
	return Message{
		update: &update,
	}
}

// OkMsg creates a response message with status = ok
func OkMsg() Message {
	return Message{
		response: &Response{true, ""},
	}
}

// ErrorMsg creates a response message with status = ok
func ErrorMsg(msg string) Message {
	return Message{
		response: &Response{false, msg},
	}
}

var splitter = regexp.MustCompile(",")
var converters = map[string]func([]string) Message{
	"skip": func(args []string) Message {
		return SkipMsg()
	},
	"join": func(args []string) Message { return JoinMsg() },
	"move": func(args []string) Message {
		fmt.Println(args)

		if len(args) != 2 {
			return Message{}
		}

		position, _ := strconv.Atoi(args[0])
		side := Left

		if args[1] == "R" {
			side = Right
		}

		return MoveMsg(Move{
			HandPosition: position,
			Side:         side,
		})
	},
}

// Convert parses a string into a Message struct
func Convert(value string) (Message, bool) {
	value = strings.TrimSpace(value)
	fields := splitter.Split(value, -1)

	if len(fields) == 0 {
		return Message{}, false
	}
	key := strings.Trim(fields[0], " ")
	if converterFunction, ok := converters[key]; ok {
		response := converterFunction(fields[1:])

		if !response.IsValid() {
			return Message{}, false
		}

		return response, true
	}

	return Message{}, false
}
