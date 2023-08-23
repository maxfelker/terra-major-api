package models

import (
	"time"

	sandboxModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
)

type Heightmap [][]float64

type TerrainChunk struct {
	ID                  string                `json:"id"`
	SandboxId           string                `json:"sandboxId"`
	Position            sandboxModels.Vector3 `json:"position"`
	Dimension           int                   `json:"dimension"`
	Height              int                   `json:"height"`
	DetailResolution    int                   `json:"detailResolution"`
	ResolutionPerPatch  int                   `json:"resolutionPerPatch"`
	HeightmapResolution int                   `json:"heightmapResolution"`
	AlphamapResolution  int                   `json:"alphamapResolution"`
	Heightmap           string                `json:"heightmap"`
	Created             time.Time             `json:"created"`
	Updated             time.Time             `json:"updated"`
}
