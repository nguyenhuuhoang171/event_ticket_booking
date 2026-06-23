package worker

import (
	"context"

	"event_ticket_booking/config"
	"event_ticket_booking/internal/worker/consumer"
	"event_ticket_booking/internal/worker/cron"
	commonModel "event_ticket_booking/model"
)

type Worker struct {
	consumerRunner *consumer.Runner
}

func InitWorker(ctx context.Context, cfg config.Config, lib commonModel.Lib) (*Worker, error) {
	consumerRunner, err := initConsumerRunner(ctx, cfg, lib)
	if err != nil {
		return nil, err
	}

	if err := initCronRunner(ctx, cfg, lib); err != nil {
		return nil, err
	}

	return &Worker{consumerRunner: consumerRunner}, nil
}

func (w *Worker) Stop() error {
	return w.consumerRunner.Close()
}

func initConsumerRunner(ctx context.Context, cfg config.Config, lib commonModel.Lib) (*consumer.Runner, error) {
	runner := consumer.NewRunner(cfg.Kafka,
		consumer.NewPaymentConsumer(cfg, lib),
	)
	if err := runner.Start(ctx); err != nil {
		return nil, err
	}
	return runner, nil
}

func initCronRunner(ctx context.Context, cfg config.Config, lib commonModel.Lib) error {
	runner := cron.NewRunner(
		cron.NewCancelExpiredBooking(cfg, lib),
	)
	return runner.Start(ctx)
}
