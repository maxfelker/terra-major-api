package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type BaseAccount struct {
	ID      string    `gorm:"type:uuid;primary_key;unique;" json:"id"`
	Email   string    `gorm:"type:varchar(100);unique" json:"email"`
	Created time.Time `gorm:"autoCreateTime" json:"created"`
	Updated time.Time `gorm:"autoUpdateTime" json:"updated"`
}

type Account struct {
	BaseAccount
	Password string `gorm:"type:varchar(100)" json:"password"`
}

type AccountResponse struct {
	BaseAccount
}

func generatePassword(passwordText string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordText), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func (account *Account) BeforeCreate(tx *gorm.DB) (err error) {
	account.ID = uuid.New().String()
	if len(account.Password) > 0 {
		account.Password = generatePassword(account.Password)
	}
	return
}

func (account *Account) BeforeSave(tx *gorm.DB) (err error) {
	if len(account.Password) > 0 {
		var hash = generatePassword(account.Password)
		tx.Statement.SetColumn("Password", hash)
	}
	return
}
