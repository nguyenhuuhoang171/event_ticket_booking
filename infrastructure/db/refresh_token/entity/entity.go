package entity

import (
	"time"

	"gorm.io/gorm"
)

type Entity struct {
	Id        uint64    `gorm:"primaryKey;column:id"`
	UserId    uint64    `gorm:"column:user_id"`
	Token     string    `gorm:"column:token"`
	ExpireAt  time.Time `gorm:"column:expire_at"`
	Status    *int      `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (e Entity) TableName() string {
	return "refresh_token"
}

type Filter struct {
	Id     uint64
	Token  string
	UserId uint64
	Status int
}

// Apply adds this filter's conditions to query. It implements base.Filter so
// the generic repository can build queries without knowing the entity.
func (f Filter) Apply(query *gorm.DB) *gorm.DB {
	if f.Id != 0 {
		query = query.Where("id = ?", f.Id)
	}
	if f.Token != "" {
		query = query.Where("token = ?", f.Token)
	}
	if f.UserId != 0 {
		query = query.Where("user_id = ?", f.UserId)
	}
	if f.Status != 0 {
		query = query.Where("status = ?", f.Status)
	}
	return query
}
