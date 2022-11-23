package clients

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
)

var accountClient *ApiClient

func InitAccountClient() {
	accountsApiUrl := viper.GetString("SIDOOH_ACCOUNTS_API_URL")
	accountClient = New(accountsApiUrl)
}

func GetAccountClient() *ApiClient {
	return accountClient
}

type AccountApiResponse struct {
	ApiResponse
	Data *Account `json:"data"`
}

func (api *ApiClient) GetOrCreateAccount(phone string) (*Account, error) {
	var apiResponse = new(AccountApiResponse)

	jsonData, err := json.Marshal(map[string]string{"phone": phone})
	dataBytes := bytes.NewBuffer(jsonData)

	err = api.NewRequest(http.MethodPost, "/accounts", dataBytes).Send(apiResponse)
	if err != nil || apiResponse.Result == 0 {
		err = api.NewRequest(http.MethodGet, "/accounts/phone/"+phone, nil).Send(apiResponse)
		if err != nil {
			return nil, err
		}
	}

	return apiResponse.Data, nil
}
