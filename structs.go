package xqueue

import (
	"time"

	"github.com/streadway/amqp"
)

type QueueOptions struct {
	Name             string
	Durable          bool
	AutoDeleteUnused bool
	Exclusive        bool
	NoWait           bool
	Arguments        amqp.Table
}

type Message struct {
	QueueName string
	Exchange  string
	Key       string
	Mandatory bool
	Immediate bool

	// Application or exchange specific fields,
	// the headers exchange will inspect this field.
	Headers amqp.Table

	// Properties
	ContentType     string    // MIME content type
	ContentEncoding string    // MIME content encoding
	DeliveryMode    uint8     // Transient (0 or 1) or Persistent (2)
	Priority        uint8     // 0 to 9
	CorrelationId   string    // correlation identifier
	ReplyTo         string    // address to to reply to (ex: RPC)
	Expiration      string    // message expiration spec
	MessageId       string    // message identifier
	Timestamp       time.Time // message timestamp
	Type            string    // message type name
	UserId          string    // creating user id - ex: "guest"
	AppId           string    // creating application id

	// The application specific payload of the message
	Body interface{}
}

type Listener struct {
	QueueName string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Arguments amqp.Table
}
