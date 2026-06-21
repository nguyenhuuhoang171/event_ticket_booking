package repository

import (
	"gorm.io/gorm"
)

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
