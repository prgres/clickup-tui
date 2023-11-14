package clickup

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

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

func (c *Client) requestGet(endpoint string, paramsQuery ...string) ([]byte, error) {
	reqUrl, err := url.Parse(c.apiUrl + endpoint)
	if err != nil {
		return nil, err
	}

	if len(paramsQuery) > 0 {
		params, err := c.parseQueryParams(paramsQuery...)
		if err != nil {
			return nil, err
		}
		reqUrl.RawQuery = params
	}

	c.logger.Infof("requestGet: %s", reqUrl.String())
	req, _ := http.NewRequest("GET", reqUrl.String(), nil)
	req.Header.Add("Authorization", c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	return body, err
}

func (c *Client) parseQueryParams(p ...string) (string, error) {
	if len(p)%2 != 0 {
		return "", fmt.Errorf("invalid number of arguments")
	}

	v := url.Values{}
	for i := 0; i < len(p); i += 2 {
		v.Add(p[i], p[i+1])
	}

	return v.Encode(), nil
}
