package entity

import (
	"time"

	"event_ticket_booking/constant"

	"gorm.io/gorm"
)

type Entity struct {
	Id           uint64         `gorm:"primaryKey;column:id"`
	Name         string         `gorm:"column:name"`
	Description  string         `gorm:"column:description"`
	DateTime     time.Time      `gorm:"column:date_time"`
	TotalTickets uint64         `gorm:"column:total_tickets"`
	TicketPrice  uint64         `gorm:"column:ticket_price"`
	Status       int            `gorm:"column:status"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	CreatedBy    uint64         `gorm:"column:created_by"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	UpdatedBy    uint64         `gorm:"column:updated_by"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
	DeletedBy    uint64         `gorm:"column:deleted_by"`
}

func (e Entity) TableName() string {
	return "event"
}

func (e *Entity) BeforeDelete(tx *gorm.DB) error {
	if e.Id == 0 {
		return constant.ErrMissingID
	}
	return tx.Model(&Entity{}).
		Where("id = ?", e.Id).
		UpdateColumns(map[string]any{
			"deleted_by": e.DeletedBy,
			"status":     constant.EVENT_STATUS_DELETED,
		}).Error
}

type Filter struct {
	Id       uint64
	Name     string
	DateTime string
	Status   int
}

func (f Filter) Apply(query *gorm.DB) *gorm.DB {
	table := Entity{}.TableName()
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
