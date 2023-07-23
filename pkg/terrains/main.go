package terrains

import (
	"math/rand"

	sandboxModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	models "github.com/mw-felker/terra-major-api/pkg/terrains/models"
	"github.com/ojrac/opensimplex-go"
)

func floatPtr(f float32) *float32 {
	return &f
}

func RandomTerrain(worldDim int) *models.TerrainChunk {
	var dimension = worldDim
	var pos = sandboxModels.Vector3{X: floatPtr(1.0), Y: floatPtr(2.0), Z: floatPtr(3.0)}
	return NewTerrainChunk(pos, dimension, dimension)
}

func NewTerrainChunk(pos sandboxModels.Vector3, width, height int) *models.TerrainChunk {
	seed := rand.Int63()
	hm := NewHeightmap(width, height, seed)

	return &models.TerrainChunk{
		Position:  pos,
		Seed:      seed,
		Heightmap: models.Heightmap(hm),
	}
}

func NewHeightmap(width, height int, seed int64) models.Heightmap {
	heightmap := make(models.Heightmap, height)
	for i := range heightmap {
		heightmap[i] = make([]float64, width)
	}

	noise := opensimplex.NewNormalized(seed)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			heightmap[y][x] = noise.Eval2(float64(x), float64(y))
		}
	}

	return heightmap
}
