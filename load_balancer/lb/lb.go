package lb

import (
	"fmt"
	"io"
	"log"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var LB *LoadBalancer

type LoadBalancer struct {
	Workers       []*Worker
	CurrentWorker int
	mux           sync.Mutex
	Logger        *log.Logger
}

func NewLoadBalancer(logger *log.Logger) *LoadBalancer {
	return &LoadBalancer{Logger: logger}
}

// Init initializes the load balancer
func Init() error {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println("Error loading environment\nContinuing..")
	}
	writer := io.Writer(os.Stdout)
	logger := log.New(writer, "LOADBALANCER : ", log.Ldate|log.Ltime|log.Lshortfile)
	LB = NewLoadBalancer(logger)
	nodesFile := filepath.Join("./", "available_nodes.txt")
	nodes_raw, err := os.ReadFile(nodesFile)
	if err != nil {
		LB.Logger.Fatal("Error reading available_nodes file : ", err)
		return err
	}

	lines := strings.Split(string(nodes_raw), "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" && isValidIPv4(trimmedLine) {
			err = LB.AddWorker(trimmedLine)
			if err != nil {
				LB.Logger.Println("Error adding worker node to LB pool: ", err)
			}
		}
	}
	LB.Logger.Println("LoadBalancer initialized successfully")
	return nil
}

// Method to add a worker node to the pool
func (lb *LoadBalancer) AddWorker(workerURL string) error {
	if workerURL == "" || strings.HasSuffix(workerURL, "\n") || strings.HasSuffix(workerURL, "\r\n") {
		return nil
	}
	parsedURL, err := parseWorkerURL(workerURL)
	if err != nil {
		return fmt.Errorf("invalid worker URL %s: %v", workerURL, err)
	}

	worker := &Worker{
		URL:          parsedURL,
		ReverseProxy: httputil.NewSingleHostReverseProxy(parsedURL),
	}

	lb.mux.Lock()
	lb.Workers = append(lb.Workers, worker)
	lb.mux.Unlock()

	lb.Logger.Printf("Added worker: %s\n", parsedURL)
	return nil
}

// Method to remove a worker node to the pool
func (lb *LoadBalancer) RemoveWorker(workerURL string) error {
	if workerURL == "" || strings.HasSuffix(workerURL, "\n") || strings.HasSuffix(workerURL, "\r\n") {
		return nil
	}
	parsedURL, err := parseWorkerURL(workerURL)
	if err != nil {
		return fmt.Errorf("invalid worker URL %s: %v", workerURL, err)
	}

	lb.mux.Lock()
	defer lb.mux.Unlock()

	for i, worker := range lb.Workers {
		if worker.URL.String() == parsedURL.String() {
			// Remove the worker
			lb.Workers = append(lb.Workers[:i], lb.Workers[i+1:]...)

			// Adjust CurrentWorker if necessary
			if lb.CurrentWorker >= len(lb.Workers) {
				lb.CurrentWorker = 0
			}

			lb.Logger.Printf("Removed worker: %s\n", parsedURL)
			return nil
		}
	}

	return fmt.Errorf("worker not found: %s", parsedURL)
}

// Method to determine the next worker node in the pool
func (lb *LoadBalancer) NextWorker() *Worker {
	lb.mux.Lock()
	defer lb.mux.Unlock()

	workerCount := len(lb.Workers)
	if workerCount == 0 {
		lb.Logger.Println("No workers available")
		return nil
	}

	// Ensure CurrentWorker is within bounds
	if lb.CurrentWorker >= workerCount {
		lb.CurrentWorker = 0
	}

	worker := lb.Workers[lb.CurrentWorker]
	lb.CurrentWorker = (lb.CurrentWorker + 1) % workerCount

	lb.Logger.Printf("Selected worker %d: %s\n", lb.CurrentWorker, worker.URL)
	return worker
}

// Function to parse and normalize worker URLs
func parseWorkerURL(workerURL string) (*url.URL, error) {
	if !strings.HasPrefix(workerURL, "http://") && !strings.HasPrefix(workerURL, "https://") {
		workerURL = "http://" + workerURL
	}

	parsedURL, err := url.Parse(workerURL)
	if err != nil {
		return nil, err
	}

	if parsedURL.Port() == "" {
		parsedURL.Host = fmt.Sprintf("%s:8080", parsedURL.Host)
	}

	return parsedURL, nil
}

// Function to check if a string is a valid IPv4 address
func isValidIPv4(ip string) bool {
	ipv4Regex := `^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	regex := regexp.MustCompile(ipv4Regex)
	return regex.MatchString(ip)
}
