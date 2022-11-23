package clients

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestInitAccountClient(t *testing.T) {
	viper.Set("SIDOOH_ACCOUNTS_API_URL", "test.test")
	InitAccountClient()

	client := GetAccountClient()

	assert.NotNil(t, t, client)
	assert.NotNil(t, t, client.client)
	assert.NotNil(t, t, client.request)
	assert.NotNil(t, t, client.cache)

	assert.Equal(t, "test.test", client.baseUrl)
}

func accountFoundRequest(t *testing.T) RoundTripFunc {
	return func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "http://localhost:8000/api/v1/accounts")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(strings.NewReader("{\"result\":1,\"data\":{\"id\":6,\"phone\":\"254110039317\",\"active\":true,\"inviter_id\":1,\"user_id\":50}}")),
			// Must be set to non-nil value, or it panics
			//Header: make(http.Header),
		}
	}
}

func accountNotFoundRequest(t *testing.T) RoundTripFunc {
	return func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "http://localhost:8000/api/v1/accounts")
		return &http.Response{
			StatusCode: 500,
			// Send response to be tested
			Body: ioutil.NopCloser(strings.NewReader("{\"result\":0,\"message\":\"Something went wrong, please try again.\"}")),
			// Must be set to non-nil value, or it panics
			//Header: make(http.Header),
		}
	}
}

func TestApiClient_GetOrCreateAccount(t *testing.T) {
	client.baseUrl = "http://localhost:8000"
	viper.Set("SIDOOH_ACCOUNTS_API_URL", client.baseUrl)

	initTestClient(accountFoundRequest(t))

	account, err := GetAccountClient().GetOrCreateAccount("")
	fmt.Println(account, err)
	if err != nil {
		return
	}

}
