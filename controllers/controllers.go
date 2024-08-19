package controllers

import (
	"wishes/models"

	"github.com/gofiber/fiber/v2"
)

func Read(c *fiber.Ctx) error {
	categories := models.GetCategories()
	return c.Render("index", fiber.Map{
		"Title":      "Wishlist",
		"Categories": categories,
	})
}

func Create(c *fiber.Ctx) error {
	return c.SendString("Create")
}

func ReadById(c *fiber.Ctx) error {
	return c.SendString("ReadById")
}

func Update(c *fiber.Ctx) error {
	return c.SendString("Update")
}

func Delete(c *fiber.Ctx) error {
	return c.SendString("Delete")
}
