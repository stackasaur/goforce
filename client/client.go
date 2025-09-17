package client

import (
	"context"
	"errors"
	"net/http"
)

type Client struct {
	context    context.Context
	httpClient *http.Client
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
	version string,
) error {
	if len(version) == 0 {
		return errors.New("version must be specified")
	}
	if !ValidateVersion(version) {
		return errors.New("invalid version")
	}
	client.version = version

	return nil
}

type ClientConfig struct {
	HttpClient *http.Client
	Context    context.Context
	Version    string
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
	if len(version) == 0 {
		version = DefaultVersion
	}
	if !ValidateVersion(version) {
		return nil, errors.New("invalid version")
	}

	client := Client{
		context:    ctx,
		httpClient: httpClient,
		version:    version,
	}

	return &client, nil
}
