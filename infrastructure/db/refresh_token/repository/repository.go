package repository

import (
	"context"
	"errors"

	"event_ticket_booking/constant"
	"event_ticket_booking/infrastructure/db/base"
	"event_ticket_booking/infrastructure/db/refresh_token/entity"

	"gorm.io/gorm"
)

// Repo is the refresh-token repository. It inherits the generic CRUD methods
// from the base repository and adds the token-specific revoke operations.
type Repo struct {
	*base.Repository[entity.Entity, entity.Filter]
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repo{
		Repository: base.NewRepository[entity.Entity, entity.Filter](db),
	}
}

// Revoke marks the tokens matching userId or token as revoked.
func (r *Repo) Revoke(ctx context.Context, params entity.Filter) error {
	if params.UserId == 0 && params.Token == "" {
		return errors.New("filter is empty")
	}
	result := r.WithContext(ctx).
		Model(&entity.Entity{}).
		Where("user_id = ? or token = ?", params.UserId, params.Token).
		Update("status", constant.REFRESH_TOKEN_STATUS_REVOKED)
	return result.Error
}
