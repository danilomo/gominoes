package gominoes

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

// StartServer starts a Gominoes match as a TCP server
func StartServer(players int, port int) {
	game := NewGame(players)
	server := game.Start()

	startServerAux(server, func() (net.Listener, error) {
		return net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	})
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

func joinPlayer(i int, c net.Conn, server *GameServer) {
	server.players[i] <- JoinMsg()
	reader := bufio.NewReader(c)
	writer := bufio.NewWriter(c)

	writer.WriteString("Welcome, player! " + fmt.Sprint(i) + "\n")
	writer.Flush()

	for {
		updateMsg := <-server.players[i]

		update, ok := updateMsg.AsUpdate()

		if !ok || update.Turn != i {
			continue
		}

		writer.WriteString("Board: " + fmt.Sprint(server.game.Board) + "\n")
		writer.WriteString("Your hand: " + fmt.Sprint(server.game.Players[i].Hand) + "\n")
		writer.WriteString("Enter your move: ")
		writer.Flush()

		for {
			line, err := reader.ReadString('\n')

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

			if msg.IsSkip() {
				server.players[i] <- msg
				response, ok := (<-server.players[i]).AsResponse()

				if ok && !response.ok {
					writer.WriteString("Error: " + response.errorMsg + "\n" +
						"Enter your move: ")
					writer.Flush()
					continue
				}

				break
			}

			move, isMove := msg.AsMove()

			if !isMove {
				writer.WriteString("Invalid move <" + line + ">\n" +
					"Enter your move: ")
				writer.Flush()
				continue
			}

			move.Player = i
			server.players[i] <- MoveMsg(*move)
			response, ok := (<-server.players[i]).AsResponse()

			if ok && !response.ok {
				writer.WriteString("Error: " + response.errorMsg + "\n" +
					"Enter your move: ")
				writer.Flush()
				continue
			}

			break
		}
	}
}
