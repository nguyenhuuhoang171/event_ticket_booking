package base

import (
	"context"
	"errors"

	"event_ticket_booking/constant"

	"gorm.io/gorm"
)

var ErrEmptyFilter = errors.New("filter is empty")

const defaultBatchSize = 100

type Filter interface {
	Apply(query *gorm.DB) *gorm.DB
}

type Repository[T any, F Filter] struct {
	Db *gorm.DB
}

func NewRepository[T any, F Filter](db *gorm.DB) *Repository[T, F] {
	return &Repository[T, F]{Db: db}
}

func (r *Repository[T, F]) WithContext(ctx context.Context) *gorm.DB {
	return r.Db.WithContext(ctx)
}

func (r *Repository[T, F]) Create(ctx context.Context, entity *T) (*T, error) {
	if err := r.WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *Repository[T, F]) CreateMany(ctx context.Context, entities []*T) ([]*T, error) {
	if len(entities) == 0 {
		return entities, nil
	}
	if err := r.WithContext(ctx).CreateInBatches(entities, defaultBatchSize).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *Repository[T, F]) GetOne(ctx context.Context, filter F, opts ...QueryOption) (*T, error) {
	if IsEmptyFilter(filter) {
		return nil, ErrEmptyFilter
	}

	query := applyOptions(filter.Apply(r.WithContext(ctx)), opts...)

	var result T
	err := query.First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository[T, F]) GetMany(ctx context.Context, filter F, opts ...QueryOption) ([]T, error) {
	if IsEmptyFilter(filter) {
		return nil, ErrEmptyFilter
	}

	query := applyOptions(filter.Apply(r.WithContext(ctx)), opts...)

	var results []T
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *Repository[T, F]) Count(ctx context.Context, filter F) (int64, error) {
	var count int64
	query := filter.Apply(r.WithContext(ctx).Model(new(T)))

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository[T, F]) Update(ctx context.Context, entity *T) (*T, error) {
	if hasZeroID(entity) {
		return nil, constant.ErrMissingID
	}
	if err := r.WithContext(ctx).Updates(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *Repository[T, F]) UpdateMany(ctx context.Context, entities []*T) ([]*T, error) {
	if len(entities) == 0 {
		return entities, nil
	}
	for _, entity := range entities {
		if hasZeroID(entity) {
			return nil, constant.ErrMissingID
		}
	}

	err := r.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, entity := range entities {
			if err := tx.Updates(entity).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return entities, nil
}

// Delete removes the row identified by entity's primary key. It refuses an
// entity with no primary key to avoid soft-deleting/deleting every row.
func (r *Repository[T, F]) Delete(ctx context.Context, entity *T) error {
	if hasZeroID(entity) {
		return constant.ErrMissingID
	}
	return r.WithContext(ctx).Delete(entity).Error
}
