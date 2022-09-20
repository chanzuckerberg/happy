package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chanzuckerberg/happy/pkg/config"
	"github.com/chanzuckerberg/happy/pkg/util"
	"github.com/pkg/errors"
)

type HappyClient struct {
	happyConfig *config.HappyConfig
	client      http.Client
}

func NewHappyClient(happyConfig *config.HappyConfig) *HappyClient {
	return &HappyClient{
		happyConfig: happyConfig,
		client:      http.Client{},
	}
}

func (c *HappyClient) Get(route string, body interface{}) (*http.Response, error) {
	return c.makeRequest(http.MethodGet, route, body)
}

func (c *HappyClient) Post(route string, body interface{}) (*http.Response, error) {
	return c.makeRequest(http.MethodPost, route, body)
}

func (c *HappyClient) Delete(route string, body interface{}) (*http.Response, error) {
	return c.makeRequest(http.MethodDelete, route, body)
}

func (c *HappyClient) makeRequest(method, route string, body interface{}) (*http.Response, error) {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert body to json")
	}
	bodyReader := bytes.NewReader(bodyJson)
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.happyConfig.GetHappyApiBaseUrl(), route), bodyReader)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *HappyClient) Do(req *http.Request) (*http.Response, error) {
	if util.GetVersion().Version != "undefined" {
		req.Header.Add("User-Agent", fmt.Sprintf("happy-cli/%s", util.GetVersion().Version))
	}
	req.Header.Add("Content-Type", "application/json")

	// fmt.Println("headers:", req.Header)

	// TODO: add auth header

	return c.client.Do(req)
}
