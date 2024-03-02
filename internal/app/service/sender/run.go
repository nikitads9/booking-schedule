package sender

import (
	"context"
	"encoding/json"
	"event-schedule/internal/app/model"
	"fmt"
	"log/slog"
	"os"

	"github.com/streadway/amqp"
)

func (s *Service) Run(ctx context.Context) {
	const op = "sender.service.Run"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("sender service initiated")

	msgChan, err := s.rabbitConsumer.Consume()
	if err != nil {
		log.Error("could not get channel to receive messages: ", err)
		os.Exit(1)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgChan:
			s.receiveEvents(msg)

			if err != nil {
				log.Error("could not receive messages: ", err)
			}
			msg.Ack(false)
		}

	}

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
			event.UserID,
			event.CreatedAt,
			event.UpdatedAt,
		))
	}

	return nil
}
