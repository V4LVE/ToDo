package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

var collection *mongo.Collection
var memoryTodos []Todo
var memoryMu sync.Mutex

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, falling back to environment variables")
	}
	fmt.Println("Starting server...")
	fmt.Println("Attempting to connect to MongoDB Atlas...")
	dbconstring := os.Getenv("dbconstring")
	if dbconstring != "" {
		clientOptions := options.Client().ApplyURI(dbconstring)
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Println("Database unavailable, using in-memory todo store:", err)
		} else if err := client.Ping(context.Background(), nil); err != nil {
			log.Println("Database unavailable, using in-memory todo store:", err)
			_ = client.Disconnect(context.Background())
		} else {
			defer client.Disconnect(context.Background())
			fmt.Println("Connected to MongoDB Atlas")
			collection = client.Database("golang_db").Collection("todos")
		}
	} else {
		log.Println("dbconstring not set, using in-memory todo store")
	}

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	app.Patch("/api/todo/markdone/:id", markDone)
	app.Patch("/api/todo/update/:id", updateTodo)
	app.Post("/api/todo/create", createTodo)
	app.Delete("/api/todo/delete/:id", deleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
	fmt.Println("We're live")
}

func getTodos(c *fiber.Ctx) error {
	if collection == nil {
		memoryMu.Lock()
		defer memoryMu.Unlock()

		return c.JSON(memoryTodos)
	}

	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}

	return c.JSON(todos)
}

func markDone(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	if collection == nil {
		memoryMu.Lock()
		defer memoryMu.Unlock()

		for index := range memoryTodos {
			if memoryTodos[index].ID == objectID {
				memoryTodos[index].Completed = true
				return c.Status(200).JSON(fiber.Map{"success": true})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "todo not found"})
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"completed": true}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"Error": "Body cannot be empty"})
	}

	if collection == nil {
		memoryMu.Lock()
		defer memoryMu.Unlock()

		for index := range memoryTodos {
			if memoryTodos[index].ID == objectID {
				memoryTodos[index].Body = todo.Body
				memoryTodos[index].Completed = todo.Completed
				return c.Status(200).JSON(memoryTodos[index])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "todo not found"})
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"body": todo.Body, "completed": todo.Completed}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return err
	}

	var updated Todo
	if err := collection.FindOne(context.Background(), filter).Decode(&updated); err != nil {
		return err
	}

	return c.Status(200).JSON(updated)
}

func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"Error": "Body cannot be empty"})
	}

	if collection == nil {
		todo.ID = primitive.NewObjectID()
		memoryMu.Lock()
		memoryTodos = append([]Todo{*todo}, memoryTodos...)
		memoryMu.Unlock()

		return c.Status(201).JSON(todo)
	}

	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)
}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	if collection == nil {
		memoryMu.Lock()
		defer memoryMu.Unlock()

		for index := range memoryTodos {
			if memoryTodos[index].ID == objectID {
				memoryTodos = append(memoryTodos[:index], memoryTodos[index+1:]...)
				return c.SendStatus(204)
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "todo not found"})
	}

	filter := bson.M{"_id": objectID}

	_, err = collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}

	return c.SendStatus(204)
}
