package main

import (
	"context"
	"epl-fantasy/src/config"
	"epl-fantasy/src/db"
	"epl-fantasy/src/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func init() {
	config.App = config.GetConfig("LOCAL")
}

func main() {
	var err error
	config.Client, err = db.InitializeMongoDB(config.App)
	if err != nil {
		log.Fatalf("Error initializing MongoDB: %v", err)
	}
	log.Println("Successfully connected to MongoDB")

	r := mux.NewRouter()
	r.HandleFunc("/epl", handlers.FetchAndStoreGameWeekData).Methods("POST")
	r.HandleFunc("/epl", handlers.GetGameData).Methods("GET")
	r.HandleFunc("/epl/players", handlers.GetBestPerformers).Methods("GET")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Starting server on port 8080")
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown(): %v", err)
	} else {
		log.Println("Server stopped")
	}

}
