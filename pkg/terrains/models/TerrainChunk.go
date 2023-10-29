package models

import (
	"time"

	"github.com/google/uuid"
	sandboxModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	"gorm.io/gorm"
)

type Heightmap [][]float64

type TerrainChunk struct {
	ID                  string                `gorm:"type:uuid;primary_key;unique;" json:"id"`
	SandboxId           string                `gorm:"type:uuid;not null;" json:"sandboxId"`
	Position            sandboxModels.Vector3 `gorm:"type:jsonb;not null" json:"position"`
	Dimension           int                   `gorm:"type:int;not null" json:"dimension"`
	Height              int                   `gorm:"type:int;not null" json:"height"`
	DetailResolution    int                   `gorm:"type:int;not null" json:"detailResolution"`
	ResolutionPerPatch  int                   `gorm:"type:int;not null" json:"resolutionPerPatch"`
	HeightmapResolution int                   `gorm:"type:int;not null" json:"heightmapResolution"`
	AlphamapResolution  int                   `gorm:"type:int;not null" json:"alphamapResolution"`
	PerlinSeed          int64                 `gorm:"type:int;not null" json:"perlinSeed"`
	PerlinAlpha         float32               `gorm:"type:float;not null" json:"perlinAlpha"`
	PerlinBeta          float32               `gorm:"type:float;not null" json:"perlinBeta"`
	PerlinN             float32               `gorm:"type:float;not null" json:"perlinN"`
	PerlinAmplitude     float32               `gorm:"type:float;not null" json:"perlinAmplitude"`
	PerlinFrequency     float32               `gorm:"type:float;not null" json:"perlinFrequency"`
	Heightmap           string                `gorm:"type:text" json:"heightmap"`
	Created             time.Time             `gorm:"autoCreateTime" json:"created"`
	Updated             time.Time             `gorm:"autoUpdateTime" json:"updated"`
}

func (chunk *TerrainChunk) BeforeCreate(tx *gorm.DB) (err error) {
	chunk.ID = uuid.New().String()
	return
}
