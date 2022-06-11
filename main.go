package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"switchboard/internal/common"
	"switchboard/internal/common/auth"
	"switchboard/internal/common/db"
	"switchboard/internal/common/middleware"

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

	r.POST("/auth/login", auth.LoginRoute)
	r.POST("/auth/signup", auth.SignUpRoute)
	r.POST("/auth/logout", auth.LogOutRoute)

	r.Use(auth.ParseAuthToken())
	r.Use(auth.RequireAuthentication())

	// temporary endpoints to test random data generator
	r.POST("/randomjson", func(c *gin.Context) {
		jsonData, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Status(500)
		}
		c.Header("Content-Type", "application/json")

		c.Stream(func(w io.Writer) bool {
			err := common.GenFakeJson(string(jsonData), w)
			if err != nil {
				log.Println(err)
			}
			return false
		})
	})

	log.Println("Starting server...")
	log.Fatal(r.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
