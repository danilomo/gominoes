package gominoes

// TODO: add logging framework

import (
	"math/rand"
)

// Gomino defines a data structure for a domino piece
type Gomino [2]int

var invalid = Gomino{-1, -1}

// matches tests if another gomino piece can be joined to the given piece.
// It returns a tuple containing the result of the test, and the other domino piece
// in the right orientation.
// E.g:
// g := Gomino{1, 6}
// g2 := Gomino{5, 6}
// g3 := Gomino{6, 3}
// t, g2 = g.matches(g2, Right) --> g2 is now {6, 5}, piece was flipped
// t, g3 = g.matches(g3, Right) --> g3 is still {6, 3}, no need to flip
func (g Gomino) matches(other Gomino, position BoardSide) (bool, Gomino) {

	if position == Left {
		if g[0] == other[0] {
			return true, other.flip()
		} else if g[0] == other[1] {
			return true, other
		}
	} else {
		if g[1] == other[0] {
			return true, other
		} else if g[1] == other[1] {
			return true, other.flip()
		}
	}

	return false, invalid
}

func (g Gomino) flip() Gomino {
	return [2]int{g[1], g[0]}
}

// Player represents a player in a game
type Player struct {
	Hand []Gomino
	Name string
}

// Game represents the current snapshot of a Gominoes game
type Game struct {
	Players []Player
	Board   []Gomino
	Turn    int
}

// BoardSide specifies the position to insert a gomino in the board (left or right)
type BoardSide int

const (
	// Left side of the board
	Left BoardSide = iota
	// Right side of the board
	Right
)

// Move defines a move performed by a player (e.g put the 2nd piece of the 3rd player in the right
// side of the board).
type Move struct {
	Player       int
	HandPosition int
	Side         BoardSide
}

// PlayMove updates a Game struct with a player move passed as argument
// It returns a tuple containing a flag to indicate if the move was valid, and
// an error message to be displayed in case the move was not valid. Game is not updated
// if the move is not valid, the turn stays the same
func (game *Game) PlayMove(move Move) (bool, string) {

	if move.Player >= len(game.Players) {
		return false, "Invalid player number"
	} else if game.Turn != move.Player {
		return false, "This is not your turn"
	} else if move.HandPosition >= len(game.Players[move.Player].Hand) || move.HandPosition < 0 {
		return false, "Invalid piece position for player"
	}

	playerPiece := game.Players[move.Player].Hand[move.HandPosition]

	if len(game.Board) == 0 {
		game.Board = append(game.Board, playerPiece)
		game.removeFromHand(move)
		game.increaseTurn()

		return true, ""
	}

	var boardIndex = 0
	if move.Side == Right {
		boardIndex = len(game.Board) - 1
	}
	boardPiece := game.Board[boardIndex]

	matches, pieceToInsert := boardPiece.matches(playerPiece, move.Side)
	if matches {
		game.insertPiece(pieceToInsert, move.Side)
		game.removeFromHand(move)
		game.increaseTurn()

		return true, ""
	}

	return false, "Piece does not match board"
}

// Skip makes the player of the current move to skip his/her turn. It should be called
// only if the player does not have a gomino that matches the left or right side of the board,
// otherwise, the turn cannot be skiped and the player must do a valid move instead.
func (game *Game) Skip(player int) bool {
	if game.canPlayMove(player) {
		return false
	}

	game.increaseTurn()
	return true
}

// Winner returns the player id for the game winner, if the game ended. If the
// game has no winner, a negative value is returned instead
func (game *Game) Winner() int {
	for i := 0; i < len(game.Players); i++ {
		if len(game.Players[i].Hand) == 0 {
			return i
		}
	}

	return -1
}

// canPlayMove returns true if the player is able to do a move for the current turn
func (game *Game) canPlayMove(player int) bool {

	if player > 0 || player >= len(game.Players) {
		return false
	}

	leftPiece := game.Board[0]
	rightPiece := game.Board[len(game.Board)-1]
	for _, p := range game.Players[player].Hand {

		matchLeft, _ := leftPiece.matches(p, Left)
		matchRight, _ := rightPiece.matches(p, Right)

		if matchLeft || matchRight {
			return true
		}
	}

	return false
}

func (game *Game) insertPiece(piece Gomino, pos BoardSide) {
	if pos == Right {
		game.Board = append(game.Board, piece)
	} else {
		game.Board = append([]Gomino{piece}, game.Board...)
	}
}

// FirstTurn returns true if no player did a move, otherwise it returns false
func (game *Game) FirstTurn() bool {
	return len(game.Board) == 0
}

func (game *Game) increaseTurn() {
	game.Turn++
	if game.Turn >= len(game.Players) {
		game.Turn = 0
	}
}

func (game *Game) removeFromHand(move Move) {
	player, pos := move.Player, move.HandPosition
	hand := &game.Players[player].Hand
	copy((*hand)[pos:], (*hand)[pos+1:])
	*hand = (*hand)[:len(*hand)-1]
}

// NewGame creates a Gominoes game in the initial state
func NewGame(numberOfPlayers int) *Game {
	board := []Gomino{}
	players := make([]Player, numberOfPlayers)

	for i := 0; i <= 6; i++ {
		for j := 0; j <= i; j++ {
			board = append(board, Gomino{i, j})
		}
	}

	rand.Shuffle(len(board), func(i, j int) { board[i], board[j] = board[j], board[i] })

	boardLen := len(board)
	for i := 0; i < boardLen; i++ {
		player := i % numberOfPlayers
		piece := board[0]
		board = board[1:]

		players[player].Hand = append(players[player].Hand, piece)
	}

	return &Game{
		Players: players,
		Board:   board,
		Turn:    0,
	}
}
