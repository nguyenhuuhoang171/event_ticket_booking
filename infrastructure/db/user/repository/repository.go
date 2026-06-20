package repository

import (
	"event_ticket_booking/infrastructure/db/base"
	"event_ticket_booking/infrastructure/db/user/entity"

	"gorm.io/gorm"
)

// Repo is the user repository. It inherits the generic CRUD methods from the
// base repository.
type Repo struct {
	*base.Repository[entity.Entity, entity.Filter]
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repo{
		Repository: base.NewRepository[entity.Entity, entity.Filter](db),
	}
}
