package xqueue_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func Test_Connection(t *testing.T) {
	assert := assert.New(t)

	if err := godotenv.Load(".env.test"); err != nil {
		t.Errorf("failed while loading '.env.test' file - %v", err)
	}

	url := os.Getenv("QUEUE_SERVER_URL")

	assert.NotEmpty(url)

}
