package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type TokenProvider interface {
	GetToken() (string, error)
}

type HappyClient struct {
	client        http.Client
	apiBaseUrl    string
	clientName    string
	clientVersion string
	tokenProvider TokenProvider
}

func NewHappyClient(clientName, clientVersion, apiBaseUrl string, tokenProvider TokenProvider) *HappyClient {
	return &HappyClient{
		apiBaseUrl:    apiBaseUrl,
		clientName:    clientName,
		clientVersion: clientVersion,
		client:        http.Client{},
		tokenProvider: tokenProvider,
	}
}

func (c *HappyClient) Get(route string, body interface{}) (*http.Response, error) {
	return c.makeRequest(http.MethodGet, route, body)
}

func (c *HappyClient) GetParsed(route string, body, result interface{}) error {
	resp, err := c.Get(route, body)
	if err != nil {
		return errors.Wrap(err, "request failed with")
	}

	return c.parseResponse(resp, result)
}

func (c *HappyClient) Post(route string, body interface{}) (*http.Response, error) {
	return c.makeRequest(http.MethodPost, route, body)
}
func (c *HappyClient) PostParsed(route string, body, result interface{}) error {
	resp, err := c.Post(route, body)
	if err != nil {
		return errors.Wrap(err, "request failed with")
	}

	return c.parseResponse(resp, result)
}

func (c *HappyClient) Delete(route string, body interface{}) (*http.Response, error) {
	return c.makeRequest(http.MethodDelete, route, body)
}
func (c *HappyClient) DeleteParsed(route string, body, result interface{}) error {
	resp, err := c.Delete(route, body)
	if err != nil {
		return errors.Wrap(err, "request failed with")
	}

	return c.parseResponse(resp, result)
}

func (c *HappyClient) parseResponse(resp *http.Response, result interface{}) error {
	fmt.Println("...> resp.StatusCode", resp.StatusCode)
	err := InspectForErrors(resp)
	if err != nil {
		return errors.Wrap(err, "response error inspection failed")
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

	err := c.addAuth(req)
	if err != nil {
		return nil, err
	}

	fmt.Println("headers:", req.Header)

	return c.client.Do(req)
}

func (c *HappyClient) addAuth(req *http.Request) error {
	fmt.Println("...> route:", req.URL.Path)

	fmt.Println("...>about to create token")

	token, err := c.tokenProvider.GetToken()
	if err != nil {
		return errors.Wrap(err, "failed to get token")
	}
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", token))
	// req.Header.Add("Cookie", fmt.Sprintf("_oauth2_proxy=%s", token))
	// req.AddCookie(&http.Cookie{
	// 	Name:  "_oauth2_proxy",
	// 	Value: token,
	// })

	return nil
}
