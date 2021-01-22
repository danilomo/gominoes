package gominoes

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

type telnetPlayer struct {
	server     *GameServer
	number     int
	connection net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
}

func (player telnetPlayer) Greeting() {
	writer := player.writer

	writer.WriteString("Welcome, player! " + fmt.Sprint(player.number) + "\n")
	writer.Flush()
}

func (player telnetPlayer) Update(updateMsg *Update) {
	writer := player.writer
	server := player.server
	number := player.number

	writer.WriteString("Board: " + fmt.Sprint(server.game.Board) + "\n")
	writer.WriteString("Your hand: " + fmt.Sprint(server.game.Players[number].Hand) + "\n")
	writer.WriteString("Enter your move: ")
	writer.Flush()
}

func (player telnetPlayer) SendResponse(responseMsg *Response) {
	writer := player.writer

	if !responseMsg.Ok {
		writer.WriteString("Error: " + responseMsg.ErrorMsg + "\n" +
			"Enter your move: ")
		writer.Flush()
	} else {
		// nothing to do here, player made a correct move
	}

}

func (player telnetPlayer) ReadMessage() Message {
	writer := player.writer
	reader := player.reader

	for {
		line, err := reader.ReadString('\n')

		fmt.Println("Read [", line, "] from player ", player.number, " and error ", err)

		if err != nil {
			// TODO: handle error
			fmt.Println("|> ", err.Error())
			continue
		}

		// discarding CR and newline
		line = line[0 : len(line)-2]

		msg, ok := Convert(line)

		if !ok {
			writer.WriteString("Invalid command <" + line + ">\n" +
				"Enter your move: ")
			writer.Flush()
			continue
		}

		if move, ok := msg.AsMove(); ok {
			move.Player = player.number

			return MoveMsg(*move)
		}

		return msg
	}
}

// StartServer starts a Gominoes match as a TCP server
func StartServer(players int, port int) *GameServer {
	game := NewGame(players)
	server := game.Start()

	startServerAux(server, func() (net.Listener, error) {
		return net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	})

	return server
}

func joinPlayer(playerNumber int, conn net.Conn, server *GameServer) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	remotePlayer := telnetPlayer{
		number:     playerNumber,
		connection: conn,
		server:     server,
		reader:     reader,
		writer:     writer,
	}

	server.StartPlayerLoop(playerNumber, &remotePlayer)
}

func startServerAux(server *GameServer, serverGen func() (net.Listener, error)) {
	players := len(server.players)
	l, err := serverGen()

	if err != nil {
		// TODO: handle error
		return
	}
	defer l.Close()

	for i := 0; i < players; i++ {
		c, err := l.Accept()

		if err != nil {
			// TODO: handle error
			continue
		}

		go joinPlayer(i, c, server)
	}

}
