package xqueue

import (
	"errors"

	"github.com/streadway/amqp"
)

func NewConnection(url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}

func NewChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	return conn.Channel()
}

func CloseConnection(conn *amqp.Connection) error {
	return conn.Close()
}

func CloseChannel(chann *amqp.Channel) error {
	return chann.Close()
}

func DeclareNewQueue(channel *amqp.Channel, options *QueueOptions) error {
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
		return errors.New("queue needs to be named")
	}

	if options.Durable {
		durable = options.Durable
	}

	if options.AutoDeleteUnused {
		auto_delete = options.AutoDeleteUnused
	}

	if options.Exclusive {
		exclusive = options.Exclusive
	}

	if options.NoWait {
		no_wait = options.NoWait
	}

	if options.Arguments != nil {
		args = options.Arguments
	}

	_, err := channel.QueueDeclare(
		name, durable, auto_delete, exclusive, no_wait, args,
	)

	return err
}

func GetNewQueue(channel *amqp.Channel, options *QueueOptions) (amqp.Queue, error) {
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
		return amqp.Queue{}, errors.New("queue needs to be named")
	}

	if options.Durable {
		durable = options.Durable
	}

	if options.AutoDeleteUnused {
		auto_delete = options.AutoDeleteUnused
	}

	if options.Exclusive {
		exclusive = options.Exclusive
	}

	if options.NoWait {
		no_wait = options.NoWait
	}

	if options.Arguments != nil {
		args = options.Arguments
	}

	return channel.QueueDeclare(
		name, durable, auto_delete, exclusive, no_wait, args,
	)
}
