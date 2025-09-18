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
		return ClientError{
			ErrorCode:        "VERSION_ERROR",
			ErrorDescription: "version must be specified",
		}
	}
	if !validateVersion(version) {
		return ClientError{
			ErrorCode:        "VERSION_ERROR",
			ErrorDescription: "invalid version",
		}
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
			return nil, errors.Join(
				ClientError{
					ErrorCode:        "TOKEN_ERROR",
					ErrorDescription: "error refreshing token",
				},
				err,
			)
		}
	}
	baseUrl, err := url.Parse(token.InstanceUrl)
	if err != nil {
		return nil, errors.Join(
			ClientError{
				ErrorCode:        "URL_ERROR",
				ErrorDescription: "error parsing url",
			},
			err,
		)
	}

	httpRequest, err := Req.SfdcRequestAsHttpRequest(
		req,
		baseUrl,
		client.version,
	)
	if err != nil {
		return nil, errors.Join(
			ClientError{
				ErrorCode:        "REQUEST_ERROR",
				ErrorDescription: "error compiling request",
			},
			err,
		)
	}

	httpRequest.Header.Set(
		"Authorization",
		fmt.Sprintf("Bearer %v", token.AccessToken),
	)

	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, errors.Join(
			ClientError{
				ErrorCode:        "HTTP_ERROR",
				ErrorDescription: "error performing request",
			},
			err,
		)
	}
	if httpResponse.StatusCode == 401 {
		// refresh token and try again

		var err error
		token, err = client.authFlow.RefreshToken(
			httpClient,
		)

		if err != nil {
			return nil, errors.Join(
				ClientError{
					ErrorCode:        "TOKEN_ERROR",
					ErrorDescription: "error refreshing token",
				},
				err,
			)
		}
		client.token = token
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
		return nil, ClientError{
			ErrorCode:        "VERSION_ERROR",
			ErrorDescription: "invalid version",
		}
	}

	if config.AuthFlow == nil {
		return nil, ClientError{
			ErrorCode:        "AUTHFLOW_ERROR",
			ErrorDescription: "authflow is required",
		}
	}

	authFlow := config.AuthFlow
	token, err := authFlow.NewToken(
		httpClient,
	)
	if err != nil {
		return nil, errors.Join(
			ClientError{
				ErrorCode:        "TOKEN_ERROR",
				ErrorDescription: "error getting token",
			},
			err,
		)
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

type ClientError struct {
	ErrorCode        string
	ErrorDescription string
}

func (err ClientError) Error() string {
	return fmt.Sprintf(
		"%s: %s",
		err.ErrorCode,
		err.ErrorDescription,
	)
}
