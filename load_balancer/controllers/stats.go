package controllers

import (
	"GoBalance/loadbalancer/lb"
	"GoBalance/loadbalancer/lib/file"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Stats handler for forwarding requests on /worker/stats route
func Stats(w http.ResponseWriter, r *http.Request) {
	stats := make(map[string]interface{})
	totalStats := lb.WorkerStats{}

	// Read IP addresses from all_nodes.txt
	ipAddresses, err := file.ReadIPAddresses("all_nodes.txt")
	if err != nil {
		http.Error(w, "Error reading IP addresses: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup
	statsChan := make(chan map[string]lb.WorkerStats, len(ipAddresses))

	// Fetching stats for each worker node using go routines
	for i, ipAddress := range ipAddresses {
		if strings.TrimSpace(ipAddress) == "" {
			continue
		}
		wg.Add(1)
		go func(i int, ipAddress string) {
			defer wg.Done()
			// Find the corresponding worker in lb.LB.Workers
			var worker *lb.Worker
			for _, w := range lb.LB.Workers {
				if w.URL.String() == fmt.Sprintf("http://%s:8080", ipAddress) {
					worker = w
					break
				}
			}
			if worker == nil {
				// If no matching worker is found, create a dummy worker with the parsed URL
				parsedURL, err := url.Parse(fmt.Sprintf("http://%s:8080", ipAddress))
				if err != nil {
					// Handle the URL parsing error (log or return)
					statsChan <- map[string]lb.WorkerStats{fmt.Sprintf("worker%d", i+1): {}}
					return
				}
				worker = &lb.Worker{URL: parsedURL}
			}
			workerStats := lb.FetchWorkerStats(worker)
			statsChan <- map[string]lb.WorkerStats{fmt.Sprintf("worker%d", i+1): workerStats}
		}(i, ipAddress)
	}

	go func() {
		wg.Wait()
		close(statsChan)
	}()

	// Calculating overall number of requests(success, fail, total) across all worker nodes
	for workerStat := range statsChan {
		for workerName, stat := range workerStat {
			stats[workerName] = stat
			totalStats.SuccessfulRequests += stat.SuccessfulRequests
			totalStats.FailedRequests += stat.FailedRequests
			totalStats.TotalRequests += stat.TotalRequests

			// Parse and add average delay time
			workerAvgDelay, err := time.ParseDuration(stat.AvgDelayTime)
			if err == nil {
				totalAvgDelay, _ := time.ParseDuration(totalStats.AvgDelayTime)
				totalAvgDelay += workerAvgDelay
				totalStats.AvgDelayTime = totalAvgDelay.String()
			}
		}
	}

	// Calculate average delay time across all workers
	if len(lb.LB.Workers) > 0 {
		totalAvgDelay, _ := time.ParseDuration(totalStats.AvgDelayTime)
		totalStats.AvgDelayTime = (totalAvgDelay / time.Duration(len(lb.LB.Workers))).String()
	}

	result := map[string]interface{}{
		"success-request":  map[string]int{"total": totalStats.SuccessfulRequests},
		"failed-request":   map[string]int{"total": totalStats.FailedRequests},
		"total-request":    map[string]int{"total": totalStats.TotalRequests},
		"avg-request-time": map[string]string{"total": totalStats.AvgDelayTime},
	}

	// Preparing the response
	for workerName, stat := range stats {
		workerStat := stat.(lb.WorkerStats)
		result["success-request"].(map[string]int)[workerName] = workerStat.SuccessfulRequests
		result["failed-request"].(map[string]int)[workerName] = workerStat.FailedRequests
		result["total-request"].(map[string]int)[workerName] = workerStat.TotalRequests
		result["avg-request-time"].(map[string]string)[workerName] = workerStat.AvgDelayTime
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
