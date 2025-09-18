package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/stackasaur/goforce/auth"
	Req "github.com/stackasaur/goforce/shared/request"
)

type Client struct {
	context    context.Context
	httpClient *http.Client
	authFlow   auth.AuthFlow
	token      auth.Token
	version    string
}

func (client *Client) GetHttpClient() *http.Client {
	return client.httpClient
}
func (client *Client) GetContext() context.Context {
	return client.context
}
func (client *Client) GetVersion() string {
	return client.version
}
func (client *Client) SetVersion(
	version int,
) error {
	if version == 0 {
		return errors.New("version must be specified")
	}
	if !validateVersion(version) {
		return errors.New("invalid version")
	}
	client.version = toVersionString(version)

	return nil
}

func (client *Client) GetUserId() string {
	userId := client.token.Id
	if len(userId) == 0 {
		return ""
	}
	splt := strings.Split(client.token.Id, "/")

	return splt[len(splt)-1]

}

func (client *Client) Send(
	req Req.SfdcRequest,
) (*http.Response, error) {
	httpClient := client.GetHttpClient()

	token := client.token

	if !token.Expiration.After(time.Now()) {
		var err error
		token, err = client.authFlow.RefreshToken(
			httpClient,
		)
		client.token = token
		if err != nil {
			return nil, err
		}
	}
	baseUrl, err := url.Parse(token.InstanceUrl)
	if err != nil {
		return nil, err
	}

	httpRequest, err := Req.SfdcRequestAsHttpRequest(
		req,
		baseUrl,
		client.version,
	)
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set(
		"Authorization",
		fmt.Sprintf("Bearer %v", token.AccessToken),
	)

	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode == 401 {
		// refresh token and try again

		var err error
		token, err = client.authFlow.RefreshToken(
			httpClient,
		)
		client.token = token
		if err != nil {
			return nil, err
		}
		httpRequest.Header.Set(
			"Authorization",
			fmt.Sprintf("Bearer %v", token.AccessToken),
		)

		return httpClient.Do(httpRequest)
	}

	return httpResponse, nil
}

type ClientConfig struct {
	HttpClient *http.Client
	Context    context.Context
	AuthFlow   auth.AuthFlow
	Version    int
}

func NewClient(
	config ClientConfig,
) (*Client, error) {
	ctx := config.Context
	if ctx == nil {
		ctx = context.Background()
	}
	httpClient := config.HttpClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	version := config.Version
	if version == 0 {
		version = DefaultVersion
	}
	if !validateVersion(version) {
		return nil, errors.New("invalid version")
	}

	if config.AuthFlow == nil {
		return nil, errors.New("authflow is required")
	}

	authFlow := config.AuthFlow
	token, err := authFlow.NewToken(
		httpClient,
	)
	if err != nil {
		return nil, err
	}

	client := Client{
		context:    ctx,
		httpClient: httpClient,
		version:    toVersionString(version),
		authFlow:   config.AuthFlow,
		token:      token,
	}

	return &client, nil
}
