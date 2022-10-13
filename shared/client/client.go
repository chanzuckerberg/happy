package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type HappyClient struct {
	client        http.Client
	apiBaseUrl    string
	clientName    string
	clientVersion string
}

func NewHappyClient(clientName, clientVersion, apiBaseUrl string) *HappyClient {
	return &HappyClient{
		apiBaseUrl:    apiBaseUrl,
		clientName:    clientName,
		clientVersion: clientVersion,
		client:        http.Client{},
	}
}

func (c *HappyClient) Get(route string, body interface{}) (*http.Response, error) {
	return c.makeRequest(http.MethodGet, route, body)
}

func (c *HappyClient) GetParsed(route string, body, result interface{}, notFoundMessage string) error {
	resp, err := c.Get(route, body)
	if err != nil {
		return errors.Wrap(err, "request failed with")
	}

	return c.parseResponse(resp, result, notFoundMessage)
}

func (c *HappyClient) Post(route string, body interface{}) (*http.Response, error) {
	return c.makeRequest(http.MethodPost, route, body)
}
func (c *HappyClient) PostParsed(route string, body, result interface{}, notFoundMessage string) error {
	resp, err := c.Post(route, body)
	if err != nil {
		return errors.Wrap(err, "request failed with")
	}

	return c.parseResponse(resp, result, notFoundMessage)
}

func (c *HappyClient) Delete(route string, body interface{}) (*http.Response, error) {
	return c.makeRequest(http.MethodDelete, route, body)
}
func (c *HappyClient) DeleteParsed(route string, body, result interface{}, notFoundMessage string) error {
	resp, err := c.Delete(route, body)
	if err != nil {
		return errors.Wrap(err, "request failed with")
	}

	return c.parseResponse(resp, result, notFoundMessage)
}

func (c *HappyClient) parseResponse(resp *http.Response, result interface{}, notFoundMessage string) error {
	err := InspectForErrors(resp, notFoundMessage)
	if err != nil {
		return err
	}

	ParseResponse(resp, &result)
	return nil
}

func (c *HappyClient) makeRequest(method, route string, body interface{}) (*http.Response, error) {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert body to json")
	}
	bodyReader := bytes.NewReader(bodyJson)
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.apiBaseUrl, route), bodyReader)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *HappyClient) Do(req *http.Request) (*http.Response, error) {
	if c.clientVersion != "undefined" {
		req.Header.Add("User-Agent", fmt.Sprintf("%s/%s", c.clientName, c.clientVersion))
	}
	req.Header.Add("Content-Type", "application/json")

	// fmt.Println("headers:", req.Header)

	// TODO: add auth header

	return c.client.Do(req)
}
