package repository

import (
	"context"
	"errors"

	"event_ticket_booking/constant"
	"event_ticket_booking/infrastructure/db/base"
	"event_ticket_booking/infrastructure/db/booking/entity"
	eventEntity "event_ticket_booking/infrastructure/db/event/entity"

	"gorm.io/gorm"
)

var (
	ErrSoldOut          = errors.New("event is sold out")
	ErrBookingNotFound  = errors.New("booking not found")
	ErrAlreadyCancelled = errors.New("booking already cancelled")
)

type Repo struct {
	*base.Repository[entity.Entity, Filter]
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repo{
		Repository: base.NewRepository[entity.Entity, Filter](db),
	}
}

func (r *Repo) GetList(ctx context.Context, filter Filter, page, size int) ([]entity.Entity, int64, error) {
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

// Reserve: increases sold_tickets + inserts the booking
func (r *Repo) Reserve(ctx context.Context, booking *entity.Entity) (*entity.Entity, error) {
	err := r.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&eventEntity.Entity{}).
			Where("id = ? AND status = ? AND sold_tickets + ? <= total_tickets",
				booking.EventId, constant.EVENT_STATUS_ACTIVE, booking.Quantity).
			UpdateColumn("sold_tickets", gorm.Expr("sold_tickets + ?", booking.Quantity))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrSoldOut
		}
		return tx.Create(booking).Error
	})
	if err != nil {
		return nil, err
	}
	return booking, nil
}

// Cancel: update booking to status CANCELLED + return event's sold_tickets
func (r *Repo) Cancel(ctx context.Context, bookingId, userId uint64) (*entity.Entity, error) {
	var booking entity.Entity
	err := r.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND user_id = ?", bookingId, userId).First(&booking).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBookingNotFound
			}
			return err
		}
		if booking.Status == constant.BOOKING_STATUS_CANCELLED {
			return ErrAlreadyCancelled
		}

		res := tx.Model(&entity.Entity{}).
			Where("id = ? AND status != ?", bookingId, constant.BOOKING_STATUS_CANCELLED).
			Updates(map[string]any{
				"status":     constant.BOOKING_STATUS_CANCELLED,
				"updated_by": userId,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrAlreadyCancelled
		}

		if err := tx.Model(&eventEntity.Entity{}).
			Where("id = ?", booking.EventId).
			UpdateColumn("sold_tickets", gorm.Expr("sold_tickets - ?", booking.Quantity)).Error; err != nil {
			return err
		}

		booking.Status = constant.BOOKING_STATUS_CANCELLED
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &booking, nil
}
