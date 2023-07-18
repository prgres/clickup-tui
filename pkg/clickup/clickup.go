package clickup

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	logger "github.com/prgrs/clickup/pkg/logger1"
)

const (
	API_URL = "https://api.clickup.com/api/v2"
)

type Client struct {
	token      string
	httpClient *http.Client
	apiUrl     string

	logger logger.Logger
}

func (c *Client) ToJson(data interface{}) string {
	json := c.ToJsonByte(data)
	return string(json)
}

func (c *Client) ToJsonByte(data interface{}) []byte {
	json, _ := json.MarshalIndent(data, "", "    ")
	return json
}

func NewDefaultClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: http.DefaultClient,
		apiUrl:     API_URL,
		logger:     logger.NewDefaultLogger(),
	}
}

func NewClient(token string, apiUrl string, logger logger.Logger) *Client {
	return &Client{
		token:      token,
		httpClient: http.DefaultClient,
		apiUrl:     apiUrl,
		logger:     logger,
	}
}
func NewDefaultClientWithLogger(token string, logger logger.Logger) *Client {
	return &Client{
		token:      token,
		httpClient: http.DefaultClient,
		apiUrl:     API_URL,
		logger:     logger,
	}
}

func (c *Client) requestGet(endpoint string) ([]byte, error) {
	reqUrl := c.apiUrl + endpoint
	req, _ := http.NewRequest("GET", reqUrl, nil)
	req.Header.Add("Authorization", c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	return body, err
}
