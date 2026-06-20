package base

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// ErrEmptyFilter is returned when a query that requires a filter is called with
// a zero-value filter, to avoid accidentally matching every row.
var ErrEmptyFilter = errors.New("filter is empty")

// defaultBatchSize is the chunk size used by CreateMany.
const defaultBatchSize = 100

// Filter is implemented by every entity's Filter struct. Apply translates the
// filter into gorm Where/Join/Group clauses so the base repository stays
// entity-agnostic.
type Filter interface {
	Apply(query *gorm.DB) *gorm.DB
}

// Repository is a generic, reusable CRUD repository over a gorm model T whose
// queries are filtered by F. Concrete repositories embed it and only add the
// methods specific to their domain.
type Repository[T any, F Filter] struct {
	Db *gorm.DB
}

// NewRepository builds a base repository for the model T / filter F pair.
func NewRepository[T any, F Filter](db *gorm.DB) *Repository[T, F] {
	return &Repository[T, F]{Db: db}
}

// WithContext returns the underlying *gorm.DB bound to ctx.
func (r *Repository[T, F]) WithContext(ctx context.Context) *gorm.DB {
	return r.Db.WithContext(ctx)
}

// Create inserts one row and returns it with DB-generated fields populated.
func (r *Repository[T, F]) Create(ctx context.Context, entity *T) (*T, error) {
	if err := r.WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

// CreateMany batch-inserts many rows (chunked to keep statements reasonable).
func (r *Repository[T, F]) CreateMany(ctx context.Context, entities []*T) ([]*T, error) {
	if len(entities) == 0 {
		return entities, nil
	}
	if err := r.WithContext(ctx).CreateInBatches(entities, defaultBatchSize).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// GetOne returns the first row matching filter, or (nil, nil) when none match.
// It refuses an empty filter to avoid returning an arbitrary row.
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

// GetMany returns every row matching filter. It refuses an empty filter to
// avoid accidentally scanning the whole table; use options for pagination,
// ordering and preloading.
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

// Count returns the number of rows matching filter.
func (r *Repository[T, F]) Count(ctx context.Context, filter F) (int64, error) {
	var count int64
	query := filter.Apply(r.WithContext(ctx).Model(new(T)))

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Update persists the non-zero fields of entity, keyed by its primary key
// (gorm Updates semantics).
func (r *Repository[T, F]) Update(ctx context.Context, entity *T) (*T, error) {
	if err := r.WithContext(ctx).Updates(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

// UpdateMany bulk-updates every row matching filter with the given column
// values and returns the number of affected rows. It refuses an empty filter.
func (r *Repository[T, F]) UpdateMany(ctx context.Context, filter F, values map[string]any) (int64, error) {
	if IsEmptyFilter(filter) {
		return 0, ErrEmptyFilter
	}
	if len(values) == 0 {
		return 0, nil
	}

	result := filter.Apply(r.WithContext(ctx).Model(new(T))).Updates(values)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
