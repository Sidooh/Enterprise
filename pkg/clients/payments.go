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

type VoucherTypesApiResponse struct {
	ApiResponse
	Data *[]VoucherType `json:"data"`
}

type VoucherTypeApiResponse struct {
	ApiResponse
	Data *VoucherType `json:"data"`
}

type FloatAccountApiResponse struct {
	ApiResponse
	Data *FloatAccount `json:"data"`
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

func (api *ApiClient) FetchVoucherTypes(accountId int) (*[]VoucherType, error) {
	var apiResponse = new(VoucherTypesApiResponse)

	err := api.NewRequest(http.MethodGet, "/voucher-types?account_id="+strconv.Itoa(accountId), nil).Send(apiResponse)

	return apiResponse.Data, err
}

func (api *ApiClient) FetchVoucherType(accountId, voucherTypeId int) (*VoucherType, error) {
	var apiResponse = new(VoucherTypeApiResponse)

	var endpoint = "/voucher-types/" + strconv.Itoa(voucherTypeId) + "?account_id=" + strconv.Itoa(accountId) + "&with_vouchers=true"
	err := api.NewRequest(http.MethodGet, endpoint, nil).Send(apiResponse)

	return apiResponse.Data, err
}

func (api *ApiClient) CreateVoucherType(accountId int, name string) (*VoucherType, error) {
	var apiResponse = new(ApiResponse)

	jsonData, err := json.Marshal(map[string]string{
		"initiator":  "ENTERPRISE",
		"name":       name,
		"account_id": strconv.Itoa(accountId),
	})
	dataBytes := bytes.NewBuffer(jsonData)

	err = api.NewRequest(http.MethodPost, "/voucher-types", dataBytes).Send(apiResponse)

	return apiResponse.Data.(*VoucherType), err
}

func (api *ApiClient) FetchFloatAccount(accountId int) (*FloatAccount, error) {
	var apiResponse = new(FloatAccountApiResponse)

	var endpoint = "http://localhost:8002/api/v1/float-accounts/" + strconv.Itoa(accountId)
	err := api.NewRequest(http.MethodGet, endpoint, nil).Send(apiResponse)

	return apiResponse.Data, err
}
