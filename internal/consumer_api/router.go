package consumer_api

import (
	"fmt"
	"switchboard/internal/common"
	"switchboard/internal/db"
	"switchboard/internal/models"
	"sync"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type RouteConfig struct {
	MockService models.MockService
	Endpoint    models.Endpoint
}

func launchEndpoints(services <-chan models.MockService, quit chan<- bool) chan RouteConfig {
	rc := make(chan RouteConfig)
	go func() {
		var wg sync.WaitGroup
		for s := range services {
			wg.Add(1)
			go func(svc models.MockService) {
				endpoints, err := db.GetEndpoints(svc.ID)
				if err != nil {
					log.Errorln(fmt.Sprintf("could not launch endpoints for mock service: %s (ID: %s)", svc.Name, svc.ID), err)
					quit <- true
				}
				for _, e := range endpoints {
					rc <- RouteConfig{
						MockService: svc,
						Endpoint:    e,
					}
				}
				wg.Done()
			}(s)
		}
		wg.Wait()
		close(rc)
	}()
	return rc
}

func launchServices(quit chan<- bool) chan models.MockService {
	services := make(chan models.MockService)
	go func() {
		mslist, err := db.GetMockServices()
		if err != nil {
			log.Errorln("could not retrieve mock services", err)
			quit <- true
		}
		for _, s := range mslist {
			services <- s
		}
		close(services)
	}()
	return services
}

func CreateRouter(name string, quit chan<- bool) *gin.Engine {
	r := gin.New()
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: common.CreateGinLogFormatter(name),
	}))
	r.Use(gin.Recovery())

	services := launchServices(quit)
	routeConfigs := launchEndpoints(services, quit)

	for rc := range routeConfigs {
		go InitRoute(r, rc)
	}
	return r
}
