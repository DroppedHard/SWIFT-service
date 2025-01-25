package api

import (
	"log"
	"net/http"

	"github.com/DroppedHard/SWIFT-service/service/swiftCode"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type APIServer struct {
	addr   string
	client *redis.Client
}

func NewAPIServer(addr string, client *redis.Client) *APIServer {
	return &APIServer{
		addr:   addr,
		client: client,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	bankDataStore := swiftCode.NewStore(s.client)
	userHandler := swiftCode.NewHandler(bankDataStore)
	userHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
