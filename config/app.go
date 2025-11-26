// package config

// import (
// 	"crud_alumni/route"

// 	"github.com/gofiber/fiber/v2"
// )

// func App() *fiber.App {
// 	app := fiber.New()

// 	// Middleware bisa ditambahkan di sini

// 	route.SetupRoutes(app)
// 	return app
// }

package config

import (
	"crud_alumni/route"

	"github.com/gofiber/fiber/v2"
)

func App() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "CRUD Alumni (MongoDB Version)",
	})
	route.SetupRoutes(app)
	return app
}
