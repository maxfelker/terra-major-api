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
)

func floatPtr(f float32) *float32 {
	return &f
}

func NewWorld(chunkCount, chunkDimension, height int, seed int64) []*models.TerrainChunk {
	var detailResolution = chunkDimension
	var resolutionPerPatch = 16
	var alphamapResolution = chunkDimension
	var heightmapResolution = chunkDimension
	var terrainHeight = height

	var chunks []*models.TerrainChunk
	for i := 0; i < chunkCount; i++ {
		for j := 0; j < chunkCount; j++ {
			globalX := float32(i * chunkDimension)
			globalZ := float32(j * chunkDimension)
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
