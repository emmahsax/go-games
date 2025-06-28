package main

import (
	"fmt"
	"os"

	"github.com/emmahsax/go-games/games/deliveryDash"
	"github.com/emmahsax/go-games/games/yourGame" // <----- Change the name of your game here
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go-games",
		Short: "Go games",
	}

	cmd.DisableAutoGenTag = true

	cmd.AddCommand(deliveryDash.NewCommand())
	cmd.AddCommand(yourGame.NewCommand()) // <----- Change the name of your game here

	return cmd
}

func main() {
	rootCmd := NewCommand()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
