package clients

import (
	"bytes"
	"encoding/json"
	"enterprise.sidooh/pkg/cache"
	"enterprise.sidooh/pkg/logger"
	"errors"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ApiClient struct {
	client  *http.Client
	request *http.Request
	baseUrl string
	cache   cache.ICache[string, string]
}

type AuthResponse struct {
	Token string `json:"access_token"`
}

type ApiResponse struct {
	Result  int         `json:"result"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

var clientCache cache.ICache[string, string]

func Init() {
	logger.ClientLog.Println("Init client")

	clientCache = cache.New[string, string]()
}

func New(baseUrl string) *ApiClient {
	logger.ClientLog.Println("New client: ", baseUrl)

	return &ApiClient{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseUrl: baseUrl,
		cache:   clientCache,
	}
}

func (api *ApiClient) getUrl(endpoint string) string {
	if strings.HasPrefix(endpoint, "http") {
		return endpoint
	}
	if !strings.HasPrefix(api.baseUrl, "http") {
		api.baseUrl = "https://" + api.baseUrl
	}
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}
	return api.baseUrl + endpoint
}

func (api *ApiClient) Send(data interface{}) error {
	//TODO: Can we encode the data for security purposes and decode when necessary? Same to response logging...
	logger.ClientLog.Println("API_REQ: ", api.request)
	start := time.Now()
	response, err := api.client.Do(api.request)
	if err != nil {
		logger.ClientLog.Error("Error sending request to API endpoint: ", err)
		return err
	}
	// Close the connection to reuse it
	defer response.Body.Close()
	logger.ClientLog.Println("API_RES - raw: ", response, time.Since(start))

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.ClientLog.Println("Couldn't parse response body: ", err)
	}
	logger.ClientLog.Println("API_RES - body: ", string(body))

	//TODO: Perform error handling in a better way
	if response.StatusCode != 200 && response.StatusCode != 201 && response.StatusCode != 401 &&
		response.StatusCode != 404 && response.StatusCode != 422 {
		if response.StatusCode < 500 {
			var errorMessage map[string][]map[string]string
			err = json.Unmarshal(body, &errorMessage)

			if len(errorMessage["errors"]) == 0 {
				var errorMessage map[string]string
				err = json.Unmarshal(body, &errorMessage)
				logger.ClientLog.Println("API_ERR - body: ", errorMessage)

				return errors.New(errorMessage["message"])
			}

			return errors.New(errorMessage["errors"][0]["message"])
		}

		return errors.New(string(body))
	}

	if response.StatusCode == 404 {
		return errors.New(string(body))
	}

	//TODO: Deal with 401
	if response.StatusCode == 401 {
		logger.ClientLog.Panic("Failed to authenticate.")
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.ClientLog.Println("Failed to unmarshal body: ", err)
	}

	return nil
}

func (api *ApiClient) setDefaultHeaders() {
	api.request.Header.Add("Accept", "application/json")
	api.request.Header.Add("Content-Type", `application/json`)
}

func (api *ApiClient) baseRequest(method string, endpoint string, body io.Reader) *ApiClient {
	endpoint = api.getUrl(endpoint)
	request, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		logger.ClientLog.Fatalf("error creating HTTP request: %v", err)
	}

	api.request = request
	api.setDefaultHeaders()

	return api
}

func (api *ApiClient) NewRequest(method string, endpoint string, body io.Reader) *ApiClient {
	if token := api.cache.GetString("token"); token != "" {
		// TODO: Check if token has expired since we should be able to decode it
		api.baseRequest(method, endpoint, body).request.Header.Add("Authorization", "Bearer "+token)
	} else {
		api.ensureAuthenticated()

		//TODO: What will happen to client if cache fails to store token? E.g. when account srv is not reachable?
		// TODO: Can we even just use a global Var?
		token = api.cache.GetString("token")
		api.baseRequest(method, endpoint, body).request.Header.Add("Authorization", "Bearer "+token)
	}

	return api
}

func (api *ApiClient) ensureAuthenticated() {
	values := map[string]string{"email": "aa@a.a", "password": "12345678"}
	jsonData, err := json.Marshal(values)

	err = api.authenticate(jsonData)
	if err != nil {
		logger.ClientLog.Printf("error authenticating: %v", err)
	}
}

func (api *ApiClient) authenticate(data []byte) error {
	var response = new(AuthResponse)

	err := api.baseRequest(http.MethodPost, viper.GetString("SIDOOH_ACCOUNTS_API_URL")+"/users/signin", bytes.NewBuffer(data)).Send(response)
	if err != nil {
		return err
	}

	if api.cache != nil {
		api.cache.Set("token", response.Token, 14*time.Minute)
	}

	return nil
}
