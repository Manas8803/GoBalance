package controller

import (
	"GoBalance/app_server/workers"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Stats handler implementation application server
// Response :
//
//	200, {"success_requests","failed_requests", "total_requests","avg_delay_time" }
//	500, Failed to get stats
func Stats(w http.ResponseWriter, r *http.Request) {
	workers.Wrkr.Logger.Printf("Stats requested")

	stats, err := workers.GetStats()
	if err != nil {
		workers.Wrkr.Logger.Printf("Failed to get stats: %v", err)
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	avgDelayTimeMs := time.Duration(stats.AvgDelayTime).Milliseconds()
	response := map[string]interface{}{
		"success_requests": stats.SuccessfulRequests,
		"failed_requests":  stats.FailedRequests,
		"total_requests":   stats.TotalRequests,
		"avg_delay_time":   fmt.Sprintf("%dms", avgDelayTimeMs),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	workers.Wrkr.Logger.Println("Stats response sent successfully")
}
