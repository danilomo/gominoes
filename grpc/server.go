package grpc

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	g "github.com/danilomo/gominoes/src"
	grpc "google.golang.org/grpc"
)

var player int

type server struct {
	games map[string]*g.GameServer
}

// StartServer starts a GRPC server that can serve multiple gominoes.GameServer matches
func StartServer(port string) {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	gamesMap := make(map[string]*g.GameServer)
	RegisterGameServiceServer(s, &server{
		games: gamesMap,
	})

	log.Printf("Starting gRPC listener on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) JoinGame(stream GameService_JoinGameServer) error {
	first, _ := stream.Recv()

	switch first.Content.(type) {
	case *Message_Join:
		join := first.GetJoin()

		if len(join.GameId) == 0 {
			return errors.New("Cannot join game with empty ID")
		}
		game, ok := s.games[join.GameId]

		if !ok {
			gameStruct := g.NewGame(4)
			game = gameStruct.Start()
			s.games[join.GameId] = game
		}

		return s.joinGameAux(game, stream)
	default:
		return errors.New("First message should be a join message")
	}
}

// TODO: this is a hack, should be removed
var playerN int64

func (s *server) joinGameAux(game *g.GameServer, stream GameService_JoinGameServer) error {
	playerNumber := game.ActivePlayers()

	player := grpcPlayer{
		stream:       stream,
		playerNumber: int(playerN),
	}

	// TODO: debug message
	fmt.Println("Player ", playerN, " joined!")

	// TODO: remove atomic.AddInt64 usage
	atomic.AddInt64(&playerN, 1)
	game.StartPlayerLoop(playerNumber, &player)

	return nil
}
