package terrains

import (
	"log"

	"github.com/google/uuid"
	sandboxModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	terrainModels "github.com/mw-felker/terra-major-api/pkg/terrains/models"
)

const (
	chunksPerGeneration = 100
	seed                = 1920
	frequency           = 0.005
	gain                = 0.001
	octaves             = 3
	lacunarity          = 2.0
	heightmapResolution = 513
	chunkDimension      = 512
	chunkHeight         = 64
	alphamapResolution  = 512
	detailResolution    = 512
	resolutionPerPatch  = 16
)

func GenerateChunksForSandbox(sandboxId string) []*terrainModels.TerrainChunk {
	if sandboxId == "" {
		log.Fatalln("Please provide a sandbox ID when storing chunks")
		return nil
	}

	chunks := generateChunks(sandboxId)
	return chunks
}

func floatPtr(f float32) *float32 {
	return &f
}

func generateChunks(sandboxId string) []*terrainModels.TerrainChunk {
	var terrainHeight = chunkHeight
	var chunks []*terrainModels.TerrainChunk
	gridSize := int(chunksPerGeneration / 2)
	offset := int(gridSize / 2)

	for i := -offset; i < offset; i++ {
		for j := -offset; j < offset; j++ {
			globalX := float32(i * chunkDimension)
			globalZ := float32(j * chunkDimension)
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
				Seed:                seed,
				Frequency:           frequency,
				Gain:                gain,
				Octaves:             octaves,
				Lacunarity:          lacunarity,
				SandboxId:           sandboxId,
			}
			chunks = append(chunks, newChunk)
		}
	}
	return chunks
}
