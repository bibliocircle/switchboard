package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"switchboard/internal/common/auth"
	"switchboard/internal/common/db"
	"switchboard/internal/common/middleware"
	"switchboard/internal/common/randomdata"
	"switchboard/internal/crud/endpoint"
	"switchboard/internal/crud/mockservice"
	"switchboard/internal/crud/scenario"
	"switchboard/internal/crud/upstream"

	"github.com/gin-gonic/gin"
	env "github.com/joho/godotenv"
)

func setUpDatabase() {
	ctx := context.Background()
	dbError := db.Connect(ctx)
	if dbError != nil {
		log.Fatalln("could not connect to the database", dbError)
	}
	db.Migrate(db.Client)
}

func setupUnauthenticatedRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.Writer.Write([]byte("pong"))
	})

	r.POST("/auth/login", auth.LoginRoute)
	r.POST("/auth/signup", auth.SignUpRoute)
	r.POST("/auth/logout", auth.LogOutRoute)
}

func setupAuthenticatedRoutes(r *gin.Engine) {
	r.POST("/endpoint", endpoint.CreateEndpointRoute)
	r.POST("/scenario", scenario.CreateScenarioRoute)
	r.POST("/upstream", upstream.CreateUpstreamRoute)
	r.POST("/mockservice", mockservice.CreateMockServiceRoute)

	r.GET("/protected", func(c *gin.Context) {
		user := c.Value("user").(*auth.User)
		c.JSON(http.StatusOK, user)
	})

	// temporary endpoints to test random data generator
	r.POST("/randomjson", func(c *gin.Context) {
		jsonData, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Status(500)
		}
		c.Header("Content-Type", "application/json")

		c.Stream(func(w io.Writer) bool {
			err := randomdata.GenFakeJson(string(jsonData), w)
			if err != nil {
				log.Println(err)
			}
			return false
		})
	})
}

func main() {
	err := env.Load()
	if err != nil {
		log.Println("could not locate or read .env file", err)
	}

	setUpDatabase()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.ConfigureCors())

	setupUnauthenticatedRoutes(r)

	r.Use(auth.ParseAuthToken())
	r.Use(auth.RequireAuthentication())

	setupAuthenticatedRoutes(r)

	log.Println("Starting server...")
	log.Fatal(r.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
