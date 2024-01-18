package agent

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "agent",
		Short: ".",
		Long:  `.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}
