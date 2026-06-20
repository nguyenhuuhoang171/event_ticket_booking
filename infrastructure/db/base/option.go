package base

import "gorm.io/gorm"

// QueryOption tweaks a read query: pagination, ordering, preloading relations,
// selecting columns, ... Pass any number of them to GetOne / GetMany.
type QueryOption func(*gorm.DB) *gorm.DB

func applyOptions(query *gorm.DB, opts ...QueryOption) *gorm.DB {
	for _, opt := range opts {
		query = opt(query)
	}
	return query
}

// WithLimit caps the number of rows returned.
func WithLimit(limit int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if limit > 0 {
			return db.Limit(limit)
		}
		return db
	}
}

// WithOffset skips the first n rows.
func WithOffset(offset int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if offset > 0 {
			return db.Offset(offset)
		}
		return db
	}
}

// WithPage applies limit/offset from a 1-based page number and page size.
func WithPage(page, size int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if size <= 0 {
			return db
		}
		if page < 1 {
			page = 1
		}
		return db.Limit(size).Offset((page - 1) * size)
	}
}

// WithOrder adds an ORDER BY clause, e.g. WithOrder("created_at desc").
func WithOrder(order string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if order != "" {
			return db.Order(order)
		}
		return db
	}
}

// WithPreload eager-loads an association in a separate query (not a SQL JOIN).
func WithPreload(association string, args ...any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(association, args...)
	}
}

// WithJoin adds a raw SQL JOIN. Use it to filter the main entity by a related
// table's columns, e.g.
//
//	WithJoin("JOIN refresh_token rt ON rt.user_id = user.id AND rt.status = ?", 1)
func WithJoin(query string, args ...any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins(query, args...)
	}
}

// WithGroup adds a GROUP BY clause (handy together with WithJoin).
func WithGroup(group string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if group != "" {
			return db.Group(group)
		}
		return db
	}
}

// WithDistinct selects distinct rows (useful when a JOIN duplicates main rows).
func WithDistinct() QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Distinct()
	}
}

// WithSelect restricts the columns fetched.
func WithSelect(columns ...string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if len(columns) > 0 {
			return db.Select(columns)
		}
		return db
	}
}
