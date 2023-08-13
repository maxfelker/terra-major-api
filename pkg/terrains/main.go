package terrains

import (
	"encoding/json"

	"github.com/aquilax/go-perlin"
	sandboxModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	models "github.com/mw-felker/terra-major-api/pkg/terrains/models"
)

const (
	// Heightmap noise
	alpha           = 1
	beta            = 2
	n               = 3
	perlinFrequency = 0.0005 // Lower value for broader features
	perlinAmplitude = 0.85   // Higher value for taller features
	// Grouping
	chunkPerGroup         = 2
	groupsPerNeighborhood = 2
	// Chunk Config
	chunkDimension      = 1024
	chunkHeight         = 256
	heightmapResolution = 1025 // (129, 257, 469, 513, 769, 1025, 2049)
	alphamapResolution  = 1024
	detailResolution    = 256
	resolutionPerPatch  = 16
)

func floatPtr(f float32) *float32 {
	return &f
}

func GenerateChunks(chunkCount, chunkDimension, height int, seed int64, offset sandboxModels.Vector3) []*models.TerrainChunk {
	var terrainHeight = height
	var chunks []*models.TerrainChunk
	for i := 0; i < chunkCount; i++ {
		for j := 0; j < chunkCount; j++ {
			globalX := float32(i*chunkDimension) + *offset.X
			globalZ := float32(j*chunkDimension) + *offset.Z
			position := sandboxModels.Vector3{
				X: &globalX,
				Y: floatPtr(0),
				Z: &globalZ,
			}
			newChunk := NewTerrainChunk(
				position,
				chunkDimension,
				terrainHeight,
				detailResolution,
				resolutionPerPatch,
				heightmapResolution,
				alphamapResolution,
				seed,
			)
			chunks = append(chunks, newChunk)
		}
	}
	return chunks
}

func NewTerrainChunk(position sandboxModels.Vector3, dimension, terrainHeight, detailResolution, resolutionPerPatch, heightmapRes, alphamapRes int, seed int64) *models.TerrainChunk {
	heightmap := NewHeightmap(heightmapRes, heightmapRes, seed, position)
	heightmapJSON, err := json.Marshal(heightmap)
	if err != nil {
		panic(err)
	}

	return &models.TerrainChunk{
		Position:            position,
		Dimension:           dimension,
		Height:              terrainHeight,
		DetailResolution:    models.Resolution{X: detailResolution, Y: resolutionPerPatch},
		HeightmapResolution: heightmapRes,
		AlphamapResolution:  alphamapRes,
		Heightmap:           string(heightmapJSON),
	}
}

func NewHeightmap(width, depth int, seed int64, pos sandboxModels.Vector3) models.Heightmap {
	heightmap := make(models.Heightmap, depth)
	for i := range heightmap {
		heightmap[i] = make([]float64, width)
	}

	perlin := perlin.NewPerlin(alpha, beta, n, seed)

	for y := 0; y < depth; y++ {
		for x := 0; x < width; x++ {
			globalX := float64(*pos.X) + float64(x)
			globalY := float64(*pos.Z) + float64(y)
			noiseValue := perlin.Noise2D(globalX*perlinFrequency, globalY*perlinFrequency)
			heightmap[y][x] = perlinAmplitude * noiseValue
		}
	}

	return heightmap
}

func CreateChunkNeighborhood(seed int64) *models.ChunkNeighborhood {
	groups := GenerateChunkGroups(seed)
	return &models.ChunkNeighborhood{
		Position: sandboxModels.Vector3{
			X: floatPtr(0),
			Y: floatPtr(0),
			Z: floatPtr(0),
		},
		Groups: groups,
	}
}

func CreateChunkGroup(seed int64, offset sandboxModels.Vector3) *models.ChunkGroup {
	chunks := GenerateChunks(chunkPerGroup, chunkDimension, chunkHeight, seed, offset)
	return &models.ChunkGroup{
		Position: sandboxModels.Vector3{
			X: floatPtr(0),
			Y: floatPtr(0),
			Z: floatPtr(0),
		},
		Chunks: chunks,
	}
}

func GenerateChunkGroups(seed int64) []*models.ChunkGroup {
	var groups []*models.ChunkGroup
	halfSize := groupsPerNeighborhood / 2
	for i := -halfSize; i < halfSize; i++ {
		for j := -halfSize; j < halfSize; j++ {
			groupX := float32(i * chunkPerGroup * chunkDimension)
			groupZ := float32(j * chunkPerGroup * chunkDimension)
			group := CreateChunkGroup(seed, sandboxModels.Vector3{
				X: &groupX,
				Y: floatPtr(0),
				Z: &groupZ,
			})
			groups = append(groups, group)
		}
	}
	return groups
}

func FlattenChunksArray(neighborhood *models.ChunkNeighborhood) []*models.TerrainChunk {
	var allChunks []*models.TerrainChunk
	for _, group := range neighborhood.Groups {
		for _, chunk := range group.Chunks {
			allChunks = append(allChunks, chunk)
		}
	}
	return allChunks
}
