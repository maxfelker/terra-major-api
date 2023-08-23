package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/mw-felker/terra-major-api/pkg/core"
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

	app := core.CreateApp()

	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(420) + 1
	seed := int64(randomNumber)
	chunkNeighborhood := terrains.CreateChunkNeighborhood(seed)
	chunks := terrains.FlattenChunksArray(chunkNeighborhood)

	containerClient, err := app.NoSQL.NewContainer("chunks")
	if err != nil {
		fmt.Printf("Failed to get container client: %v\n", err)
		return
	}

	pk := azcosmos.NewPartitionKeyString(sandboxId)
	batch := containerClient.NewTransactionalBatch(pk)

	for _, chunk := range chunks {
		chunk.SandboxId = sandboxId
		marshalled, err := json.Marshal(chunk)
		if err != nil {
			fmt.Printf("Failed to marshal chunk: %v\n", err)
			return
		}
		batch.CreateItem(marshalled, nil)
	}
	ctx := context.Background()
	batchResponse, err := containerClient.ExecuteTransactionalBatch(ctx, batch, nil)
	if err != nil {
		fmt.Printf("Failed to execute transactional batch: %v\n", err)
		return
	}

	if batchResponse.Success {
		for index, operation := range batchResponse.OperationResults {
			fmt.Printf("Operation %v completed with status code %v consumed %v RU\n", index, operation.StatusCode, operation.RequestCharge)
		}
	} else {
		for index, operation := range batchResponse.OperationResults {
			if operation.StatusCode != http.StatusFailedDependency {
				fmt.Printf("Transaction failed due to operation %v which failed with status code %v\n", index, operation.StatusCode)
			}
		}
	}
}
