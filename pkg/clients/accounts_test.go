package clients

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitAccountClient(t *testing.T) {
	viper.Set("SIDOOH_ACCOUNTS_API_URL", "test.test")

	client := InitAccountClient()

	assert.NotNil(t, t, client)
	assert.NotNil(t, t, client.client)
	assert.NotNil(t, t, client.request)
	assert.NotNil(t, t, client.cache)

	assert.Equal(t, "test.test", client.baseUrl)
}
