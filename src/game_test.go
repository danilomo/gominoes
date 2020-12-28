package gominoes

import (
	"reflect"
	"testing"
)

func TestPlacePieceEmptyBoard(t *testing.T) {
	game := Game{
		Players: []Player{
			{
				Hand: []Gomino{{1, 2}, {2, 2}, {3, 1}, {6, 4}},
				Name: "Player 1",
			},
		},
	}
	move := Move{Player: 0, HandPosition: 1, Side: Left}
	success, errMsg := game.PlayMove(move)

	expected := Game{
		Players: []Player{
			{
				Hand: []Gomino{{1, 2}, {3, 1}, {6, 4}},
				Name: "Player 1",
			},
		},
		Board: []Gomino{{2, 2}},
	}

	if !success || len(errMsg) > 0 {
		t.Error("Unable to play a successful move with a empty board")
	}

	if !reflect.DeepEqual(game, expected) {
		t.Error("PlayMove on an empty board did not produce the expected result")
	}
}

func TestPlayMoves(t *testing.T) {
	game := Game{
		Players: []Player{
			{
				Hand: []Gomino{{1, 2}, {2, 2}, {3, 1}, {6, 4}},
				Name: "Player 1",
			},
			{
				Hand: []Gomino{{5, 5}, {6, 6}, {1, 4}},
				Name: "Player 2",
			},
			{
				Hand: []Gomino{{3, 3}, {5, 6}, {2, 5}, {1, 1}},
				Name: "Player 3",
			},
		},
		Board: []Gomino{{3, 4}, {4, 2}},
	}

	moves := []Move{
		{Player: 0, HandPosition: 2, Side: Left},
		{Player: 1, HandPosition: 2, Side: Left},
		{Player: 2, HandPosition: 2, Side: Right},
		{Player: 0, HandPosition: 2, Side: Left},
	}

	for _, move := range moves {
		success, msg := game.PlayMove(move)

		if !success || len(msg) > 0 {
			t.Error("Unable to play a successful move: " + msg)
			return
		}
	}

	expected := Game{
		Players: []Player{
			{
				Hand: []Gomino{{1, 2}, {2, 2}},
				Name: "Player 1",
			},
			{
				Hand: []Gomino{{5, 5}, {6, 6}},
				Name: "Player 2",
			},
			{
				Hand: []Gomino{{3, 3}, {5, 6}, {1, 1}},
				Name: "Player 3",
			},
		},
		Board: []Gomino{{6, 4}, {4, 1}, {1, 3}, {3, 4}, {4, 2}, {2, 5}},
		Turn:  1,
	}

	if !reflect.DeepEqual(game, expected) {
		t.Error("Game is in an unexpected state after moves")
	}

}
