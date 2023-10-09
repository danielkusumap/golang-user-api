package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	URI := fmt.Sprintf("mongodb+srv://%s:%s@%s.zevxmja.mongodb.net/", dbUser, dbPass, dbName)
	clientOptions := options.Client().ApplyURI(URI)

	var err error

	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to MongoDB!")
}

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	router.Use(cors.New(config))
	
	routes.ConfigureUserRoutes(router, client)
	// config.AllowOrigins = []string{"http://localhost:3000"} // Replace with your frontend domain
	// router.GET("/users", getUsers)
	// router.POST("/api/register", registerUser)
	// router.POST("/api/login", loginUser)
	// router.PUT("/users", updateUser)

	router.Run(":8080")
}
