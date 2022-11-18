package clients

import (
	"github.com/spf13/viper"
)

func InitNotifyClient() *ApiClient {
	apiUrl := viper.GetString("SIDOOH_NOTIFY_API_URL")
	return New(apiUrl)
}
