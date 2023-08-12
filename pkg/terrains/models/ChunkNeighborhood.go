package models

import (
	"time"

	"github.com/google/uuid"
	sandboxModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	"gorm.io/gorm"
)

type ChunkNeighborhood struct {
	ID       string                `gorm:"type:uuid;primary_key;unique;" json:"id"`
	Position sandboxModels.Vector3 `gorm:"type:jsonb;not null" json:"position"`
	Groups   []*ChunkGroup         `gorm:"foreignKey:ChunkNeighborhoodID" json:"groups"`
	Created  time.Time             `gorm:"autoCreateTime" json:"created"`
	Updated  time.Time             `gorm:"autoUpdateTime" json:"updated"`
}

func (neighborhood *ChunkNeighborhood) BeforeCreate(tx *gorm.DB) (err error) {
	neighborhood.ID = uuid.New().String()
	return
}
