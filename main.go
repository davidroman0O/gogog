package main

import (
	"fmt"
	"os"

	cli "github.com/davidroman0O/gogog/cmd"
	"github.com/davidroman0O/gogog/cmd/agent"
	"github.com/davidroman0O/gogog/cmd/login"
	"github.com/davidroman0O/gogog/cmd/logout"
	"github.com/davidroman0O/gogog/cmd/server"
)

func main() {
	var rootCmd = cli.Cmd()
	rootCmd.AddCommand(agent.Cmd())
	rootCmd.AddCommand(login.Cmd())
	rootCmd.AddCommand(logout.Cmd())
	rootCmd.AddCommand(server.Cmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
