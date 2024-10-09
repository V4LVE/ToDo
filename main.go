package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
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

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:" + err.Error())
	}

	port := os.Getenv("PORT")

	todos := []Todo{}

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	// Creates a todo item
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

	// Marks todo item as completed
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

	//Delete a todo item
	app.Delete("/api/todo/delete/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
			}

		}

		return c.Status(200).JSON(fiber.Map{"success": true})
	})

	log.Fatal(app.Listen(":" + port))
}
