package main

import (
	"compress/bzip2"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Define a Go struct that represents the JSON objects.
type ExampleData struct {
	Model string `json:"model"`
}

const workerCount = 8

func main() {
	startTime := time.Now()
	// Replace this URL with the actual URL to your .json.bz2 file
	url := "https://quiz.storpool.com/bigf.json.bz2"

	// 1) Stream the file from URL
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("HTTP request failed with status: %v", resp.Status)
	}

	// 2) Wrap the HTTP response body with a bzip2 decompressor
	bz2Reader := bzip2.NewReader(resp.Body)

	// 3) Create a JSON decoder that reads from the decompressed stream
	decoder := json.NewDecoder(bz2Reader)

	// --- A) If your JSON is a large array: [ {...}, {...}, ... ]
	//

	// First, read the opening bracket `[`:
	t, err := decoder.Token()
	if err != nil {
		log.Fatalf("Error reading opening JSON token: %v", err)
	}
	if delim, ok := t.(json.Delim); !ok || delim != '[' {
		log.Fatalf("Expected [ at the beginning of JSON array")
	}

	// We'll keep a map of model -> count
	modelCount := make(map[string]int)

	var mu sync.Mutex

	// Channel for decoded JSON objects
	dataChannel := make(chan ExampleData, 10000)

	// Worker pool for concurrent processing
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localCount := make(map[string]int) // Local map to reduce contention

			for record := range dataChannel {
				localCount[record.Model]++
			}

			// Merge local count into the global map
			mu.Lock()
			for model, count := range localCount {
				modelCount[model] += count
			}
			mu.Unlock()
		}()
	}

	// Read JSON objects and send them to workers
	for decoder.More() {
		var record ExampleData
		if err := decoder.Decode(&record); err != nil {
			log.Fatalf("Error decoding JSON: %v", err)
		}
		dataChannel <- record
	}
	close(dataChannel)

	// Wait for workers to finish processing
	wg.Wait()

	// // Loop through array elements
	// for decoder.More() {
	// 	var record ExampleData
	// 	if err := decoder.Decode(&record); err != nil {
	// 		log.Fatalf("Error decoding JSON array element: %v", err)
	// 	}
	// 	// Process the record
	// 	// fmt.Printf("Record from array: %+v\n", record)
	// 	modelCount[record.Model]++

	// }

	// Finally, read the closing bracket `]`:
	t, err = decoder.Token()
	if err != nil {
		log.Fatalf("Error reading closing JSON token: %v", err)
	}
	if delim, ok := t.(json.Delim); !ok || delim != ']' {
		log.Fatalf("Expected ] at the end of JSON array")
	}

	elapsedTime := time.Since(startTime)
	// 7) Print results
	fmt.Printf("Total unique models: %d\n", len(modelCount))
	fmt.Println("Counts per model:")
	for model, count := range modelCount {
		fmt.Printf("  %s: %d\n", model, count)
	}
	fmt.Printf("\nExecution time: %s\n", elapsedTime)

}
