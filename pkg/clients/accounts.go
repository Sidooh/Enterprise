package clients

import (
	"bytes"
	"encoding/json"
	"enterprise.sidooh/pkg/services"
	"github.com/spf13/viper"
	"net/http"
)

func InitAccountClient() *ApiClient {
	accountsApiUrl := viper.GetString("SIDOOH_ACCOUNTS_API_URL")
	return New(accountsApiUrl)
}

func (api *ApiClient) GetOrCreateAccount(phone string) (*services.Account, error) {
	var apiResponse = new(services.AccountApiResponse)

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
