package consumer_api

import (
	"fmt"
	"net/http"
	"switchboard/internal/common"
	"switchboard/internal/db"
	"switchboard/internal/models"
	"sync"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type HandlerConfig struct {
	Method      string
	Path        string
	HandlerFunc gin.HandlerFunc
}

type ServiceRegistryEntry struct {
	MockService *models.MockService
	Handlers    []HandlerConfig
}

type ConsumerApiRouter struct {
	name            string
	serviceRegistry map[string]*gin.Engine
	Router          *gin.Engine
}

func New(name string) ConsumerApiRouter {
	r := ConsumerApiRouter{
		name:            name,
		serviceRegistry: map[string]*gin.Engine{},
		Router:          gin.New(),
	}

	r.Router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: common.CreateGinLogFormatter(name),
	}))
	r.Router.Use(gin.Recovery())
	return r
}

func (r ConsumerApiRouter) initEndpoints(services <-chan *models.MockService, quit chan<- bool) chan ServiceRegistryEntry {
	se := make(chan ServiceRegistryEntry)
	go func() {
		var wg sync.WaitGroup
		for s := range services {
			wg.Add(1)
			go func(svc *models.MockService) {
				// may need to initialise middleware here
				endpoints, err := db.GetEndpoints(svc.ID)
				if err != nil {
					log.Errorln(fmt.Sprintf("could not launch endpoints for mock service: %s (ID: %s)", svc.Name, svc.ID), err)
					quit <- true
				}
				handlers := []HandlerConfig{}
				for _, e := range endpoints {
					handlers = append(handlers, HandlerConfig{
						Method:      e.Method,
						Path:        fmt.Sprintf("/ws/:workspaceId/%s%s", svc.ID, e.Path),
						HandlerFunc: CreateRoute(svc.ID, e.ID),
					})
				}
				se <- ServiceRegistryEntry{
					MockService: svc,
					Handlers:    handlers,
				}
				wg.Done()
			}(s)
		}
		wg.Wait()
		close(se)
	}()
	return se
}

func (r ConsumerApiRouter) initServices(quit chan<- bool) chan *models.MockService {
	services := make(chan *models.MockService)
	go func() {
		mslist, err := db.GetMockServices()
		if err != nil {
			log.Errorln("could not retrieve mock services", err)
			quit <- true
		}
		for _, s := range *mslist {
			services <- s
		}
		close(services)
	}()
	return services
}

func (r ConsumerApiRouter) initServiceRegistry(entries <-chan ServiceRegistryEntry) {
	msRouter := gin.New()
	for entry := range entries {
		for _, h := range entry.Handlers {
			msRouter.Handle(h.Method, h.Path, h.HandlerFunc)
		}
		r.serviceRegistry[entry.MockService.ID] = msRouter
	}
}

func (r ConsumerApiRouter) InitialiseRouter(reload chan bool, quit chan<- bool) {
	services := r.initServices(quit)
	entries := r.initEndpoints(services, quit)
	r.initServiceRegistry(entries)

	r.Router.Any("/ws/:workspaceId/:mockServiceId/*path", func(c *gin.Context) {
		msID := c.Param("mockServiceId")
		msRouter := r.serviceRegistry[msID]

		if msRouter != nil {
			msRouter.ServeHTTP(c.Writer, c.Request)
			return
		}
		c.Writer.WriteHeader(http.StatusTeapot)
	})
}
