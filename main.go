package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"store/pkg/redis"
	httpServer "store/pkg/server"
	"store/pkg/utils"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	redisAddr := utils.GetEnv("REDIS_ADDR", "localhost:6379")
	port := utils.GetEnv("STORE_PORT", "8080")

	redisPassword := utils.GetEnv("REDIS_PASS", "")

	redisTimeout, err := strconv.Atoi(utils.GetEnv("REDIS_TIMEOUT", "10"))
	if err != nil {
		// Fail fast
		log.Fatalf("Can not convert TIMEOUT env var to int")
	}
	redisPoolSize, err := strconv.Atoi(utils.GetEnv("REDIS_POOLSIZE", "100"))
	if err != nil {
		// Fail fast
		log.Fatalf("Can not convert TIMEOUT env var to int")
	}

	timeout, err := strconv.Atoi(utils.GetEnv("TIMEOUT", "10"))

	if err != nil {
		// Fail fast
		log.Fatalf("Can not convert TIMEOUT env var to int")
	}

	repo := redis.NewOrdersRepo(redisAddr,
		redisPassword,
		redisPoolSize,
		redisTimeout,
	)
	server := httpServer.NewServer(repo)
	mux := http.NewServeMux()
	setupRoutes(mux, server)

	fmt.Println("start")

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * time.Duration(timeout),
		ReadTimeout:  time.Second * time.Duration(timeout),
		IdleTimeout:  time.Second * time.Duration(timeout),
		Handler:      mux,
	}

	// Run our http server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout))
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	repo.Shutdown(ctx)

	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

}

func setupRoutes(mux *http.ServeMux, server *httpServer.Server) {
	mux.HandleFunc("/create", server.Create)
	mux.HandleFunc("/order/", server.Get)
	mux.Handle("/metrics", promhttp.Handler())
}
