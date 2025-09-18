package sobject

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stackasaur/goforce/auth"
	"github.com/stackasaur/goforce/client"
)

func TestCreateSObject(t *testing.T) {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	tokenEndpoint := os.Getenv("TOKEN_ENDPOINT")

	authFlow := auth.ClientCredentialsFlow{
		ClientId:      clientId,
		ClientSecret:  clientSecret,
		TokenEndpoint: tokenEndpoint,
	}

	sfdcClient, err := client.NewClient(
		client.ClientConfig{
			Version:  60,
			AuthFlow: authFlow,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	createSObjectRequest := CreateSObjectRequest{
		SObjectApiName: "Account",
		Fields: Account{
			Name: fmt.Sprintf(
				"test_%d",
				time.Now().UnixMilli(),
			),
		},
	}

	recordId, err := CreateSObject(
		*sfdcClient,
		createSObjectRequest,
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(recordId)

	deleteSObjectRequest := DeleteSObjectRequest{
		SObjectApiName: "Account",
		RecordId:       recordId,
	}
	err = DeleteSObject(
		*sfdcClient,
		deleteSObjectRequest,
	)
	if err != nil {
		t.Fatal(err)
	}
}
