package main

import (
	"GoBalance/loadbalancer/controllers"
	"GoBalance/loadbalancer/lb"
	"GoBalance/loadbalancer/lib/middleware"
	"log"
	"net/http"
)

// Initialize the load balancer
func init() {
	retries := 2
	var err error
	for i := 0; i < retries; i++ {
		err = lb.Init()
		if err == nil {
			break
		}
		log.Printf("Error initializing load balancer (attempt %d/%d): %v\n", i+1, retries, err)
	}
	if err != nil {
		log.Fatal("Failed to initialize load balancer after all retries")
	}
}

func main() {
	if lb.LB == nil {
		log.Println("Unable to intialize the load balancer")
		return
	}
	// Setup the routes with middleware
	http.HandleFunc("/api/v1/hello", middleware.ScalingMiddleware(controllers.Hello))
	http.HandleFunc("/worker/stats", controllers.Stats)

	// Start the server
	lb.LB.Logger.Println("Load Balancer started on :2000")
	err := http.ListenAndServe(":2000", nil)
	if err != nil {
		lb.LB.Logger.Fatal("Error starting server: ", err)
	}
}
