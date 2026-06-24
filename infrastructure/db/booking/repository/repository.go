package repository

import (
	"context"
	"errors"

	"event_ticket_booking/constant"
	"event_ticket_booking/infrastructure/db/base"
	"event_ticket_booking/infrastructure/db/booking/entity"
	eventEntity "event_ticket_booking/infrastructure/db/event/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repo struct {
	*base.Repository[entity.Entity, Filter]
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repo{
		Repository: base.NewRepository[entity.Entity, Filter](db),
	}
}

func (r *Repo) GetListPaging(ctx context.Context, filter Filter, page, size int64) ([]entity.Entity, int64, error) {
	var total int64
	countQuery := filter.Apply(r.WithContext(ctx).Model(&entity.Entity{}))
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if size <= 0 {
		size = constant.DEFAULT_SIZE
	}
	if page < 1 {
		page = 1
	}

	var items []entity.Entity
	listQuery := filter.Apply(r.WithContext(ctx).Model(&entity.Entity{}))
	if err := listQuery.
		Order(entity.Entity{}.TableName() + ".id DESC").
		Limit(int(size)).
		Offset(int((page - 1) * size)).
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
			return constant.ErrSoldOut
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
				return constant.ErrBookingNotFound
			}
			return err
		}
		if booking.Status == constant.BOOKING_STATUS_CANCELLED {
			return constant.ErrAlreadyCancelled
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
			return constant.ErrAlreadyCancelled
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

// Confirm: chuyển booking PENDING sang CONFIRMED. Trả ErrNotPending nếu booking không còn PENDING.
func (r *Repo) Confirm(ctx context.Context, bookingId uint64) (*entity.Entity, error) {
	var booking entity.Entity
	err := r.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", bookingId).First(&booking).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return constant.ErrBookingNotFound
			}
			return err
		}
		if booking.Status != constant.BOOKING_STATUS_PENDING {
			return constant.ErrNotPending
		}

		res := tx.Model(&entity.Entity{}).
			Where("id = ? AND status = ?", bookingId, constant.BOOKING_STATUS_PENDING).
			Updates(map[string]any{
				"status":     constant.BOOKING_STATUS_CONFIRMED,
				"updated_by": constant.SYSTEM_USER_ID,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return constant.ErrNotPending
		}

		booking.Status = constant.BOOKING_STATUS_CONFIRMED
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// Cancel các booking PENDING (hết hạn thanh toán) + trả lại sold_tickets,
// tất cả trong 1 transaction. Trả về đúng những booking thực sự bị huỷ (để caller trả vé Redis).
// Dùng SELECT ... FOR UPDATE để khoá đúng tập PENDING, tránh huỷ nhầm booking vừa được CONFIRMED.
func (r *Repo) CancelBookings(ctx context.Context, bookingIds []uint64) ([]entity.Entity, error) {
	if len(bookingIds) == 0 {
		return nil, nil
	}

	var cancelled []entity.Entity
	err := r.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Khoá + lấy các booking còn PENDING trong danh sách
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id IN ? AND status = ?", bookingIds, constant.BOOKING_STATUS_PENDING).
			Find(&cancelled).Error; err != nil {
			return err
		}
		if len(cancelled) == 0 {
			return nil
		}

		ids := make([]uint64, len(cancelled))
		for i := range cancelled {
			ids[i] = cancelled[i].Id
		}

		// CancelMany
		if err := tx.Model(&entity.Entity{}).
			Where("id IN ?", ids).
			Updates(map[string]any{
				"status":     constant.BOOKING_STATUS_CANCELLED,
				"updated_by": constant.SYSTEM_USER_ID,
			}).Error; err != nil {
			return err
		}

		// Trả lại sold_tickets, gộp theo event
		eventQty := make(map[uint64]uint64)
		for i := range cancelled {
			eventQty[cancelled[i].EventId] += cancelled[i].Quantity
		}
		for eventId, qty := range eventQty {
			if err := tx.Model(&eventEntity.Entity{}).
				Where("id = ?", eventId).
				UpdateColumn("sold_tickets", gorm.Expr("sold_tickets - ?", qty)).Error; err != nil {
				return err
			}
		}

		for i := range cancelled {
			cancelled[i].Status = constant.BOOKING_STATUS_CANCELLED
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return cancelled, nil
}

// GetStats trả về thống kê vé đã bán và doanh thu (chỉ từ booking CONFIRMED).
// eventId = 0: tất cả event chưa xoá; eventId > 0: lọc theo 1 event.
func (r *Repo) GetStats(ctx context.Context, eventId uint64) ([]entity.Entity, error) {
	query := r.WithContext(ctx).
		Table("event e").
		Select(`
			e.id AS event_id,
			COALESCE(SUM(b.quantity), 0) AS tickets_sold,
			COALESCE(SUM(b.quantity), 0) * e.ticket_price AS estimated_revenue
		`).
		Joins("LEFT JOIN booking b ON b.event_id = e.id AND b.status = ?", constant.BOOKING_STATUS_CONFIRMED).
		Where("e.status = ?", constant.EVENT_STATUS_ACTIVE)

	if eventId > 0 {
		query = query.Where("e.id = ?", eventId)
	}

	var results []entity.Entity
	err := query.
		Group("e.id, e.ticket_price").
		Order("e.id ASC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
