package main

import (
	"fmt"
	"os"
	"switchboard/internal/consumer_api"
	"switchboard/internal/db"
	"switchboard/internal/management_api"

	"github.com/gin-gonic/gin"
	env "github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

const (
	CONSUMER_API_NAME   = "CONSUMER_API"
	MANAGEMENT_API_NAME = "MANAGEMENT_API"
)

type serverConfig struct {
	name   string
	router *gin.Engine
	port   string
}

func startConsumerAPI(reload chan bool, quit chan<- bool) {
	port := os.Getenv("CONSUMER_API_PORT")
	c_api := consumer_api.New(CONSUMER_API_NAME)
	c_api.InitialiseRouter(reload, quit)
	fmt.Printf("%s is starting on port %s...\n", CONSUMER_API_NAME, port)
	if err := c_api.Router.Run(fmt.Sprintf(":%s", port)); err != nil {
		fmt.Println("Could not start consumer API", err)
	}
}

func startManagementAPI(reload chan bool, quit chan<- bool) {
	port := os.Getenv("MANAGEMENT_API_PORT")
	m_api := management_api.CreateRouter(MANAGEMENT_API_NAME, reload, quit)
	fmt.Printf("%s is starting on port %s...\n", MANAGEMENT_API_NAME, port)
	if err := m_api.Run(fmt.Sprintf(":%s", port)); err != nil {
		fmt.Println("Could not start management API", err)
	}
}

func main() {
	err := env.Load()
	if err != nil {
		log.Println("could not locate or read .env file", err)
	}

	db.Setup()
	gin.SetMode(gin.ReleaseMode)

	reload := make(chan bool)
	quit := make(chan bool)
	go startConsumerAPI(reload, quit)
	go startManagementAPI(reload, quit)

	if <-quit {
		log.Fatalln("an error occurred while starting the servers")
	}
}
