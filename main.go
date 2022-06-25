package main

import (
	"fmt"
	"log"
	"os"
	"switchboard/internal/consumer_api"
	"switchboard/internal/db"
	"switchboard/internal/management_api"
	"sync"

	"github.com/gin-gonic/gin"
	env "github.com/joho/godotenv"
)

type RouterFactory func(string) *gin.Engine

type serverConfig struct {
	createRouterFn RouterFactory
	servName       string
	port           string
}

func startServer(sc *serverConfig) error {
	r := sc.createRouterFn(sc.servName)
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

	var wg sync.WaitGroup
	wg.Add(len(servers))

	for _, s := range servers {
		go func(sc serverConfig) {
			startServer(&sc)
			wg.Done()
		}(s)
	}

	wg.Wait()
}
