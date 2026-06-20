package repository

import (
	"context"

	"event_ticket_booking/infrastructure/db/base"
	"event_ticket_booking/infrastructure/db/refresh_token/entity"
)

// IRepository is the refresh-token repository contract: the generic CRUD set
// inherited from base, plus the token-specific revoke operation.
type IRepository interface {
	base.IRepository[entity.Entity, entity.Filter]

	// Revoke marks the tokens matching userId OR token as revoked.
	Revoke(ctx context.Context, params entity.Filter) error
}
