package client

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hetiansu5/urlquery"
	"github.com/pkg/errors"
)

type TokenProvider interface {
	GetToken() (string, error)
}

type AWSCredentialsProvider interface {
	GetCredentials(context.Context) (aws.Credentials, error)
}

type HappyClient struct {
	client          http.Client
	apiBaseUrl      string
	clientName      string
	clientVersion   string
	tokenProvider   TokenProvider
	awsCredProvider AWSCredentialsProvider
}

func NewHappyClient(clientName, clientVersion, apiBaseUrl string, tokenProvider TokenProvider, awsCredProvider AWSCredentialsProvider) *HappyClient {
	return &HappyClient{
		apiBaseUrl:      apiBaseUrl,
		clientName:      clientName,
		clientVersion:   clientVersion,
		client:          *http.DefaultClient,
		tokenProvider:   tokenProvider,
		awsCredProvider: awsCredProvider,
	}
}

func (c *HappyClient) Get(route string, payload interface{}) (*http.Response, error) {
	return c.makeRequestWithQueryString(http.MethodGet, route, payload)
}

func (c *HappyClient) GetParsed(route string, body, result interface{}) error {
	resp, err := c.Get(route, body)
	if err != nil {
		return errors.Wrap(err, "request failed with")
	}

	return c.parseResponse(resp, result)
}

func (c *HappyClient) Post(route string, body interface{}) (*http.Response, error) {
	return c.makeRequestWithBody(http.MethodPost, route, body)
}

func (c *HappyClient) PostParsed(route string, body, result interface{}) error {
	resp, err := c.Post(route, body)
	if err != nil {
		return errors.Wrap(err, "request failed with")
	}

	return c.parseResponse(resp, result)
}

func (c *HappyClient) Delete(route string, body interface{}) (*http.Response, error) {
	return c.makeRequestWithBody(http.MethodDelete, route, body)
}

func (c *HappyClient) DeleteParsed(route string, body, result interface{}) error {
	resp, err := c.Delete(route, body)
	if err != nil {
		return errors.Wrap(err, "request failed with")
	}

	return c.parseResponse(resp, result)
}

func (c *HappyClient) parseResponse(resp *http.Response, result interface{}) error {
	err := InspectForErrors(resp)
	if err != nil {
		return errors.Wrap(err, "response error inspection failed")
	}

	return ParseResponse(resp, &result)
}

func (c *HappyClient) makeRequestWithQueryString(method, route string, payload interface{}) (*http.Response, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.apiBaseUrl, route), nil)
	if err != nil {
		return nil, err
	}

	if payload != nil {
		queryBytes, err := urlquery.Marshal(payload)
		if err != nil {
			return nil, errors.Wrap(err, "could not convert body to json")
		}
		req.URL.RawQuery = string(queryBytes)
	}

	return c.Do(req)
}

func (c *HappyClient) makeRequestWithBody(method, route string, body interface{}) (*http.Response, error) {
	var bodyReader *bytes.Reader
	if body != nil {
		bodyJson, err := json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(err, "could not convert body to json")
		}
		bodyReader = bytes.NewReader(bodyJson)
	}

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

	return c.client.Do(req)
}

func (c *HappyClient) addAuth(req *http.Request) error {
	token, err := c.tokenProvider.GetToken()
	if err != nil {
		return errors.Wrap(err, "failed to get oidc token")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	creds, err := c.awsCredProvider.GetCredentials(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to get aws credentials")
	}
	req.Header.Add("x-aws-access-key-id", b64.StdEncoding.EncodeToString([]byte(creds.AccessKeyID)))
	req.Header.Add("x-aws-secret-access-key", b64.StdEncoding.EncodeToString([]byte(creds.SecretAccessKey)))
	req.Header.Add("x-aws-session-token", creds.SessionToken) // SessionToken is already base64 encoded

	return nil
}
