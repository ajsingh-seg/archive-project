package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// JSONData defines the structure of the JSON data to be generated
type JSONData struct {
	Channel string `json:"channel"`
	Context struct {
		Library struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"library"`
		Personas struct {
			ComputationClass string `json:"computation_class"`
			ComputationID    string `json:"computation_id"`
			ComputationKey   string `json:"computation_key"`
			Namespace        string `json:"namespace"`
			SpaceID          string `json:"space_id"`
		} `json:"personas"`
		Protocols struct {
			SourceID string `json:"sourceId"`
		} `json:"protocols"`
	} `json:"context"`
	Integrations struct {
		All        bool `json:"All"`
		Warehouses struct {
			All bool `json:"all"`
		} `json:"Warehouses"`
		Webhooks bool `json:"Webhooks"`
	} `json:"integrations"`
	MessageID         string     `json:"messageId"`
	OriginalTimestamp time.Time  `json:"originalTimestamp"`
	ProjectID         string     `json:"projectId"`
	ReceivedAt        time.Time  `json:"receivedAt"`
	SentAt            *time.Time `json:"sentAt"`
	Timestamp         time.Time  `json:"timestamp"`
	Traits            struct {
		Email   string `json:"email"`
		LWTest2 bool   `json:"lw_test_2"`
	} `json:"traits"`
	Type     string `json:"type"`
	UserID   string `json:"userId"`
	Version  int    `json:"version"`
	WriteKey string `json:"writeKey"`
}

// randomString generates a random string of the given length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// generateRandomJSON generates a JSONData object with random values
func generateRandomJSON() JSONData {
	return JSONData{
		Channel: "server",
		Context: struct {
			Library struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"library"`
			Personas struct {
				ComputationClass string `json:"computation_class"`
				ComputationID    string `json:"computation_id"`
				ComputationKey   string `json:"computation_key"`
				Namespace        string `json:"namespace"`
				SpaceID          string `json:"space_id"`
			} `json:"personas"`
			Protocols struct {
				SourceID string `json:"sourceId"`
			} `json:"protocols"`
		}{
			Library: struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			}{
				Name:    "unknown",
				Version: "unknown",
			},
			Personas: struct {
				ComputationClass string `json:"computation_class"`
				ComputationID    string `json:"computation_id"`
				ComputationKey   string `json:"computation_key"`
				Namespace        string `json:"namespace"`
				SpaceID          string `json:"space_id"`
			}{
				ComputationClass: "audience",
				ComputationID:    "aud_" + randomString(12),
				ComputationKey:   "lw_test_" + randomString(2),
				Namespace:        "spa_" + randomString(16),
				SpaceID:          "spa_" + randomString(16),
			},
			Protocols: struct {
				SourceID string `json:"sourceId"`
			}{
				SourceID: randomString(12),
			},
		},
		Integrations: struct {
			All        bool `json:"All"`
			Warehouses struct {
				All bool `json:"all"`
			} `json:"Warehouses"`
			Webhooks bool `json:"Webhooks"`
		}{
			All: false,
			Warehouses: struct {
				All bool `json:"all"`
			}{
				All: false,
			},
			Webhooks: true,
		},
		MessageID:         "personas_" + randomString(12),
		OriginalTimestamp: time.Now(),
		ProjectID:         randomString(12),
		ReceivedAt:        time.Now(),
		SentAt:            nil,
		Timestamp:         time.Now(),
		Traits: struct {
			Email   string `json:"email"`
			LWTest2 bool   `json:"lw_test_2"`
		}{
			Email:   "test_" + randomString(10) + "@gmail.com",
			LWTest2: false,
		},
		Type:     "identify",
		UserID:   "test_" + randomString(10),
		Version:  2,
		WriteKey: randomString(32),
	}
}

// writeJSONToFile writes the JSON data to a file
func writeJSONToFile(filePath string, data JSONData) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// generateData generates data and writes it to a file until the target size is reached
func generateData(filePath string, targetSize int64) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var currentSize int64
	for currentSize < targetSize {
		jsonData := generateRandomJSON()
		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			return err
		}
		n, err := file.WriteString(string(jsonBytes) + "\n")
		if err != nil {
			return err
		}
		currentSize += int64(n)
	}

	return nil
}

func main() {
	buckets := map[string]int64{
		"100KB":  100 * 1024,
		"<1MB":   1024 * 1024,
		"1-10MB": 10 * 1024 * 1024,
	}

	outputDir := "random_json_files"
	os.Mkdir(outputDir, 0755)

	for bucket, size := range buckets {
		filePath := fmt.Sprintf("%s/%s.txt", outputDir, bucket) // Changed extension to .txt
		err := generateData(filePath, size)
		if err != nil {
			fmt.Printf("Error generating data for %s: %v\n", bucket, err)
		} else {
			fmt.Printf("Generated file for %s: %s\n", bucket, filePath)
		}
	}
}
