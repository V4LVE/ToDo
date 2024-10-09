package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

// Structs are like objects in c#
type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {

	fmt.Println("Golang backend is live")
	app := fiber.New()

	todos := []Todo{}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "Welcome to the react-go todo list"})
	})

	app.Post("/api/todo/create", func(c *fiber.Ctx) error {
		todo := new(Todo)

		// Check for errors doing the user parse This binds the JSON data youre sending with the request
		if err := c.BodyParser(todo); err != nil {
			return err
		}

		// check for null or empty in Todo Body
		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(201).JSON(todo)
	})

	app.Patch("/api/todo/markdone/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "The todo with id: " + id + " was not found"})
	})

	log.Fatal(app.Listen(":4000"))
}
