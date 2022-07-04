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

type RouterFactory func(string, chan<- bool) *gin.Engine

type serverConfig struct {
	createRouterFn RouterFactory
	servName       string
	port           string
}

func startServer(sc *serverConfig, quit chan<- bool) error {
	r := sc.createRouterFn(sc.servName, quit)
	log.Printf("Starting server %s on port %s..\n", sc.servName, sc.port)
	return r.Run(fmt.Sprintf(":%s", sc.port))
}

func main() {
	err := env.Load()
	if err != nil {
		log.Println("could not locate or read .env file", err)
	}

	db.RunMigrations()
	gin.SetMode(gin.ReleaseMode)

	servers := []serverConfig{
		{
			createRouterFn: management_api.CreateRouter,
			servName:       "MANAGEMENT_API",
			port:           os.Getenv("MANAGEMENT_API_PORT"),
		},
		{
			createRouterFn: consumer_api.CreateRouter,
			servName:       "CONSUMER_API",
			port:           os.Getenv("CONSUMER_API_PORT"),
		},
	}

	quit := make(chan bool)
	for _, s := range servers {
		go func(sc serverConfig) {
			startServer(&sc, quit)
		}(s)
	}

	if <-quit {
		log.Fatalln("an error occurred while starting the servers")
	}
}
