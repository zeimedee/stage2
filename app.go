package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/zeimedee/stage2/mailer"
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
	Email := os.Getenv("EMAIL")
	Pass := os.Getenv("PASSWORD")

	mes := new(Mes)
	err := c.BodyParser(mes)
	mailer.Check(err)

	email := Email
	password := Pass
	recipient := mes.Email
	cc := []string{}
	path, err := os.Getwd()
	mailer.Check(err)

	sender := mailer.NewSender(email, password)

	msg, err := sender.WriteMessage(mes.Name, path+"mailTemplate/mailTemplate.html")
	if err != nil {
		fmt.Println(err)
	}

	body := sender.WriteEmail(recipient, mes.Subject, msg, cc)

	mail, err := sender.Mail(mes.Name, string(body), recipient)
	mailer.Check(err)

	return c.Render("index", fiber.Map{
		"Sub":  mes.Name,
		"Sent": mail,
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
