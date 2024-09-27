package xqueue_test

import (
	"encoding/json"
	"os"
	"sync"
	"testing"

	xqueue "github.com/avila-r/message-broker"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	queue_name = "test-queue"
)

func Test_Connection(t *testing.T) {
	assert := assert.New(t)

	if err := godotenv.Load(".env.test"); err != nil {
		t.Errorf("failed while loading '.env.test' file - %v", err)
	}

	// Get RabbitMQ URL
	url := os.Getenv("QUEUE_SERVER_URL")

	assert.NotEmpty(url)

	// Connect to RabbitMQ
	rabbitmq, err := xqueue.NewClient(url)

	if err != nil {
		t.Errorf("failed to connect in rabbitmq - %v", err)
	}

	t.Log("Successfully connected to RabbitMQ")

	// Declare new queue
	rabbitmq.NewQueue(&xqueue.QueueOptions{
		Name: queue_name,
	})

	// Create new consumer and listen
	// by messages in created queue
	//
	// OBS: This channel is receive-only (<-chan)
	channel, _ := rabbitmq.NewConsumer(&xqueue.Listener{
		QueueName: queue_name,
	})

	// Declare wait group
	var w sync.WaitGroup

	w.Add(1)

	// Produce a message
	rabbitmq.Produce(xqueue.Message{
		QueueName: queue_name,
		Body:      "your payload",
	})

	// Create goroutine to receive channel
	go func() {
		// Finish wait group
		defer w.Done()

		// Receive message from consumer channel
		message := <-channel

		// Declare variable to receive unmarshalled payload
		var payload any

		// Unmarshal payload's bytes to 'payload' variable
		if err := json.Unmarshal(message.Body, &payload); err != nil {
			t.Errorf("error while unmarshalling - %v", err)
		}

		t.Logf("payload received: %v", payload)
	}()

	// Wait for the message to be received at goroutine
	w.Wait()
}
