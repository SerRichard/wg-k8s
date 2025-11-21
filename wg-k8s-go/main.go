package main

import (
	"os"
	
	"github.com/spf13/cobra"
)

func WireguardK8S() *cobra.Command {
	command := &cobra.Command{
		Use:   "wgk8s",
		Short: "Wireguard K8s Toolkit.",
		Long: `
Command line tool for interacting with the Wireguard K8s toolkit. Currently, you can 
use the wizard setup remote locations and output the config update for your other sites.`,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	// command.AddCommand(server.ServerCommand())
	// command.AddCommand(daemon.DaemonCommand())

	command.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	return command	
}

func main() {

	if err := WireguardK8S().Execute(); err != nil {
		os.Exit(1)
	}

}
