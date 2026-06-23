package entity

import (
	"time"
)

type Entity struct {
	Id        uint64    `gorm:"primaryKey;column:id"`
	EventId   uint64    `gorm:"column:event_id"`
	UserId    uint64    `gorm:"column:user_id"`
	Quantity  uint64    `gorm:"column:quantity"`
	Status    int       `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	CreatedBy uint64    `gorm:"column:created_by"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	UpdatedBy uint64    `gorm:"column:updated_by"`

	// Custom fields
	TicketsSold      uint64 `gorm:"column:tickets_sold;->"`
	EstimatedRevenue uint64 `gorm:"column:estimated_revenue;->"`
}

func (e Entity) TableName() string {
	return "booking"
}
