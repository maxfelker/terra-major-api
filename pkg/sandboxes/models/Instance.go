package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	GM_UUID = "f1c9f0a6-695f-424a-9424-60a5e96032df"
)

type Instance struct {
	ID          string    `gorm:"type:uuid;primary_key;unique;" json:"id"`
	SandboxId   string    `gorm:"type:uuid;not null" json:"sandboxId"`
	CharacterId string    `gorm:"type:uuid;not null" json:"characterId"`
	PrefabName  string    `gorm:"type:string;not null" json:"prefabName"`
	Position    Vector3   `gorm:"type:jsonb;not null" json:"position"`
	Rotation    Vector3   `gorm:"type:jsonb;not null" json:"rotation"`
	Created     time.Time `gorm:"autoCreateTime" json:"created"`
	Updated     time.Time `gorm:"autoUpdateTime" json:"updated"`
}

func (instance *Instance) BeforeCreate(tx *gorm.DB) (err error) {
	instance.ID = uuid.New().String()
	return
}
