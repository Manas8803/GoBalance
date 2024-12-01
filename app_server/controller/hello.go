package controller

import (
	"GoBalance/app_server/workers"
	"encoding/json"
	"math/rand/v2"
	"net/http"
)

// Hello handler implementation application server
// Response :
//
//	200, {"message" : "hello-world"}
//	500, Internal Server Error
func Hello(w http.ResponseWriter, r *http.Request) {
	workers.Wrkr.Logger.Printf("Received request on hello route")

	// Determine the response(Success or Failure)
	success := float64(rand.IntN(10))/10.00 >= workers.Wrkr.FailurePercent

	// Update the stats
	workers.Wrkr.UpdateStats(success)

	// Return a response
	if !success {
		workers.Wrkr.Logger.Printf("Request failed (simulated failure)")
		http.Error(w, "(Simulated) Internal Server Error", http.StatusInternalServerError)
		return
	}

	workers.Wrkr.Logger.Printf("Request successful")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "hello-world"})
}
