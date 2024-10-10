package controllers

import (
	"GoBalance/loadbalancer/lb"
	"net/http"
	"time"
)

// Hello handler for forwarding requests on api/v1/hello route
func Hello(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the next worker node from the available pool of nodes
	worker := lb.LB.NextWorker()
	if worker == nil {
		lb.LB.Logger.Println("No available workers")
		http.Error(w, "No available workers", http.StatusServiceUnavailable)
		return
	}

	url := worker.URL.String()

	// Ping the worker node to check if it is healthy
	resp, err := http.Get(url + "/ping")
	if err != nil {
		lb.LB.Logger.Println(err)
		http.Error(w, "Unable to perform ping : "+err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		lb.LB.Logger.Println("Health check unsuccessful for ", url)
		// We can implement a logic to remove the worker node from the pool here.
		http.Error(w, "Unable to reach server", http.StatusServiceUnavailable)
		return
	}

	lb.LB.Logger.Printf("Worker at %s passed health check after %dms", url, time.Since(startTime).Milliseconds())

	worker.ReverseProxy.ServeHTTP(w, r)
}
