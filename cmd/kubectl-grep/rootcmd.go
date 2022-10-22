package main

import (
	"log"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "kubectl lookup STRING",
	SilenceUsage: true, // for when RunE returns an error
	Example:      "kubectl lookup STRING \n",
	Args:         cobra.MinimumNArgs(1),
	RunE:         run,
}

func run(command *cobra.Command, args []string) error {
	containerState(args[0])

	return nil
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
