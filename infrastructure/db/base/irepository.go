package base

import "context"

// IRepository is the generic CRUD contract every concrete repository inherits.
// Concrete repos embed this (parameterised with their entity + filter) and add
// only their domain-specific methods.
type IRepository[T any, F Filter] interface {
	// Create inserts one row and returns it with DB-generated fields populated.
	Create(ctx context.Context, entity *T) (*T, error)
	// CreateMany batch-inserts many rows.
	CreateMany(ctx context.Context, entities []*T) ([]*T, error)

	// GetOne returns the first row matching filter, or (nil, nil) if none match.
	GetOne(ctx context.Context, filter F, opts ...QueryOption) (*T, error)
	// GetMany returns every row matching filter (supports pagination/order/...).
	GetMany(ctx context.Context, filter F, opts ...QueryOption) ([]T, error)
	// Count returns the number of rows matching filter.
	Count(ctx context.Context, filter F) (int64, error)

	// Update persists the non-zero fields of entity, keyed by its primary key.
	Update(ctx context.Context, entity *T) (*T, error)
	// UpdateMany bulk-updates rows matching filter and returns rows affected.
	UpdateMany(ctx context.Context, filter F, values map[string]any) (int64, error)
}
