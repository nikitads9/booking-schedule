package rabbit

import (
	"io"

	"booking-schedule/internal/config"

	"github.com/streadway/amqp"
)

const (
	exchangeName = ""
	contentType  = "text/plain"
)

// Producer ...
type Producer interface {
	Publish(msg []byte) error
	Close() error
}

type producer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      *amqp.Queue

	queueName string

	closeFuncs []io.Closer
}

// NewProducer ...
func NewProducer(config *config.RabbitProducer) (Producer, error) {
	closeFuncs := make([]io.Closer, 0, 2)

	conn, err := amqp.Dial(config.DSN)
	if err != nil {
		return nil, err
	}
	closeFuncs = append(closeFuncs, conn)

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	closeFuncs = append(closeFuncs, ch)

	q, err := ch.QueueDeclare(
		config.QueueName, // name
		true,             // durable
		false,            // autodelete
		false,            // exclusive
		false,            // nowait
		nil,              // args
	)
	if err != nil {
		return nil, err
	}

	return &producer{
		connection: conn,
		channel:    ch,
		queue:      &q,
		queueName:  config.QueueName,
		closeFuncs: closeFuncs,
	}, nil
}

// Publish ...
func (p *producer) Publish(msg []byte) error {
	err := p.channel.Publish(
		exchangeName,
		p.queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  contentType,
			Body:         msg,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// Close ...
func (p *producer) Close() error {
	for _, closeFunc := range p.closeFuncs {
		closeFunc.Close() //nolint:errcheck
	}

	return nil
}
