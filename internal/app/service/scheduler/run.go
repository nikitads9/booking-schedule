package scheduler

import (
	"booking-schedule/internal/app/model"
	"booking-schedule/internal/logger/sl"
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func (s *Service) Run(ctx context.Context) {
	const op = "service.scheduler.Run"

	log := s.log.With(
		slog.String("op", op),
	)
	log.Info("scheduler initiated")

	ticker := time.NewTicker(s.checkPeriod)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.handleBookings(ctx)
		}
	}
}

func (s *Service) handleBookings(ctx context.Context) {
	const op = "service.scheduler.handleBookings"
	wg := &sync.WaitGroup{}

	log := s.log.With(
		slog.String("op", op),
	)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	log.Debug("started handling")

	wg.Add(2)

	go func(*sync.WaitGroup) {
		defer wg.Done()
		bookings, err := s.getBookings(ctx)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("failed to get bookings", sl.Err(err))
			return
		}

		if len(bookings) == 0 {
			log.Debug("no bookings to send")
			return
		}

		span.AddEvent("bookings to send acquired", trace.WithAttributes(attribute.Int("quantity", len(bookings))))

		for _, val := range bookings {
			err = s.sendBooking(val)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				log.Error("failed to send booking:", sl.Err(err))
			}
		}
		span.AddEvent("bookings sent")

	}(wg)

	go func(*sync.WaitGroup) {
		defer wg.Done()
		err := s.cleanUpOldBookings(ctx)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("failed to clean up old bookings", sl.Err(err))
			return
		}
		log.Debug("old bookings handled")
	}(wg)

	wg.Wait()
	span.AddEvent("bookings handled")
	log.Debug("finished handling bookings")
}

func (s *Service) getBookings(ctx context.Context) ([]*model.BookingInfo, error) {
	const op = "service.scheduler.getBookings"

	log := s.log.With(
		slog.String("op", op),
	)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	end := time.Now()
	end = time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), end.Minute(), 0, 0, end.Location())
	start := end.Add(-s.checkPeriod)

	bookings, err := s.bookingRepository.GetBookingListByDate(ctx, start, end)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("failed to get list by date", sl.Err(err))
		return nil, err
	}

	return bookings, nil
}

func (s *Service) cleanUpOldBookings(ctx context.Context) error {
	const op = "scheduler.service.cleanUpOldBookings"

	log := s.log.With(
		slog.String("op", op),
	)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	err := s.bookingRepository.DeleteBookingsBeforeDate(ctx, time.Now().Add(-s.bookingTTL))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("failed to clean up old bookings", sl.Err(err))
		return err
	}

	return nil
}

func (s *Service) sendBooking(booking *model.BookingInfo) error {
	data, err := json.Marshal(booking)
	if err != nil {
		return err
	}
	err = s.rabbitProducer.Publish(data)
	if err != nil {
		return err
	}

	return nil
}
