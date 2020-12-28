package main

import (
	"time"

	gominoes "github.com/danilomo/gominoes/src"
)

func main() {

	gominoes.StartServer(4, 8001)

	time.Sleep(3600 * time.Second)
}
