package grpc

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	g "github.com/danilomo/gominoes/src"
	grpc "google.golang.org/grpc"
)

// PlayGame aaa
func PlayGame(address, gameID string, playerNum int) error {

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := NewGameServiceClient(conn)

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	clientDeadline := time.Now().Add(10000 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)

	defer cancel()

	stream, error := c.JoinGame(ctx)

	if error != nil {
		fmt.Println("Error: ", error)
		return error
	}

	startGameLoop(stream, gameID, playerNum)

	return nil
}

func startGameLoop(stream GameService_JoinGameClient, gameID string, playerNum int) {
	reader := bufio.NewReader(os.Stdin)

	stream.Send(&Message{Content: &Message_Join{&Join{
		GameId: gameID,
	}}})

	for {
		serverMessage, error := stream.Recv()

		if error != nil {
			// TODO: better error treatment
			fmt.Println("Error: ", error)
			continue
		}

		update, ok := serverMessage.AsUpdate()

		// TODO: debug message
		fmt.Println("Update: ", update.Turn, ", ", ok)

		if ok && update.Turn == playerNum {
			playMove(stream, reader)
		}
	}
}

func playMove(stream GameService_JoinGameClient, reader *bufio.Reader) {
	fmt.Println("It is your turn. Make your move: ")

	for {
		text, err := reader.ReadString('\n')

		if err != nil {
			// TODO: better error treatment
			fmt.Println("Error: ", err)
			continue
		}

		text = strings.Replace(text, "\n", "", -1)

		message, ok := g.Convert(text)
		_, isMove := message.AsMove()

		if !(isMove || message.IsSkip()) {
			fmt.Println("Invalid move. Try again: ")
			continue
		}

		stream.Send(ToProtobuf(message))

		responseMessage, err := stream.Recv()

		if err != nil {
			// TODO: better error treatment
			fmt.Println("Error: ", err)
			continue
		}

		response, ok := responseMessage.AsResponse()

		if !ok {
			// TODO: better error treatment
			fmt.Println("Invalid response from server")
			continue
		}

		if !response.Ok {
			// TODO: better error treatment
			fmt.Println("Invalid move: ", response.ErrorMsg)
			continue
		}

		break
	}
}
