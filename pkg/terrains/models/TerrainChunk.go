package models

import "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"

type TerrainChunk struct {
	Position  models.Vector3 `json:"position"`
	Seed      int64          `json:"seed"`
	Heightmap Heightmap      `json:"heightmap"`
}
