// src/models/user.go
package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
    ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();" json:"id"`
    Name         string    `json:"name"`
    Email        string    `json:"email"`
    Password     string    `json:"password"`
	OnlineStatus bool      `gorm:"default:false" json:"online_status"` // Default value set to false
	ProfileImage string    `gorm:"default:''" json:"profile_image"`
    gorm.Model
}


func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return
}