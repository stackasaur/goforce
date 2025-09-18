package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

type ClientCredentialsFlow struct {
	ClientId      string
	ClientSecret  string
	TokenEndpoint string
}

func (flow ClientCredentialsFlow) NewToken(
	httpClient *http.Client,
) (Token, error) {
	payload := make(url.Values)
	payload["client_id"] = []string{flow.ClientId}
	payload["client_secret"] = []string{flow.ClientSecret}
	payload["grant_type"] = []string{"client_credentials"}

	res, err := httpClient.PostForm(
		flow.TokenEndpoint,
		payload,
	)
	if err != nil {
		return Token{}, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		var tokenResponse TokenResponse
		json.NewDecoder(res.Body).Decode(&tokenResponse)

		expiration := time.UnixMilli(tokenResponse.IssuedAt).Add(
			time.Hour,
		)
		return Token{
			AccessToken: tokenResponse.AccessToken,
			InstanceUrl: tokenResponse.InstanceUrl,
			Expiration:  expiration,
			Id:          tokenResponse.Id,
		}, nil
	}
	var errorResponse AuthError
	decodeError := json.NewDecoder(res.Body).Decode(&errorResponse)
	if decodeError != nil {
		return Token{}, errors.Join(
			AuthError{
				ErrorCode:        "DECODING_ERROR",
				ErrorDescription: "error decoding sfdc auth response",
			},
			decodeError,
		)
	}
	return Token{}, errorResponse
}
func (flow ClientCredentialsFlow) RefreshToken(
	httpClient *http.Client,
) (Token, error) {
	return flow.NewToken(httpClient)
}
