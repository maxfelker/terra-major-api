package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Sandbox struct {
	ID          string    `gorm:"type:uuid;primary_key;unique;" json:"id"`
	CharacterId string    `gorm:"type:uuid;not null;unique;" json:"characterId"`
	AccountId   string    `gorm:"type:uuid;not null;" json:"accountId"`
	Created     time.Time `gorm:"autoCreateTime" json:"created"`
	Updated     time.Time `gorm:"autoUpdateTime" json:"updated"`
}

func (sandbox *Sandbox) BeforeCreate(tx *gorm.DB) (err error) {
	sandbox.ID = uuid.New().String()
	return
}
