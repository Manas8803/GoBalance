package middleware

import (
	"GoBalance/loadbalancer/lb"
	"GoBalance/loadbalancer/lib/file"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var (
	maxConcurrentRequests int64
	minPoolSize           int64
	maxPoolSize           int64
	limiter               chan struct{}
	mu                    sync.Mutex
	currentRequests       int64
	fileMutex             sync.Mutex
)

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println("Error loading environment\nContinuing anyway..")
	}
	poolEnv := os.Getenv("POOL")
	if poolEnv == "" {
		maxConcurrentRequests = 20
	} else {
		parsed, err := strconv.ParseInt(poolEnv, 10, 64)
		if err != nil {
			log.Printf("Error parsing POOL environment variable: %v. Using default value of 20.", err)
			maxConcurrentRequests = 20
		} else {
			maxConcurrentRequests = parsed
		}
	}
	limiter = make(chan struct{}, maxConcurrentRequests)
	log.Printf("Max concurrent requests is set to: %d", maxConcurrentRequests)

	minPoolEnv := os.Getenv("WORKER")
	if minPoolEnv == "" {
		minPoolSize = 2
	} else {
		parsed, err := strconv.ParseInt(minPoolEnv, 10, 64)
		if err != nil {
			log.Printf("Error parsing WORKER environment variable: %v. Using default value of 2.", err)
			minPoolSize = 2
		} else {
			minPoolSize = parsed
		}
	}
	log.Printf("Min pool size set to: %d", minPoolSize)

	maxPoolEnv := os.Getenv("MAX_WORKER")
	if maxPoolEnv == "" {
		maxPoolSize = minPoolSize
	} else {
		parsed, err := strconv.ParseInt(maxPoolEnv, 10, 64)
		if err != nil {
			log.Printf("Error parsing WORKER environment variable: %v. Using default value of 2.", err)
			maxPoolSize = minPoolSize
		} else {
			maxPoolSize = parsed
		}
	}
	log.Printf("Max pool size set to: %d", maxPoolSize)
}

// Middleware that handles scaling based on the number of requests
func ScalingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case limiter <- struct{}{}:
			updateActiveRequests(1)

			// Check for scale-up logic
			checkScaleUp()

			next.ServeHTTP(w, r)
			defer func() {
				<-limiter
				updateActiveRequests(-1)

				// Check for scale-down logic
				checkScaleDown()
			}()

		default:
			// Too many requests, return 429 error
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
		}
	}
}

// Function to update the active request count
func updateActiveRequests(delta int64) {
	mu.Lock()
	defer mu.Unlock()
	currentRequests += delta
}

// Function to check if we need to scale up
func checkScaleUp() {
	mu.Lock()
	defer mu.Unlock()

	halfMax := maxConcurrentRequests / 2
	if currentRequests >= halfMax && int64(len(lb.LB.Workers)) < maxPoolSize {
		lb.LB.Logger.Printf("Scaling up, active requests: %d, current workers: %d", currentRequests, len(lb.LB.Workers))
		scaleUp()
	}
}

// Function to check if we need to scale down
func checkScaleDown() {
	mu.Lock()
	defer mu.Unlock()

	halfMax := maxConcurrentRequests / 2
	if currentRequests <= halfMax && int64(len(lb.LB.Workers)) > minPoolSize {
		lb.LB.Logger.Printf("Scaling down, active requests: %d, current workers: %d", currentRequests, len(lb.LB.Workers))
		scaleDown()
	}
}

// Function to add worker nodes from the pool of standy workers
func scaleUp() {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	ip, err := file.ReadFirstLineAndRemove("standby_nodes.txt")
	if err != nil {
		lb.LB.Logger.Printf("Error reading from standby_nodes.txt: %v", err)
		return
	}

	if ip == "" {
		lb.LB.Logger.Println("No standby nodes available for scaling up.")
		return
	}

	go func() {
		err := lb.LB.AddWorker(ip)
		if err != nil {
			lb.LB.Logger.Printf("Error adding worker %s: %v", ip, err)
			// If failed to add, put it back in standby
			file.AppendToFile("standby_nodes.txt", ip)
		} else {
			lb.LB.Logger.Printf("Successfully scaled up. Added worker: %s", ip)
			file.AppendToFile("available_nodes.txt", ip)
		}
	}()
}

// Function to remove worker nodes from the pool of available nodes
func scaleDown() {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	if int64(len(lb.LB.Workers)) <= minPoolSize {
		lb.LB.Logger.Printf("Cannot scale down. Current workers (%d) at or below min pool size (%d)", len(lb.LB.Workers), minPoolSize)
		return
	}

	ip, err := file.ReadLastLineAndRemove("available_nodes.txt")
	if err != nil {
		lb.LB.Logger.Printf("Error reading from available_nodes.txt: %v", err)
		return
	}

	if ip == "" {
		lb.LB.Logger.Println("No available nodes to scale down.")
		return
	}

	go func() {
		err := lb.LB.RemoveWorker(ip)
		if err != nil {
			lb.LB.Logger.Printf("Error removing worker %s: %v", ip, err)
			// If failed to remove, put it back in available
			file.AppendToFile("available_nodes.txt", ip)
		} else {
			lb.LB.Logger.Printf("Successfully scaled down. Removed worker: %s", ip)
			file.AppendToFile("standby_nodes.txt", ip)
		}
	}()
}
