# Gominoes

Gominoes is a Dominoes game implemented in the Go programming language. It implements the game rules as an ADT (Abstract Data Type, struct + methods) and also provide concurrency operations to implement a distributed game. It is a Cobra-based application, providing three commands:

* start - starts the gRPC game server 
* start-telnet - starts the TCP text-based game server
* play - client application for the gRPC game server
