package models

import "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"

type Resolution struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type TerrainChunk struct {
	Position            models.Vector3 `json:"position"`
	Dimension           int            `json:"dimension"`
	Height              int            `json:"height"`
	DetailResolution    Resolution     `json:"detailResolution"`
	HeightmapResolution int            `json:"heightmapResolution"`
	AlphamapResolution  int            `json:"alphamapResolution"`
	Heightmap           Heightmap      `json:"heightmap"`
}
