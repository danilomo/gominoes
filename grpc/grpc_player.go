package grpc

import (
	"fmt"
	"time"

	g "github.com/danilomo/gominoes/src"
)

type grpcPlayer struct {
	stream       GameService_JoinGameServer
	playerNumber int
}

func (player *grpcPlayer) Greeting() {
	// nothing to do on greeting
}

func (player *grpcPlayer) Update(update *g.Update) {
	message := g.UpdateMsg(*update)
	updateProto := ToProtobuf(message)

	player.stream.SendMsg(updateProto)
}

func (player *grpcPlayer) SendResponse(response *g.Response) {
	message := g.ResponseMsg(*response)
	responseProto := ToProtobuf(message)

	player.stream.SendMsg(responseProto)
}

func (player *grpcPlayer) ReadMessage() g.Message {
	for {
		msg, err := player.stream.Recv() // TODO: error handling

		fmt.Println("Leu: ", msg, ", ", err, ", playernum: ", player.playerNumber)

		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		return msg
	}

}
