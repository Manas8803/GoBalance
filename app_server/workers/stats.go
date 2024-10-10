package workers

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Function to read worker_stats.json file and return the current stats of the worker node.
func GetStats() (*Stats, error) {
	statsFile := filepath.Join(Wrkr.StatsDir, "worker_stats.json")

	statsJSON, err := os.ReadFile(statsFile)
	if err != nil {
		if os.IsNotExist(err) {
			defaultStats := &Stats{
				SuccessfulRequests: 0,
				FailedRequests:     0,
				TotalRequests:      0,
				AvgDelayTime:       0,
			}
			defaultStatsJSON, err := json.Marshal(defaultStats)
			if err != nil {
				Wrkr.Logger.Println("Error marshalling default stats data", err)
				return nil, err
			}
			err = os.WriteFile(statsFile, defaultStatsJSON, 0644)
			if err != nil {
				Wrkr.Logger.Println("Error creating worker_stats.json with default values", err)
				return nil, err
			}

			return defaultStats, nil
		}

		Wrkr.Logger.Println("Error reading stats file", err)
		return nil, err
	}

	var stats Stats
	err = json.Unmarshal(statsJSON, &stats)
	if err != nil {
		Wrkr.Logger.Println("Error unmarshalling JSON (stats) data", err)
		return nil, err
	}

	return &stats, nil
}
