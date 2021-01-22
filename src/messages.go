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
type Message interface {
	// IsValid returns true if the message is a proper message or false
	// if all fields are set to false and nill
	IsValid() bool

	// IsSkip Returns true if the message represents a skip signal
	IsSkip() bool

	// IsJoin Returns true if the message represents a join signal
	IsJoin() bool

	// AsMove cast the message to a Move structure, if the message is a valid move message
	AsMove() (*Move, bool)

	// AsUpdate cast the message to a Update structure, if the message is a valid update message
	AsUpdate() (*Update, bool)

	// AsResponse cast the message to a Response structure, if the message is a valid response message
	AsResponse() (*Response, bool)
}

type messageStruct struct {
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
	Ok       bool
	ErrorMsg string
}

// IsValid returns true if the message is a proper message or false
// if all fields are set to false and nill
func (msg messageStruct) IsValid() bool {
	return !(msg.move == nil && msg.update == nil && !msg.skip && !msg.join)
}

// IsSkip Returns true if the message represents a skip signal
func (msg messageStruct) IsSkip() bool {
	return msg.skip
}

// IsJoin Returns true if the message represents a join signal
func (msg messageStruct) IsJoin() bool {
	return msg.skip
}

// AsMove cast the message to a Move structure, if the message is a valid move message
func (msg messageStruct) AsMove() (*Move, bool) {
	if msg.move != nil {
		return msg.move, true
	}

	return nil, false
}

// AsUpdate cast the message to a Update structure, if the message is a valid update message
func (msg messageStruct) AsUpdate() (*Update, bool) {
	if msg.update != nil {
		return msg.update, true
	}

	return nil, false
}

// AsResponse cast the message to a Response structure, if the message is a valid response message
func (msg messageStruct) AsResponse() (*Response, bool) {
	if msg.response != nil {
		return msg.response, true
	}

	return nil, false
}

// SkipMsg creates a new skip message
func SkipMsg() Message {
	return &messageStruct{
		skip: true,
	}
}

// JoinMsg creates a new join message
func JoinMsg() Message {
	return &messageStruct{
		join: true,
	}
}

// MoveMsg creates a new move message
func MoveMsg(move Move) Message {
	return &messageStruct{
		move: &move,
	}
}

// ResponseMsg creates a new response message
func ResponseMsg(response Response) Message {
	return &messageStruct{
		response: &response,
	}
}

// UpdateMsg creates a new update message
func UpdateMsg(update Update) Message {
	return &messageStruct{
		update: &update,
	}
}

// OkMsg creates a response message with status = ok
func OkMsg() Message {
	return &messageStruct{
		response: &Response{true, ""},
	}
}

// ErrorMsg creates a response message with status = ok
func ErrorMsg(msg string) Message {
	return &messageStruct{
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
			return &messageStruct{}
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

// Convert parses a string into a Message struct.
// Examples:
// * move,1,L -> inserts the second (0 based) piece for the hand for the current player
//               in the left position of the board
// * skip -> represents a "skip" move (the user has no pieces that match the left or right side
//           of the board)
func Convert(value string) (Message, bool) {
	value = strings.TrimSpace(value)
	fields := splitter.Split(value, -1)

	if len(fields) == 0 {
		return &messageStruct{}, false
	}
	key := strings.Trim(fields[0], " ")
	if converterFunction, ok := converters[key]; ok {
		response := converterFunction(fields[1:])

		if !response.IsValid() {
			return &messageStruct{}, false
		}

		return response, true
	}

	return &messageStruct{}, false
}
