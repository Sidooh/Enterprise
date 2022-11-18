package clients

import (
	"github.com/spf13/viper"
)

func InitAccountClient() *ApiClient {
	accountsApiUrl := viper.GetString("SIDOOH_ACCOUNTS_API_URL")
	return New(accountsApiUrl)
}
