package main

import (
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

	encoder := json.NewEncoder(os.Stdout)
	err := encoder.Encode(chunks)
	if err != nil {
		fmt.Printf("Failed to write chunks to JSON output: %v\n", err)
		return
	}

	endTime := time.Now()
	elapsedTime := float64(int(endTime.Sub(startTime).Seconds()*10+0.5)) / 10

	fmt.Fprintf(os.Stderr, "Number of chunks created: %d\n", len(chunks))
	fmt.Fprintf(os.Stderr, "Time taken to create the data: %.1fs\n", elapsedTime)
}
