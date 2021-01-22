package gominoes

// RemotePlayer represents a remote client as a Gominoes player. The methods from this interface
// should implement the network calls that will push events to the player (Greeting, Update, SendResponse)
// and that methods that will pull a move from a player in his/her turn.
// It should use a bi-directional networking protocol like normal sockets, websockets, GRPC stream, etc.
type RemotePlayer interface {
	// Greeting sends the initial greeting message to the player
	Greeting()

	// Update sends the update message to each player after each turn in the Gominoes match
	Update(updateMsg *Update)

	// SendResponse sends the message in response to a move, indicating the move was accepted by
	// the game server, or if it is an erroneous move (in this case the turn does not advance
	// and the user is asked for a move again)
	SendResponse(responseMsg *Response)

	// Reads a move or skip message from the user, when it is his/her turn
	ReadMessage() Message
}

// StartPlayerLoop starts an event loop for a remote Gominoes player.
// The playerNumber should be a number between 0 and 3, which identifies the player in a Gominoes match
// The player argument represents an object which is able to read from and write to the remote connection,
// for the remote client, whatever the network protocol is
func (server *GameServer) StartPlayerLoop(playerNumber int, player RemotePlayer) {
	playerChannel := server.players[playerNumber]
	playerChannel <- JoinMsg()
	player.Greeting()

	for {
		updateMsg := <-playerChannel
		update, ok := updateMsg.AsUpdate()

		if !ok || update.Turn != playerNumber {
			continue
		}

		player.Update(update)

		for {
			messageFromPlayer := player.ReadMessage()

			if messageFromPlayer.IsSkip() {
				playerChannel <- messageFromPlayer
				response, ok := (<-playerChannel).AsResponse()

				if ok && !response.Ok {
					player.SendResponse(response)
					continue
				}

				break
			}

			move, isMove := messageFromPlayer.AsMove()

			if !isMove {
				errorResponse := &Response{Ok: false, ErrorMsg: "Invalid move."}
				player.SendResponse(errorResponse)

				continue
			}

			playerChannel <- MoveMsg(*move)
			response, ok := (<-playerChannel).AsResponse()

			if ok && !response.Ok {
				player.SendResponse(response)
				continue
			}

			break
		}
	}
}
