package lb

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Worker struct {
	URL          *url.URL
	ReverseProxy *httputil.ReverseProxy
}

type WorkerStats struct {
	SuccessfulRequests int    `json:"success_requests"`
	FailedRequests     int    `json:"failed_requests"`
	TotalRequests      int    `json:"total_requests"`
}

// Function to fetch worker stats from a given worker node
func FetchWorkerStats(worker *Worker) WorkerStats {
	resp, err := http.Get(worker.URL.String() + "/worker/stats")
	if err != nil {
		LB.Logger.Printf("Error fetching stats from worker %s: %v", worker.URL.String(), err)
		return WorkerStats{}
	}
	defer resp.Body.Close()

	var stats WorkerStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		LB.Logger.Printf("Error decoding stats from worker %s: %v", worker.URL.String(), err)
		return WorkerStats{}
	}

	return stats
}
