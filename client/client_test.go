package client

import (
	"os"
	"testing"

	"github.com/stackasaur/goforce/auth"
)

func TestClient(t *testing.T) {
	t.Log("testing")
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	tokenEndpoint := os.Getenv("TOKEN_ENDPOINT")

	authFlow := auth.ClientCredentialsFlow{
		ClientId:      clientId,
		ClientSecret:  clientSecret,
		TokenEndpoint: tokenEndpoint,
	}

	sfdcClient, err := NewClient(
		ClientConfig{
			Version:  "60.0",
			AuthFlow: authFlow,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	userId := sfdcClient.GetUserId()

	if len(userId) == 0 {
		t.Fatal(
			"expected userId to be populated",
		)
	}
}
