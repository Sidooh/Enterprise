package clients

import (
	"github.com/spf13/viper"
)

func InitPaymentClient() *ApiClient {
	apiUrl := viper.GetString("SIDOOH_PAYMENTS_API_URL")
	return New(apiUrl)
}
