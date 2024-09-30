package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        int    `json:"_id bson:"_id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

var collection *mongo.Collection

func main(){
	fmt.Println("API without DB")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading", err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Atlas")

	collection = client.Database("golang_db").Collection("todos")
	app := fiber.New()

	app.Get("/api/todos", getTodos)
	// app.Post("/api/todos", createTodos)
	// app.Patch("/api/todos/:id", updateTodos)
	// app.Delete("/api/todos/:id", deleteTodos)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "4000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + PORT))
}


func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	//search cursor in context of MONGODB
	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}
	
	//defer is a keyword we use to postpone the execution of a function.
	//it is an advancement feature
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