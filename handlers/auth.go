package handlers

import "github.com/gofiber/fiber/v2"

// The client `cli` will send it's own cookies to allow the agent to be able to connect
func PostCookies(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
