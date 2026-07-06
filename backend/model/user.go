package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `gorm:"primaryKey;size:36" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Password  string    `gorm:"size:128;not null" json:"-"`
	RealName  string    `gorm:"size:64" json:"realName"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u *User) SetPassword(raw string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) CheckPassword(raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(raw)) == nil
}

func FindUserByUsername(username string) (*User, error) {
	var user User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUserByID(id string) (*User, error) {
	var user User
	err := DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *User) error {
	return DB.Create(user).Error
}
