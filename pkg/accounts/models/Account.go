package models

import (
	"fmt"
	"strings"
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
	passwordBytes := []byte(strings.TrimSpace(passwordText))
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func (account *Account) BeforeCreate(tx *gorm.DB) (err error) {
	account.ID = uuid.New().String()
	account.Email = strings.TrimSpace(account.Email)
	account.Password = generatePassword(account.Password)
	return
}

func (account *Account) BeforeUpdate(tx *gorm.DB) (err error) {
	fmt.Println("Saving account")
	if tx.Statement.Changed("Password") {
		fmt.Println("Password change " + account.Password)
		account.Password = generatePassword(account.Password)
	}
	return
}
