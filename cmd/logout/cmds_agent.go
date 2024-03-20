package logout

import (
	"fmt"

	"github.com/davidroman0O/gogog/data"
	"github.com/davidroman0O/gogog/types"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove your login informations",
		Long:  `.`,
		Run: func(cmd *cobra.Command, args []string) {
			if !data.Has[types.GogAuthenticationChrome]() {
				fmt.Println("You have no authentication information")
				return
			}
			if err := data.Delete[types.GogAuthenticationChrome](); err != nil {
				fmt.Printf("Failed to delete your authentication information: %v\n", err)
				return
			}
		},
	}
}
