package repository

import (
	"context"

	"event_ticket_booking/infrastructure/db/base"
	"event_ticket_booking/infrastructure/db/booking/entity"
)

type IRepository interface {
	base.IRepository[entity.Entity, Filter]
	Reserve(ctx context.Context, booking *entity.Entity) (*entity.Entity, error)
	Cancel(ctx context.Context, bookingId, userId uint64) (*entity.Entity, error)
	GetList(ctx context.Context, filter Filter, page, size int) ([]entity.Entity, int64, error)
}
