package repository

import (
	"event_ticket_booking/infrastructure/db/booking/entity"

	"gorm.io/gorm"
)

type Filter struct {
	Id      uint64
	EventId uint64
	UserId  uint64
	Status  int
}

func (f Filter) Apply(query *gorm.DB) *gorm.DB {
	table := entity.Entity{}.TableName()
	if f.Id != 0 {
		query = query.Where(table+".id = ?", f.Id)
	}
	if f.EventId != 0 {
		query = query.Where(table+".event_id = ?", f.EventId)
	}
	if f.UserId != 0 {
		query = query.Where(table+".user_id = ?", f.UserId)
	}
	if f.Status != 0 {
		query = query.Where(table+".status = ?", f.Status)
	}
	return query
}
