package controller

import (
	"GoBalance/app_server/workers"
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"time"
)

// Hello handler implementation application server
// Response :
//
//	200, {"message" : "hello-world"}
//	500, Internal Server Error
func Hello(w http.ResponseWriter, r *http.Request) {
	workers.Wrkr.Logger.Printf("Received request on hello route")

	// Simulate a delay
	delay := workers.Wrkr.Delay()

	// Determine the response(Success or Failure)
	success := float64(rand.IntN(10))/10.00 >= workers.Wrkr.FailurePercent

	// Augemented delay to match the actual response time
	time.Sleep(delay - 150)

	// Update the stats
	workers.Wrkr.UpdateStats(success, delay)

	// Return a response
	if !success {
		workers.Wrkr.Logger.Printf("Request failed (simulated failure)")
		http.Error(w, "(Simulated) Internal Server Error", http.StatusInternalServerError)
		return
	}

	workers.Wrkr.Logger.Printf("Request successful, responded after %v", delay)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "hello-world"})
}
