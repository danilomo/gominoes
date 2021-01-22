package cmd

import (
	"github.com/danilomo/gominoes/grpc"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the Gominoes game server in the GRPC mode",
	Long:  ``,
	Run:   startServer,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func startServer(cmd *cobra.Command, args []string) {
	grpc.StartServer(args[0])
}
