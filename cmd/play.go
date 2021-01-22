package cmd

import (
	"fmt"
	"strconv"

	"github.com/danilomo/gominoes/grpc"
	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Client application for Gominoes in the GRPC mode.",
	Long:  ``,
	Run:   playGame,
}

func init() {
	rootCmd.AddCommand(playCmd)
}

func playGame(cmd *cobra.Command, args []string) {
	playerNum, err := strconv.Atoi(args[2])

	if err != nil {
		fmt.Println("Deu pau ", err)
		return
	}

	grpc.PlayGame(args[0], args[1], playerNum)
}
