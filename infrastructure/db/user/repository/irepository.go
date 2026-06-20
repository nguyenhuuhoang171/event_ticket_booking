package repository

import (
	"event_ticket_booking/infrastructure/db/base"
	"event_ticket_booking/infrastructure/db/user/entity"
)

// IRepository is the user repository contract: the generic CRUD set inherited
// from base, plus any user-specific methods added below.
type IRepository interface {
	base.IRepository[entity.Entity, entity.Filter]
}
