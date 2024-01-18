package main

import (
	"fmt"
	"os"

	cli "github.com/davidroman0O/gogog/cmds"
	"github.com/davidroman0O/gogog/cmds/agent"
	"github.com/davidroman0O/gogog/cmds/login"
)

func main() {
	var rootCmd = cli.Cmd()
	rootCmd.AddCommand(agent.Cmd())
	rootCmd.AddCommand(login.Cmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
