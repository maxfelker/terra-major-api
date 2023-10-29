package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/google/uuid"
	sandboxModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	terrainModels "github.com/mw-felker/terra-major-api/pkg/terrains/models"
)

const (
	// Heightmap noise
	seed            = int64(13)
	alpha           = 1
	beta            = 2
	n               = 3
	perlinFrequency = 0.0005 // Lower value for broader features
	perlinAmplitude = 0.85   // Higher value for taller features
	// Grouping
	chunkPerGroup         = 2
	groupsPerNeighborhood = 2
	// Chunk Config
	heightmapResolution = 1025 // (129, 257, 469, 513, 769, 1025, 2049)
	chunkDimension      = 1024 // must be one smaller than heightMap resolution
	chunkHeight         = 256
	alphamapResolution  = 1024
	detailResolution    = 1024
	resolutionPerPatch  = 16 // https://docs.unity3d.com/ScriptReference/TerrainData.SetDetailResolution.html
)

func floatPtr(f float32) *float32 {
	return &f
}

func main() {
	var sandboxId string
	flag.StringVar(&sandboxId, "sandboxId", "", "UUIDv4 of the sandbox")
	flag.Parse()

	if sandboxId == "" {
		fmt.Println("Please provide a sandbox ID using the -sandboxId flag.")
		return
	}

	//startTime := time.Now()
	chunkNeighborhood := CreateChunkNeighborhood()
	chunks := FlattenChunksArray(chunkNeighborhood)

	for _, chunk := range chunks {
		chunk.SandboxId = sandboxId
	}

	encoder := json.NewEncoder(os.Stdout)
	err := encoder.Encode(chunks)
	if err != nil {
		fmt.Printf("Failed to write chunks to JSON output: %v\n", err)
		return
	}
	/*
		endTime := time.Now()
		elapsedTime := float64(int(endTime.Sub(startTime).Seconds()*10+0.5)) / 10

		fmt.Fprintf(os.Stderr, "Number of chunks created: %d\n", len(chunks))
		fmt.Fprintf(os.Stderr, "Time taken to create the data: %.1fs\n", elapsedTime)*/
}

func GenerateChunks(offset sandboxModels.Vector3) []*terrainModels.TerrainChunk {
	var terrainHeight = chunkHeight
	var chunks []*terrainModels.TerrainChunk
	for i := 0; i < chunkPerGroup; i++ {
		for j := 0; j < chunkPerGroup; j++ {
			globalX := float32(i*chunkDimension) + *offset.X
			globalZ := float32(j*chunkDimension) + *offset.Z
			position := sandboxModels.Vector3{
				X: &globalX,
				Y: floatPtr(0),
				Z: &globalZ,
			}
			newChunk := &terrainModels.TerrainChunk{
				ID:                  uuid.New().String(),
				Position:            position,
				Dimension:           chunkDimension,
				Height:              terrainHeight,
				DetailResolution:    detailResolution,
				ResolutionPerPatch:  resolutionPerPatch,
				HeightmapResolution: heightmapResolution,
				AlphamapResolution:  alphamapResolution,
				PerlinSeed:          seed,
				PerlinAlpha:         alpha,
				PerlinBeta:          beta,
				PerlinN:             n,
				PerlinAmplitude:     perlinAmplitude,
				PerlinFrequency:     perlinFrequency,
			}
			chunks = append(chunks, newChunk)
		}
	}
	return chunks
}

func CreateChunkNeighborhood() *terrainModels.ChunkNeighborhood {
	groups := GenerateChunkGroups()
	return &terrainModels.ChunkNeighborhood{
		Position: sandboxModels.Vector3{
			X: floatPtr(0),
			Y: floatPtr(0),
			Z: floatPtr(0),
		},
		Groups: groups,
	}
}

func CreateChunkGroup(offset sandboxModels.Vector3) *terrainModels.ChunkGroup {
	chunks := GenerateChunks(offset)
	return &terrainModels.ChunkGroup{
		Position: sandboxModels.Vector3{
			X: floatPtr(0),
			Y: floatPtr(0),
			Z: floatPtr(0),
		},
		Chunks: chunks,
	}
}

func GenerateChunkGroups() []*terrainModels.ChunkGroup {
	var groups []*terrainModels.ChunkGroup
	halfSize := groupsPerNeighborhood / 2
	for i := -halfSize; i < halfSize; i++ {
		for j := -halfSize; j < halfSize; j++ {
			groupX := float32(i * chunkPerGroup * chunkDimension)
			groupZ := float32(j * chunkPerGroup * chunkDimension)
			group := CreateChunkGroup(sandboxModels.Vector3{
				X: &groupX,
				Y: floatPtr(0),
				Z: &groupZ,
			})
			groups = append(groups, group)
		}
	}
	return groups
}

func FlattenChunksArray(neighborhood *terrainModels.ChunkNeighborhood) []*terrainModels.TerrainChunk {
	var allChunks []*terrainModels.TerrainChunk
	for _, group := range neighborhood.Groups {
		for _, chunk := range group.Chunks {
			allChunks = append(allChunks, chunk)
		}
	}
	return allChunks
}
