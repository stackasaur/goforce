package sobject

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stackasaur/goforce/auth"
	"github.com/stackasaur/goforce/client"
)

type Account struct {
	Id   string `json:",omitempty"`
	Name string
}

func TestSObjectFunctions(t *testing.T) {
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

	name1 := fmt.Sprintf(
		"test_%x",
		time.Now().UnixMilli(),
	)
	name2 := name1 + "_2"

	createSObjectRequest := CreateSObjectRequest{
		SObjectApiName: "Account",
		Fields: Account{
			Name: name1,
		},
	}

	recordId, err := CreateSObject(
		*sfdcClient,
		createSObjectRequest,
	)
	if err != nil {
		t.Fatal(err)
	}

	getSObjectRequest := GetSObjectRequest{
		SObjectApiName: "Account",
		RecordId:       recordId,
		Fields:         "Id,Name",
	}

	acct, err := GetSObject[Account](
		*sfdcClient,
		getSObjectRequest,
	)
	if err != nil {
		t.Fatal(err)
	}

	if acct.Name != name1 {
		t.Fatalf(
			"expected name: %s, received: %s",
			name1,
			acct.Name,
		)
	}

	updateSObjectRequest := UpdateSObjectRequest{
		SObjectApiName: "Account",
		RecordId:       recordId,
		Fields: Account{
			Name: name2,
		},
	}

	err = UpdateSObject(
		*sfdcClient,
		updateSObjectRequest,
	)
	if err != nil {
		t.Fatal(err)
	}

	acct, err = GetSObject[Account](
		*sfdcClient,
		getSObjectRequest,
	)
	if err != nil {
		t.Fatal(err)
	}

	if acct.Name != name2 {
		t.Fatalf(
			"expected name: %s, received: %s",
			name2,
			acct.Name,
		)
	}

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

	_, err = GetSObject[Account](
		*sfdcClient,
		getSObjectRequest,
	)
	if err == nil {
		t.Fatal(
			"expected record to be deleted",
		)
	}
}
