package main

import (
	"fmt"
	"log"
	"net/http"
	"store/pkg/redis"
	httpServer "store/pkg/server"
	"store/pkg/utils"
)

func main() {
	reddisAddr := utils.GetEnv("REDIS_ADDR", "localhost:6379")
	port := utils.GetEnv("STORE_PORT", ":8080")

	repo := redis.NewOrdersRepo(reddisAddr)
	server := httpServer.NewServer(repo)
	mux := http.NewServeMux()
	setupRoutes(mux, server)

	fmt.Println("start")
	log.Fatal(http.ListenAndServe(port, mux))
}

func setupRoutes(mux *http.ServeMux, server *httpServer.Server) {
	mux.HandleFunc("/create", server.Create)
	mux.HandleFunc("/order/", server.Get)
}
