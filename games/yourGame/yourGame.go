package yourGame // <----- Change the name of your game here

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "your-game", // <----- Change how users call your game here

		// Aliases to make calling your game easier
		// Aliases: []string{""},

		// Brief description of your game
		// Short:   "",

		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}

	return cmd
}
