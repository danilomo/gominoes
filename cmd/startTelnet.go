package cmd

import (
	gominoes "github.com/danilomo/gominoes/src"
	"github.com/spf13/cobra"
)

// startTelnetCmd represents the startTelnet command
var startTelnetCmd = &cobra.Command{
	Use:   "start-telnet",
	Short: "Starts the Gominoes server in the Telnet mode",
	Long:  ``,
	Run:   startTelnet,
}

func init() {
	rootCmd.AddCommand(startTelnetCmd)
}

func startTelnet(cmd *cobra.Command, args []string) {
	gameServer := gominoes.StartServer(4, 8001)
	gameServer.Wait()
}
