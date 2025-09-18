package auth

import (
	"net/http"
	"os"
	"testing"
)

func TestClientCredentials(t *testing.T) {
	t.Log("testing")
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	tokenEndpoint := os.Getenv("TOKEN_ENDPOINT")

	authFlow := ClientCredentialsFlow{
		ClientId:      clientId,
		ClientSecret:  clientSecret,
		TokenEndpoint: tokenEndpoint,
	}

	httpClient := http.Client{}

	token, err := authFlow.NewToken(&httpClient)
	if err != nil {
		t.Fatal(err)
	}

	userId := token.Id

	t.Logf(
		"userId: %s",
		userId,
	)

	if len(userId) == 0 {
		t.Fatal(
			"expected userId to be populated",
		)
	}
}
