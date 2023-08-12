package terrains

import (
	"encoding/json"

	"github.com/aquilax/go-perlin"
	sandboxModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	models "github.com/mw-felker/terra-major-api/pkg/terrains/models"
)

const (
	alpha           = 2.
	beta            = 2.
	n               = 3
	perlinHeight    = 1
	perlinDamper    = 1000.0
	perlinFrequency = 0.005 // Lower value for broader features
	perlinAmplitude = 0.85  // Higher value for taller features

	chunkDimension   = 128
	height           = 64
	chunkCount       = 2
	neighborhoodSize = 4
)

func floatPtr(f float32) *float32 {
	return &f
}

func GenerateChunks(chunkCount, chunkDimension, height int, seed int64, offset sandboxModels.Vector3) []*models.TerrainChunk {
	var detailResolution = chunkDimension
	var resolutionPerPatch = 16
	var alphamapResolution = chunkDimension
	var heightmapResolution = chunkDimension
	var terrainHeight = height

	var chunks []*models.TerrainChunk
	for i := 0; i < chunkCount; i++ {
		for j := 0; j < chunkCount; j++ {
			globalX := float32(i*chunkDimension) + *offset.X
			globalZ := float32(j*chunkDimension) + *offset.Z
			pos := sandboxModels.Vector3{
				X: &globalX,
				Y: floatPtr(0),
				Z: &globalZ,
			}
			newChunk := NewTerrainChunk(
				pos,
				chunkDimension+1,
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

func NewTerrainChunk(pos sandboxModels.Vector3, dimension, terrainHeight, detailResolution, resolutionPerPatch, heightmapRes, alphamapRes int, seed int64) *models.TerrainChunk {
	heightmap := NewHeightmap(dimension, dimension, seed, pos)
	heightmapJSON, err := json.Marshal(heightmap)
	if err != nil {
		panic(err)
	}

	return &models.TerrainChunk{
		Position:            pos,
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

			// Scale the coordinates by the frequency
			noiseValue := perlin.Noise2D(globalX*perlinFrequency, globalY*perlinFrequency)

			// Scale the noise value by the amplitude
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
	chunks := GenerateChunks(chunkCount, chunkDimension, height, seed, offset)
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
	for i := 0; i < neighborhoodSize; i++ {
		for j := 0; j < neighborhoodSize; j++ {
			groupX := float32(i * chunkCount * chunkDimension)
			groupZ := float32(j * chunkCount * chunkDimension)
			group := CreateChunkGroup(seed, sandboxModels.Vector3{
				X: &groupX,
				Y: floatPtr(0),
				Z: &groupZ,
			})
			// Update the position of the chunk group based on its position in the neighborhood grid
			group.Position.X = &groupX
			group.Position.Z = &groupZ
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
