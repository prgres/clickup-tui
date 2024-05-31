package clickup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

const (
	API_URL = "https://api.clickup.com/api/v2"
)

type Client struct {
	httpClient *http.Client
	logger     *slog.Logger
	token      string
	apiUrl     string
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
		logger:     slog.Default(),
	}
}

func NewClient(token string, apiUrl string, logger *slog.Logger) *Client {
	return &Client{
		token:      token,
		httpClient: http.DefaultClient,
		apiUrl:     apiUrl,
		logger:     logger,
	}
}

func NewDefaultClientWithLogger(token string, logger *slog.Logger) *Client {
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

	c.logger.Debug("Sending GET request", "request", reqUrl.String())
	req, err := http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.token)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

func (c *Client) requestPut(endpoint string, data []byte, paramsQuery ...string) ([]byte, error) {
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

	c.logger.Debug("Sending PUT request", "request", reqUrl.String())
	req, err := http.NewRequest("PUT", reqUrl.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.token)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return io.ReadAll(res.Body)
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

type RequestGet interface {
	Error() string
}

func (c *Client) get(url string, objmap RequestGet) error {
	errMsg := "Error occurs while getting resources from url: %s. Error: %s. Raw data: %s"
	errApiMsg := errMsg + " API response: %s"

	rawData, err := c.requestGet(url)
	if err != nil {
		return fmt.Errorf(errMsg, url, err, "none")
	}

	if err := json.Unmarshal(rawData, objmap); err != nil {
		return fmt.Errorf(errApiMsg, url, err, string(rawData))
	}

	if objmap.Error() != "" {
		return fmt.Errorf(
			errMsg, url, "API response contains error.", string(rawData))
	}

	return nil
}

func (c *Client) update(url string, requestUpdate interface{}, objmap interface{}) error {
	errMsg := "Error occurs while getting resources from url: %s. Error: %s. Raw data: %s"
	errApiMsg := errMsg + " API response: %s"

	requestJson, err := json.Marshal(requestUpdate)
	if err != nil {
		return err
	}

	rawData, err := c.requestPut(url, requestJson)
	if err != nil {
		return fmt.Errorf(errMsg, url, err, "none")
	}

	if err := json.Unmarshal(rawData, objmap); err != nil {
		return fmt.Errorf(errApiMsg, url, err, string(rawData))
	}

	return nil
}
