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
	SandboxId           string                `gorm:"type:uuid;not null" json:"sandboxId"`
	Position            sandboxModels.Vector3 `gorm:"type:jsonb;not null" json:"position"`
	Dimension           int                   `json:"dimension" gorm:"type:int;not null"`
	Height              int                   `json:"height" gorm:"type:int;not null"`
	DetailResolution    int                   `json:"detailResolution" gorm:"type:int;not null"`
	ResolutionPerPatch  int                   `json:"resolutionPerPatch" gorm:"type:int;not null"`
	HeightmapResolution int                   `json:"heightmapResolution" gorm:"type:int;not null"`
	AlphamapResolution  int                   `json:"alphamapResolution" gorm:"type:int;not null"`
	Heightmap           string                `json:"heightmap" gorm:"type:jsonb;not null"`
	Created             time.Time             `gorm:"autoCreateTime" json:"created"`
	Updated             time.Time             `gorm:"autoUpdateTime" json:"updated"`
}

func (chunk *TerrainChunk) BeforeCreate(tx *gorm.DB) (err error) {
	chunk.ID = uuid.New().String()
	return
}
