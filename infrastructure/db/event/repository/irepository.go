package repository

import (
	"context"

	"event_ticket_booking/infrastructure/db/base"
	"event_ticket_booking/infrastructure/db/event/entity"
)

type IRepository interface {
	base.IRepository[entity.Entity, Filter]
	GetList(ctx context.Context, filter Filter, page, size int) ([]entity.Entity, int64, error)
}
