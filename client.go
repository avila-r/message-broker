package xqueue

import (
	"encoding/json"
	"errors"

	"github.com/streadway/amqp"
)

var (
	ErrUnnamedQueue = errors.New("queue needs to be named")
)

type RabbitMQ struct {
	Connection     *amqp.Connection
	DefaultChannel *amqp.Channel
}

func NewClient(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)

	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	response := &RabbitMQ{
		Connection:     conn,
		DefaultChannel: channel,
	}

	return response, nil
}

func (q *RabbitMQ) CloseConnection() {
	q.Connection.Close()
}

func (q *RabbitMQ) CloseDefaultChannel() {
	q.DefaultChannel.Close()
}

func (q *RabbitMQ) NewQueue(options *QueueOptions) error {
	var (
		name        string
		durable     bool       = true
		auto_delete bool       = false
		exclusive   bool       = false
		no_wait     bool       = false
		args        amqp.Table = nil
	)

	if options.Name != "" {
		name = options.Name
	} else {
		return ErrUnnamedQueue
	}

	if v := options.Durable; v {
		durable = v
	}

	if v := options.AutoDeleteUnused; v {
		auto_delete = v
	}

	if v := options.Exclusive; v {
		exclusive = v
	}

	if v := options.NoWait; v {
		no_wait = v
	}

	if a := options.Arguments; a != nil {
		args = a
	}

	_, err := q.DefaultChannel.QueueDeclare(
		name, durable, auto_delete, exclusive, no_wait, args,
	)

	return err
}

func (q *RabbitMQ) Produce(message Message) error {
	// Default values
	var (
		exchange     string = ""
		queue_name   string
		mandatory    bool   = false
		immediate    bool   = false
		content_type string = "text/plain"
	)

	if e := message.Exchange; e != "" {
		exchange = e
	}

	if q := message.QueueName; q != "" {
		queue_name = q
	} else {
		return ErrUnnamedQueue
	}

	if m := message.Mandatory; m {
		mandatory = m
	}

	if i := message.Immediate; i {
		immediate = i
	}

	if c := message.ContentType; c != "" {
		content_type = c
	}

	json, err := json.Marshal(message.Body)

	if err != nil {
		return err
	}

	publishing := amqp.Publishing{
		ContentType: content_type,
		Body:        json,
	}

	return q.DefaultChannel.Publish(
		exchange, queue_name, mandatory, immediate, publishing,
	)
}

func (q *RabbitMQ) NewConsumer(listener *Listener) (<-chan amqp.Delivery, error) {
	var (
		queue_name       string
		consumer         string     = ""
		disable_auto_ack bool       = false
		exclusive        bool       = false
		no_local         bool       = false
		no_wait          bool       = false
		args             amqp.Table = nil
	)

	if n := listener.QueueName; n != "" {
		queue_name = n
	} else {
		return nil, ErrUnnamedQueue
	}

	consumer = listener.Consumer

	if a := listener.DisableAutoAck; a {
		disable_auto_ack = a
	}

	if e := listener.Exclusive; e {
		exclusive = e
	}

	if n := listener.NoLocal; n {
		no_local = n
	}

	if n := listener.NoWait; n {
		no_wait = n
	}

	if a := listener.Arguments; a != nil {
		args = a
	}

	channel, err := q.DefaultChannel.Consume(
		queue_name, consumer, disable_auto_ack, exclusive, no_local, no_wait, args,
	)

	return channel, err
}
