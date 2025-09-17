package auth

import (
	"fmt"
	"net/http"
	"time"
)

// abstraction of the token.

type Token struct {
	Id          string
	AccessToken string
	InstanceUrl string
	Expiration  time.Time
}

// individual implementations of an AuthFlow are expected to handle managing
// token lifecycles and refresh tokens (if applicable). The AuthFlow interface
// is expected to simply call NewToken when no token exists and refresh token
// when a 401 error is returned or the Expiration of the token is reached.
type AuthFlow interface {
	NewToken(
		httpClient *http.Client,
	) (Token, error)
	RefreshToken(
		httpClient *http.Client,
	) (Token, error)
}

// default endpoints for ease
const ProductionTokenEndpoint string = "https://login.salesforce.com/services/oauth2/token"
const ProductionAuthorizationEndpoint string = "https://login.salesforce.com/services/oauth2/authorize"
const SandboxTokenEndpoint string = "https://test.salesforce.com/services/oauth2/token"
const SandboxAuthorizationEndpoint string = "https://test.salesforce.com/services/oauth2/authorize"

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	InstanceUrl  string `json:"instance_url"`
	IssuedAt     int64  `json:"issued_at"`
	Id           string `json:"id"`
}

type AuthError struct {
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (err AuthError) Error() string {
	return fmt.Sprintf(
		"%s: %s",
		err.ErrorCode,
		err.ErrorDescription,
	)
}
