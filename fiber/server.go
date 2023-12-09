package fiber

import (
	"fmt"

	db "github.com/drithh/smart-classroom/database"
	"github.com/gofiber/fiber/v2"
)

func SetupFiber() {
	app := fiber.New()
	db.ConnectDB()
	defer db.CloseDB()

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
