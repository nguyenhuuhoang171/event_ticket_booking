package cron

import (
	"context"
	"log"
	"reflect"

	"event_ticket_booking/config"
	"event_ticket_booking/constant"
	"event_ticket_booking/internal/worker/payment"
	commonModel "event_ticket_booking/model"
)

type CancelExpiredBookingJob struct {
	processor payment.Processor
	enable    bool
	spec      string
}

func NewCancelExpiredBooking(cfg config.Config, lib commonModel.Lib) Job {
	return CancelExpiredBookingJob{
		processor: payment.NewProcessor(cfg, lib),
		enable:    cfg.Cron.CancelExpiredBookingJob.Enable,
		spec:      cfg.Cron.CancelExpiredBookingJob.Spec,
	}
}

func (j CancelExpiredBookingJob) Enabled() bool { return j.enable }

func (j CancelExpiredBookingJob) Name() string { return reflect.TypeOf(j).Name() }

func (j CancelExpiredBookingJob) Spec() string { return j.spec }

func (j CancelExpiredBookingJob) Run(ctx context.Context) error {
	log.Printf("[Cron] %s started", j.Name())

	for {
		bookings, err := j.processor.GetExpiredPendingBookings(ctx)
		if err != nil {
			return err
		}
		if len(bookings) == 0 {
			break
		}

		err = j.processor.FailPayment(ctx, bookings)
		if err != nil {
			return err
		}

		if int64(len(bookings)) < constant.MAX_SIZE {
			break
		}
	}
	return nil
}
