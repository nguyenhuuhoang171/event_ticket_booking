package repository

import (
	"context"

	"event_ticket_booking/constant"
	"event_ticket_booking/infrastructure/db/base"
	"event_ticket_booking/infrastructure/db/event/entity"

	"gorm.io/gorm"
)

type Repo struct {
	*base.Repository[entity.Entity, entity.Filter]
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repo{
		Repository: base.NewRepository[entity.Entity, entity.Filter](db),
	}
}

func (r *Repo) GetList(ctx context.Context, filter entity.Filter, page, size int) ([]entity.Entity, int64, error) {
	var total int64
	countQuery := filter.Apply(r.WithContext(ctx).Model(&entity.Entity{}))
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if size <= 0 {
		size = constant.DEFAULT_PAGE_SIZE
	}
	if page < 1 {
		page = 1
	}

	var items []entity.Entity
	listQuery := filter.Apply(r.WithContext(ctx).Model(&entity.Entity{}))
	if err := listQuery.
		Order(entity.Entity{}.TableName() + ".id DESC").
		Limit(size).
		Offset((page - 1) * size).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
