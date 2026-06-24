package usecase

import (
	"context"
	"net/http"

	"event_ticket_booking/constant"
	"event_ticket_booking/internal/domain/event/dto"
	commonModel "event_ticket_booking/model"
)

// GetStats trả về thống kê theo eventId.
// eventId = 0: tất cả event chưa xoá; eventId > 0: một event (404 nếu không tồn tại hoặc đã xoá).
func (u Usecase) GetStats(ctx context.Context, eventId uint64) (*dto.ListStatsResponse, error) {
	stats, err := u.bookingRepo.GetStats(ctx, eventId)
	if err != nil {
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}
	if eventId > 0 && len(stats) == 0 {
		return nil, commonModel.NewError(http.StatusNotFound, "Event not found")
	}

	items := make([]dto.StatsResponse, len(stats))
	for i := range stats {
		items[i] = dto.StatsResponse{
			EventId:          stats[i].EventId,
			TicketsSold:      stats[i].TicketsSold,
			EstimatedRevenue: stats[i].EstimatedRevenue,
		}
	}

	return &dto.ListStatsResponse{
		Items: items,
		Total: int64(len(items)),
	}, nil
}
