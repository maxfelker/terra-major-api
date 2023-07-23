package terrains

import (
	"math/rand"

	"github.com/aquilax/go-perlin"
	sandboxModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	models "github.com/mw-felker/terra-major-api/pkg/terrains/models"
)

const (
	alpha        = 2.
	beta         = 2.
	n            = 3
	perlinHeight = 1.5
	perlinDamper = 0.25
)

func floatPtr(f float32) *float32 {
	return &f
}

func NewWorld(worldDim, chunkDim int) []*models.TerrainChunk {
	var detailResolutionX = 128
	var detailResolutionY = 16
	var alphamapResolution = 128
	var heightmapResolution = 128

	var chunks []*models.TerrainChunk
	for i := 0; i < worldDim; i++ {
		for j := 0; j < worldDim; j++ {
			pos := sandboxModels.Vector3{
				X: floatPtr(float32(i * chunkDim)),
				Y: floatPtr(0),
				Z: floatPtr(float32(j * chunkDim)),
			}
			chunks = append(chunks, NewTerrainChunk(pos, chunkDim, chunkDim, detailResolutionX, detailResolutionY, heightmapResolution, alphamapResolution))
		}
	}
	return chunks
}

func NewTerrainChunk(pos sandboxModels.Vector3, dimension, height, detailResX, detailResY, heightmapRes, alphamapRes int) *models.TerrainChunk {
	seed := rand.Int63()
	hm := NewHeightmap(dimension, dimension, seed)

	return &models.TerrainChunk{
		Position:            pos,
		Dimension:           dimension,
		Height:              height,
		DetailResolution:    models.Resolution{X: detailResX, Y: detailResY},
		HeightmapResolution: heightmapRes,
		AlphamapResolution:  alphamapRes,
		Heightmap:           hm,
	}
}

func NewHeightmap(width, height int, seed int64) models.Heightmap {
	heightmap := make(models.Heightmap, height)
	for i := range heightmap {
		heightmap[i] = make([]float64, width)
	}

	perlin := perlin.NewPerlin(alpha, beta, n, seed)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			heightmap[y][x] = perlinHeight * perlin.Noise2D(float64(x)/perlinDamper, float64(y)/perlinDamper)
		}
	}

	return heightmap
}
