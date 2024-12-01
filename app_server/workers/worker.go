package workers

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type Worker struct {
	ID             string
	FailurePercent float64
	Stats          *Stats
	StatsDir       string
	Logger         *log.Logger
	Rng            *rand.Rand
}

type Stats struct {
	SuccessfulRequests int `json:"successful_requests"`
	FailedRequests     int `json:"failed_requests"`
	TotalRequests      int `json:"total_requests"`
	mutex              sync.Mutex
}

var Wrkr *Worker

func NewWorkerNode(worker_id string, failure_percent float64, logger *log.Logger, rng *rand.Rand) *Worker {
	return &Worker{
		ID:             worker_id,
		FailurePercent: failure_percent / 100,
		Stats:          &Stats{},
		StatsDir:       os.Getenv("WORKER_DIR"),
		Logger:         logger,
		Rng:            rng,
	}
}

// Initialize Worker Node
func Init() error {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println("Error loading environment\nContinuing..")
	}

	writer := io.Writer(os.Stdout)

	worker_id := os.Getenv("WORKER_ID")
	log.Println(worker_id)

	logger := log.New(writer, "WORKER - "+worker_id+" : ", log.Ldate|log.Ltime|log.Lshortfile)

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	failurePercentStr := os.Getenv("FAIL_PERCENT")

	failure_percent, err := strconv.ParseFloat(failurePercentStr, 64)
	if err != nil {
		logger.Println("Error Parsing FAIL_PERCENT: ", err)
		logger.Println("Using default FAIL_PERCENT of 10%")
		failure_percent = 10.00
	}

	Wrkr = NewWorkerNode(worker_id, failure_percent, logger, rng)

	if err := os.MkdirAll(Wrkr.StatsDir, 0755); err != nil {
		Wrkr.Logger.Fatalf("Failed to create stats directory: %v", err)
		return err
	}

	Wrkr.Logger.Printf("Worker %s initialized successfully", worker_id)
	return nil
}

// Synchronized function to current update statistics of worker node and write it to the stats file
func (w *Worker) UpdateStats(success bool) {

	// Acquire and release lock for accessing stats
	w.Stats.mutex.Lock()
	defer w.Stats.mutex.Unlock()

	w.Stats.TotalRequests++

	if success {
		w.Stats.SuccessfulRequests++
	} else {
		w.Stats.FailedRequests++
	}

	// Write stats to file
	statsFile := filepath.Join(w.StatsDir, "worker_stats.json")
	statsJSON, _ := json.Marshal(w.Stats)
	if err := os.WriteFile(statsFile, statsJSON, 0644); err != nil {
		w.Logger.Printf("Failed to write stats to file: %v", err)
	}
}
