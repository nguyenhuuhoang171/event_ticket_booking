package entity

import (
	"time"

	"gorm.io/gorm"
)

type Entity struct {
	Id             uint64    `gorm:"primaryKey;column:id"`
	Email          string    `gorm:"column:email"`
	HashedPassword string    `gorm:"column:hashed_password"`
	Role           string    `gorm:"column:role"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (e Entity) TableName() string {
	return "user"
}

type Filter struct {
	Id    uint64
	Email string
}

// Apply adds this filter's conditions to query. It implements base.Filter so
// the generic repository can build queries without knowing the entity.
func (f Filter) Apply(query *gorm.DB) *gorm.DB {
	if f.Id != 0 {
		query = query.Where("id = ?", f.Id)
	}
	if f.Email != "" {
		query = query.Where("email = ?", f.Email)
	}
	return query
}
