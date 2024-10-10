package controller

import (
	"GoBalance/app_server/workers"
	"encoding/json"
	"net/http"
)

// Ping handler implementation application server
// Response :
//
//	200, {"message" : "pong"}
func Ping(w http.ResponseWriter, r *http.Request) {
	workers.Wrkr.Logger.Printf("Request received on health check route")

	// Add critical health checks of different resources like db connections
	// Meant for future improvements
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
}
