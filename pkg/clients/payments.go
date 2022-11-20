package clients

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
)

var paymentClient *ApiClient

func InitPaymentClient() {
	apiUrl := viper.GetString("SIDOOH_PAYMENTS_API_URL")
	paymentClient = New(apiUrl)
}

func GetPaymentClient() *ApiClient {
	return paymentClient
}

func (api *ApiClient) CreateFloatAccount(enterpriseId, accountId int) (*FloatAccount, error) {
	var apiResponse = new(FloatAccountApiResponse)

	jsonData, err := json.Marshal(map[string]string{
		"initiator":  "ENTERPRISE",
		"reference":  strconv.Itoa(enterpriseId),
		"account_id": strconv.Itoa(accountId),
	})
	dataBytes := bytes.NewBuffer(jsonData)

	err = api.NewRequest(http.MethodPost, "/float-accounts", dataBytes).Send(apiResponse)

	return apiResponse.Data, err
}
