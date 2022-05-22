package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"switchboard/db"
	"switchboard/middleware"

	"github.com/gin-gonic/gin"
	env "github.com/joho/godotenv"
)

func main() {
	err := env.Load()
	if err != nil {
		log.Println("could not locate or read .env file", err)
	}
	ctx := context.Background()
	dbError := db.Connect(ctx)
	if dbError != nil {
		log.Fatalln("could not connect to the database", dbError)
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.ConfigureCors())

	r.GET("/ping", func(c *gin.Context) {
		c.Writer.Write([]byte("pong"))
	})

	log.Println("Starting server...")
	log.Fatal(r.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
