// Package gominoes provides the implementation of the synchronous and asynchronous code
// for the Gominoes game. It provides some ADTs that can be used to create a concurrent and remote
// gominoes server in any network protocol:
// * GameServer   - represents a virtual game server (or game room, as you wish to name it)
//             for a single Gominoes match
// * Message      - represents a message exchanged between game server and remote players. It is an union type
//			   shaped as an interface
// * RemotePlayer - represents a remote player, written using any network protocol (sockets, websockets, grpc, etc)
// This package contains also a TCP text-based server implementation (available at gominoes.StartServer)
package gominoes
