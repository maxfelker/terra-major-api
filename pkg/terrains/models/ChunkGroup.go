package models

import (
	"time"

	"github.com/google/uuid"
	sandboxModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"

	"gorm.io/gorm"
)

type ChunkGroup struct {
	ID       string                `gorm:"type:uuid;primary_key;unique;" json:"id"`
	Position sandboxModels.Vector3 `gorm:"type:jsonb;not null" json:"position"`
	Chunks   []*TerrainChunk       `gorm:"foreignKey:ChunkGroupID" json:"chunks"`
	Created  time.Time             `gorm:"autoCreateTime" json:"created"`
	Updated  time.Time             `gorm:"autoUpdateTime" json:"updated"`
}

func (chunkGroup *ChunkGroup) BeforeCreate(tx *gorm.DB) (err error) {
	chunkGroup.ID = uuid.New().String()
	return
}
