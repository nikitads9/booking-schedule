package rabbit

import (
	"booking-schedule/internal/config"
	"io"

	"github.com/streadway/amqp"
)

const consumerName = "sender"

type Consumer interface {
	Consume() (<-chan amqp.Delivery, error)
	Close() error
}

type consumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      *amqp.Queue
	queueName  string

	closeFuncs []io.Closer
}

func NewConsumer(config *config.RabbitConsumer) (Consumer, error) {
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

	return &consumer{
		connection: conn,
		channel:    ch,
		queue:      &q,
		queueName:  config.QueueName,
		closeFuncs: closeFuncs,
	}, nil
}

func (c *consumer) Consume() (<-chan amqp.Delivery, error) {
	msgChan, err := c.channel.Consume(
		c.queueName,
		consumerName,
		false, // autoAck
		true,  // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgChan, nil
}

func (c *consumer) Close() error {
	for _, closeFunc := range c.closeFuncs {
		err := closeFunc.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
