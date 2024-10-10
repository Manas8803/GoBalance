package main

import (
	"GoBalance/app_server/controller"
	"GoBalance/app_server/workers"
	"log"
	"net/http"
)

func main() {
	// Intialize the worker node
	retries := 2
	err := workers.Init()
	for i := 0; i < retries && err != nil; i++ {
		log.Println("Error initializing worker node : ", err)
		log.Println("Retrying...")
	}

	if err != nil {
		log.Println("Failed to initialize worker node : ", err)
		return
	}
	// Setup routes
	http.HandleFunc("/api/v1/hello", controller.Hello)
	http.HandleFunc("/worker/stats", controller.Stats)
	http.HandleFunc("/ping", controller.Ping)

	// Start the server
	workers.Wrkr.Logger.Println("Server is running on :8080")
	workers.Wrkr.Logger.Fatal(http.ListenAndServe(":8080", nil))
}
