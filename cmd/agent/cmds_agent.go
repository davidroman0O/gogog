package agent

import (
	"github.com/davidroman0O/gogog/data"
	"github.com/davidroman0O/gogog/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

/// note: I want to do some grpc but man i have no time these days

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "agent",
		Short: ".",
		Long:  `.`,
		Run: func(cmd *cobra.Command, args []string) {
			data.SetRuntime(false)

			app := fiber.New()

			app.Post("/cookies", handlers.PostCookies)

			app.Listen(":3000")
		},
	}
}
