package sender

import (
	"context"
	"encoding/json"
	"event-schedule/internal/app/model"
	"fmt"
	"log/slog"
	"sync"

	"github.com/streadway/amqp"
)

func (s *Service) Run(ctx context.Context) error {
	const op = "sender.service.Run"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("sender service initiated")

	msgChan, err := s.rabbitConsumer.Consume()
	if err != nil {
		log.Error("could not get channel to receive messages: ", err)
		return err
	}

	/* 	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-msgChan:
			s.receiveEvents(msg)
			if err != nil {
				return err
			}
		}

	} */

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() error {
		defer wg.Done()
		for msg := range msgChan {
			err := s.receiveEvents(msg)
			if err != nil {
				log.Error("could not receive messages: ", err)
				return err
			}
		}

		return nil
	}()

	wg.Wait()

	return nil

}

func (s *Service) receiveEvents(msg amqp.Delivery) error {
	const op = "sender.service.receiveEvents"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info(fmt.Sprintf("Received a message: %s", msg.Body))

	var events []*model.EventInfo
	err := json.Unmarshal(msg.Body, &events)
	if err != nil {
		return err
	}

	for _, event := range events {
		log.Info(fmt.Sprintf(
			"Event:  %d \n "+
				"SuiteID: %d \n "+
				"StartDate: %v \n "+
				"EndDate: :%v \n "+
				"NotifyAt: %v \n "+
				"OwnerID: %d \n "+
				"CreatedAt: %v \n "+
				"UpdatedAt: %v \n\n ",
			event.ID,
			event.SuiteID,
			event.StartDate,
			event.EndDate,
			event.NotifyAt,
			event.OwnerID,
			event.CreatedAt,
			event.UpdatedAt,
		))
	}

	return nil
}
