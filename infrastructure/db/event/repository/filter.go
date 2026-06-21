package repository

import (
	"event_ticket_booking/infrastructure/db/event/entity"

	"gorm.io/gorm"
)

type Filter struct {
	Id       uint64
	Name     string
	DateTime string
	Status   int
}

func (f Filter) Apply(query *gorm.DB) *gorm.DB {
	table := entity.Entity{}.TableName()
	if f.Id != 0 {
		query = query.Where(table+".id = ?", f.Id)
	}
	if f.Name != "" {
		query = query.Where(table+".name = ?", f.Name)
	}
	if f.DateTime != "" {
		query = query.Where(table+".date_time = ?", f.DateTime)
	}
	if f.Status != 0 {
		query = query.Where(table+".status = ?", f.Status)
	}
	return query
}
