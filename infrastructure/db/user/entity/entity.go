package entity

import (
	"time"
)

type Entity struct {
	Id             uint64    `gorm:"primaryKey;column:id"`
	Email          string    `gorm:"column:email"`
	Name           string    `gorm:"column:name"`
	HashedPassword string    `gorm:"column:hashed_password"`
	Role           string    `gorm:"column:role"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	CreatedBy      uint64    `gorm:"column:created_by"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
	UpdatedBy      uint64    `gorm:"column:updated_by"`
}

func (e Entity) TableName() string {
	return "user"
}
