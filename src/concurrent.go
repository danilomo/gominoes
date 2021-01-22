package gominoes

import (
	"fmt"
	"sync"
)

// GameServer is a handler for a Gominoes match running assynchronously.
// It waits until all players join the match, then it starts the game loop.
// To register new players, get a reference to the channel players[i] and
// exchange messages according to the protocol.
type GameServer struct {
	game          *Game
	players       []chan Message
	wg            *sync.WaitGroup
	activePlayers int
}

// Wait blocks the current goroutine until the gominoes match finishes
func (server *GameServer) Wait() {
	server.wg.Wait()
}

// ActivePlayers returns the number of players that joined the game server
func (server *GameServer) ActivePlayers() int {
	return server.activePlayers
}

// starts the game loop. Should be started as a goroutine: go start();
func (server *GameServer) start() {
	chans := server.players
	game := server.game

	for i := 0; i < len(server.players); i++ {
		<-chans[i] // TODO: verify it is a join message
		server.activePlayers = server.activePlayers + 1
	}

	// TODO: debug message
	fmt.Println("Jogo começou")

	update := UpdateMsg(Update{Turn: game.Turn})

	for game.Winner() < 0 {
		for i := 0; i < len(chans); i++ {
			chans[i] <- update
		}

		for {
			// TODO: debug message
			fmt.Println("O tabuleiro é: ", game.Board)
			fmt.Println("Mão do jogador: ", game.Players[game.Turn].Hand)

			msg := <-chans[game.Turn]

			// TODO: debug message
			fmt.Println("Jogador ", game.Turn, " enviou ", msg)

			if msg.IsSkip() {
				turn := game.Turn
				skipped := game.Skip(game.Turn)

				if !skipped {
					chans[game.Turn] <- ErrorMsg("Cannot skip, player has pieces to play")
					continue
				}

				chans[turn] <- OkMsg()
				update = UpdateMsg(Update{Turn: game.Turn})
				break
			}

			move, isMove := msg.AsMove()

			if !isMove {
				chans[game.Turn] <- ErrorMsg("Invalid message.")
				continue
			}

			// TODO: debug message
			fmt.Println("Player ", game.Turn, " played ", move)

			turn := game.Turn
			ok, err := game.PlayMove(*move)

			if !ok {
				chans[game.Turn] <- ErrorMsg(err)
				continue
			}

			chans[turn] <- OkMsg()
			update = UpdateMsg(Update{Turn: game.Turn, Move: move})

			if game.Winner() >= 0 {
				server.wg.Done()
				return
			}

			// TODO: debug message
			fmt.Println("Successful move. Now the board is: ", game.Board)

			break
		}
	}

	// TODO: update players about who is the winner

	server.wg.Done()
}

// Start starts a new gominoes match as a goroutine, and returns the GameServer
// handler that is connected to it
func (game *Game) Start() *GameServer {
	numPlayers := len(game.Players)
	channels := make([]chan Message, numPlayers)

	for i := 0; i < numPlayers; i++ {
		channels[i] = make(chan Message)
	}

	server := &GameServer{
		game:    game,
		players: channels,
	}

	server.wg = &sync.WaitGroup{}
	server.wg.Add(1)
	go server.start()

	return server
}
