package base

import "context"

type IRepository[T any, F Filter] interface {
	Create(ctx context.Context, entity *T) (*T, error)
	CreateMany(ctx context.Context, entities []*T) ([]*T, error)

	GetOne(ctx context.Context, filter F, opts ...QueryOption) (*T, error)
	GetMany(ctx context.Context, filter F, opts ...QueryOption) ([]T, error)
	Count(ctx context.Context, filter F) (int64, error)

	Update(ctx context.Context, entity *T) (*T, error)
	UpdateMany(ctx context.Context, entities []*T) ([]*T, error)

	Delete(ctx context.Context, entity *T) error
}
