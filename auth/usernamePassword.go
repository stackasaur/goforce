package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type UsernamePasswordFlow struct {
	ClientId      string
	ClientSecret  string
	Username      string
	Password      string
	SecurityToken string
	TokenEndpoint string
}

func (flow *UsernamePasswordFlow) NewToken(
	httpClient *http.Client,
) (Token, error) {
	payload := make(url.Values)
	payload["client_id"] = []string{flow.ClientId}
	payload["client_secret"] = []string{flow.ClientSecret}
	payload["username"] = []string{flow.Username}
	payload["password"] = []string{
		fmt.Sprintf(
			"%s%s",
			flow.Password,
			flow.SecurityToken,
		),
	}
	payload["grant_type"] = []string{"password"}

	res, err := httpClient.PostForm(
		flow.TokenEndpoint,
		payload,
	)
	if err != nil {
		return Token{}, errors.Join(
			errors.New("error creating request"),
			err,
		)
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
			errors.New("error decoding sfdc auth response"),
			decodeError,
		)
	}
	return Token{}, errors.Join(
		errors.New("sfdc auth error"),
		errorResponse,
	)
}
func (flow *UsernamePasswordFlow) RefreshToken(
	httpClient *http.Client,
) (Token, error) {
	return flow.NewToken(httpClient)
}
