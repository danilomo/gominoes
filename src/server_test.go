package gominoes

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strings"
	"testing"

	"github.com/danilomo/gominoes/faketcp"
)

func TestServer(t *testing.T) {
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

	gameServer := game.Start()

	moves := []string{
		strings.Join([]string{
			"move,2,L\r\n",
			"move,2,L\r\n",
			"skip\r\n",
		}, ""),
		strings.Join([]string{
			"move,2,L\r\n",
			"move,1,L\r\n",
			"move,0,R\r\n",
		}, ""),
		strings.Join([]string{
			"move,2,R\r\n",
			"move,1,L\r\n",
		}, ""),
	}

	fakeServer := fakeServer(moves)

	startServerAux(gameServer, func() (net.Listener, error) { return fakeServer, nil })

	gameServer.Wait()

	expected := Game{
		Players: []Player{
			{
				Hand: []Gomino{{1, 2}, {2, 2}},
				Name: "Player 1",
			},
			{
				Hand: []Gomino{},
				Name: "Player 2",
			},
			{
				Hand: []Gomino{{3, 3}, {1, 1}},
				Name: "Player 3",
			},
		},
		Board: []Gomino{{5, 6}, {6, 6}, {6, 4}, {4, 1}, {1, 3}, {3, 4}, {4, 2}, {2, 5}, {5, 5}},
		Turn:  2,
	}

	if !reflect.DeepEqual(game, expected) {
		json, _ := json.Marshal(game)
		fmt.Println(string(json))
		t.Error("Game is in an unexpected state after moves")
	}

	if 1 != game.Winner() {
		t.Error("game.Winner() returned an unexpected result", game.Winner())
	}
}

func fakeServer(moves []string) net.Listener {
	var i int
	generator := func() net.Conn {
		con := faketcp.FakeClient(moves[i])
		i++
		return con
	}

	return faketcp.FakeServer(generator)
}
