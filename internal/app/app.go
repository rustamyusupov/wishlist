package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/sirupsen/logrus"
)

type Wish struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Link     string `json:"link"`
	Price    string `json:"price"`
	Currency string `json:"currency"`
}

var Wishes = []Wish{
	{
		ID:       "1",
		Name:     "Apple Macbook Pro 13",
		Link:     "https://www.apple.com/macbook-pro-13/",
		Price:    "2000",
		Currency: "$",
	},
	{
		ID:       "2",
		Name:     "Apple iPhone 12",
		Link:     "https://www.apple.com/iphone-12/",
		Price:    "1000",
		Currency: "$",
	},
	{
		ID:       "3",
		Name:     "Apple Watch Series 6",
		Link:     "https://www.apple.com/apple-watch-series-6/",
		Price:    "500",
		Currency: "$",
	},
	{
		ID:       "4",
		Name:     "Jack and Jones T-shirt with a cat",
		Link:     "https://www.jackjones.com/",
		Price:    "50",
		Currency: "$",
	},
}

func Run() {
	engine := html.New("./web", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":  "Wishlist",
			"Wishes": Wishes,
		})
	})

	app.Static("/", "./web")

	logrus.Fatal(app.Listen(":3000"))
}
