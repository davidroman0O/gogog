package main

import (
	"fmt"
	"os"

	cli "github.com/davidroman0O/gogog/cmds"
	"github.com/davidroman0O/gogog/cmds/agent"
	"github.com/davidroman0O/gogog/cmds/login"
	"github.com/davidroman0O/gogog/cmds/logout"
	"github.com/davidroman0O/gogog/cmds/web"
)

func main() {
	var rootCmd = cli.Cmd()
	rootCmd.AddCommand(agent.Cmd())
	rootCmd.AddCommand(login.Cmd())
	rootCmd.AddCommand(logout.Cmd())
	rootCmd.AddCommand(web.Cmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
