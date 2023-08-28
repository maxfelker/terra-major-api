package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	terrains "github.com/mw-felker/terra-major-api/pkg/terrains"
)

func main() {
	var sandboxId string
	flag.StringVar(&sandboxId, "sandboxId", "", "UUIDv4 of the sandbox")
	flag.Parse()

	if sandboxId == "" {
		fmt.Println("Please provide a sandbox ID using the -sandboxId flag.")
		return
	}

	startTime := time.Now()

	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(420) + 1
	seed := int64(randomNumber)
	chunkNeighborhood := terrains.CreateChunkNeighborhood(seed)
	chunks := terrains.FlattenChunksArray(chunkNeighborhood)

	for _, chunk := range chunks {
		chunk.SandboxId = sandboxId
	}

	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("chunks_%d.json", timestamp)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Failed to create JSON file: %v\n", err)
		return
	}
	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	encoder := json.NewEncoder(bufferedWriter)
	err = encoder.Encode(chunks)
	if err != nil {
		fmt.Printf("Failed to write chunks to JSON file: %v\n", err)
		return
	}
	bufferedWriter.Flush()

	fileInfo, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("Failed to get file info: %v\n", err)
		return
	}

	endTime := time.Now()
	elapsedTime := float64(int(endTime.Sub(startTime).Seconds()*10+0.5)) / 10
	fileSizeMB := float64(fileInfo.Size()) / 1048576.0

	fmt.Printf("Number of chunks created: %d\n", len(chunks))
	fmt.Printf("Time taken to create the file: %.1fs\n", elapsedTime)
	fmt.Printf("Size of the file: %.2f MB\n", fileSizeMB)
}
