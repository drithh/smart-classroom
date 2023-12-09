package fiber

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func SetupFiber(db *sqlx.DB) {
	app := fiber.New()

	// Define your Fiber routes and handlers here
	// For example, a simple route to handle HTTP requests
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})

	// Start the Fiber web server
	go func() {
		err := app.Listen(":3000")
		if err != nil {
			fmt.Printf("Error starting server: %v", err)
		}
	}()

}
