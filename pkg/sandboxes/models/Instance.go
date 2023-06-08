package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Instance struct {
	ID          string    `gorm:"type:uuid;primary_key;unique;" json:"id"`
	SandboxId   string    `gorm:"type:uuid;not null" json:"sandboxId"`
	LayerId     string    `gorm:"type:uuid;not null" json:"layerId"`
	CharacterId string    `gorm:"type:uuid;not null" json:"characterId"`
	PrefabName  string    `gorm:"type:uuid;not null" json:"prefabName"`
	Health      int       `gorm:"type:int; not null; default: 100;" json:"health"`
	Position    Vector3   `gorm:"type:varchar(100);not null" json:"position"`
	Rotation    Vector3   `gorm:"type:varchar(100);not null" json:"rotation"`
	Created     time.Time `gorm:"autoCreateTime" json:"created"`
	Updated     time.Time `gorm:"autoUpdateTime" json:"updated"`
}

func (instance *Instance) BeforeCreate(tx *gorm.DB) (err error) {
	instance.ID = uuid.New().String()
	return
}
