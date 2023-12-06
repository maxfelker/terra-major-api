package terrains

import (
	"log"

	"github.com/google/uuid"
	sandboxModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	terrainModels "github.com/mw-felker/terra-major-api/pkg/terrains/models"
)

const (
	// Heightmap noise
	seed       = 1004
	frequency  = 0.01  // Lower value for broader features
	gain       = 0.001 // Higher value for taller features
	octaves    = 3
	lacunarity = 2.0
	// Grouping
	chunkPerGroup         = 2
	groupsPerNeighborhood = 2
	// Chunk Config
	heightmapResolution = 513 // (129, 257, 469, 513, 769, 1025, 2049)
	chunkDimension      = 512 // must be one smaller than heightMap resolution
	chunkHeight         = 256
	alphamapResolution  = 512
	detailResolution    = 512
	resolutionPerPatch  = 16 // https://docs.unity3d.com/ScriptReference/TerrainData.SetDetailResolution.html
)

func GenerateChunksForSandbox(sandboxId string) []*terrainModels.TerrainChunk {

	if sandboxId == "" {
		log.Fatalln("Please provide a sandbox ID when storing chunks")
		return nil
	}

	chunkNeighborhood := createChunkNeighborhood()
	chunks := flattenChunksArray(chunkNeighborhood)

	for _, chunk := range chunks {
		chunk.SandboxId = sandboxId
	}

	return chunks
}

func floatPtr(f float32) *float32 {
	return &f
}

func generateChunks(offset sandboxModels.Vector3) []*terrainModels.TerrainChunk {
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
				Seed:                seed,
				Frequency:           frequency,
				Gain:                gain,
				Octaves:             octaves,
				Lacunarity:          lacunarity,
			}
			chunks = append(chunks, newChunk)
		}
	}
	return chunks
}

func createChunkNeighborhood() *terrainModels.ChunkNeighborhood {
	groups := generateChunkGroups()
	return &terrainModels.ChunkNeighborhood{
		Position: sandboxModels.Vector3{
			X: floatPtr(0),
			Y: floatPtr(0),
			Z: floatPtr(0),
		},
		Groups: groups,
	}
}

func createChunkGroup(offset sandboxModels.Vector3) *terrainModels.ChunkGroup {
	chunks := generateChunks(offset)
	return &terrainModels.ChunkGroup{
		Position: sandboxModels.Vector3{
			X: floatPtr(0),
			Y: floatPtr(0),
			Z: floatPtr(0),
		},
		Chunks: chunks,
	}
}

func generateChunkGroups() []*terrainModels.ChunkGroup {
	var groups []*terrainModels.ChunkGroup
	halfSize := groupsPerNeighborhood / 2
	for i := -halfSize; i < halfSize; i++ {
		for j := -halfSize; j < halfSize; j++ {
			groupX := float32(i * chunkPerGroup * chunkDimension)
			groupZ := float32(j * chunkPerGroup * chunkDimension)
			group := createChunkGroup(sandboxModels.Vector3{
				X: &groupX,
				Y: floatPtr(0),
				Z: &groupZ,
			})
			groups = append(groups, group)
		}
	}
	return groups
}

func flattenChunksArray(neighborhood *terrainModels.ChunkNeighborhood) []*terrainModels.TerrainChunk {
	var allChunks []*terrainModels.TerrainChunk
	for _, group := range neighborhood.Groups {
		for _, chunk := range group.Chunks {
			allChunks = append(allChunks, chunk)
		}
	}
	return allChunks
}
