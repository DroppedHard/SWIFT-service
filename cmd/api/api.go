package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DroppedHard/SWIFT-service/service/swiftCode"
	"github.com/DroppedHard/SWIFT-service/utils"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type APIServer struct {
	host   string
	port   string
	client *redis.Client
}

func NewAPIServer(host string, port string, client *redis.Client) *APIServer {
	return &APIServer{
		host:   host,
		port:   port,
		client: client,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix(utils.ApiPrefix).Subrouter()

	bankDataStore := swiftCode.NewStore(s.client)
	userHandler := swiftCode.NewHandler(bankDataStore)
	userHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", fmt.Sprintf("%s:%s", s.host, s.port))

	return http.ListenAndServe(fmt.Sprintf(":%s", s.port), router)
}
