package models

import (
	"time"

	"github.com/google/uuid"
	sandboxModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	"gorm.io/gorm"
)

type Heightmap [][]float64

type Resolution struct {
	X int `json:"x" gorm:"type:int"`
	Y int `json:"y" gorm:"type:int"`
}

type TerrainChunk struct {
	ID                  string                `gorm:"type:uuid;primary_key;unique;" json:"id"`
	Position            sandboxModels.Vector3 `gorm:"type:jsonb;not null" json:"position"`
	Dimension           int                   `json:"dimension" gorm:"type:int;not null"`
	Height              int                   `json:"height" gorm:"type:int;not null"`
	DetailResolution    Resolution            `json:"detailResolution" gorm:"type:jsonb;not null"`
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
