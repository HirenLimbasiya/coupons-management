package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	app := fiber.New()

	mongoEndpoint := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(mongoEndpoint)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
		os.Exit(1)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("MongoDB connection test failed:", err)
		os.Exit(1)
	}
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber with MongoDB!")
	})

	app.Listen(":2121")
}
