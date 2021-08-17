package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"log"
	"os"
)

type Mes struct {
	Name    string `json:"name" form:"name"`
	Email   string `json:"email" form:"email"`
	Subject string `json:"subject" form:"subject"`
	Message string `json:"message" form:"message"`
}

func hi(c *fiber.Ctx) error {
	return c.SendString("hi")
}

func hello(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title":   "hello world",
		"Message": "alexander",
	})
}

func sub(c *fiber.Ctx) error {
	mes := new(Mes)
	err := c.BodyParser(mes)
	if err != nil {
		log.Fatal(err)
	}
	return c.Render("index", fiber.Map{
		"Sub": mes.Name,
	})
}

func handlers(app *fiber.App) {
	app.Get("/hi", hi)

	app.Get("/", hello)

	app.Post("/", sub)
}

func main() {
	port := os.Getenv("PORT")
	engine := html.New("./public", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	handlers(app)

	log.Fatal(app.Listen(":" + port))
}
