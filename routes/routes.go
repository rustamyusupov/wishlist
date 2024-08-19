package routes

import (
	"wishes/controllers"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	app.Post("/", controllers.Create)
	app.Get("/", controllers.Read)
	app.Get("/:id", controllers.ReadById)
	app.Patch("/:id", controllers.Update)
	app.Delete("/:id", controllers.Delete)
}
