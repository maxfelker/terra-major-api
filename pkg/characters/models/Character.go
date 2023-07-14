package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Character struct {
	ID           string    `gorm:"type:uuid;primary_key;unique;" json:"id"`
	AccountId    string    `gorm:"type:uuid;not null;unique;" json:"accountId"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Bio          string    `gorm:"type:text;default:''" json:"bio"`
	Age          int       `gorm:"type:int;default:25" json:"age"`
	Strength     int       `gorm:"type:int;default:5" json:"strength"`
	Intelligence int       `gorm:"type:int;default:5" json:"intelligence"`
	Endurance    int       `gorm:"type:int;default:5" json:"endurance"`
	Agility      int       `gorm:"type:int;default:5" json:"agility"`
	Created      time.Time `gorm:"autoCreateTime" json:"created"`
	Updated      time.Time `gorm:"autoUpdateTime" json:"updated"`
}

func (character *Character) BeforeCreate(tx *gorm.DB) (err error) {
	character.ID = uuid.New().String()
	return
}
