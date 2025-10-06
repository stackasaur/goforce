package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stackasaur/goforce/auth"
	"github.com/stackasaur/goforce/client"
	composite "github.com/stackasaur/goforce/rest/composite"
	sobject "github.com/stackasaur/goforce/rest/sobject"
)

type Account struct {
	Id   string `json:",omitempty"`
	Name string
}

func TestCompositeFunctions(t *testing.T) {
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

	uniqueName := fmt.Sprintf(
		"test_%x",
		time.Now().UnixMilli(),
	)
	t.Run(
		"Create and Delete",
		func(t *testing.T) {

			createSObjectRequest := sobject.CreateSObjectRequest{
				Version:        "60.0",
				SObjectApiName: "Account",
				Fields: Account{
					Name: uniqueName,
				},
			}
			createSubrequest, err := composite.SubRequest(
				createSObjectRequest,
				&composite.SubRequestOptions{
					ReferenceId: "refAccount",
				},
			)
			if err != nil {
				t.Fatal(err)
			}

			deleteSObjectRequest := sobject.DeleteSObjectRequest{
				Version:        "60.0",
				SObjectApiName: "Account",
				RecordId:       "@{refAccount.id}",
			}
			deleteSubrequest, err := composite.SubRequest(
				deleteSObjectRequest,
				&composite.SubRequestOptions{
					ReferenceId: "refAccount2",
				},
			)
			if err != nil {
				t.Fatal(err)
			}

			compositeRequest := composite.CompositeRequest{
				Version:   "60.0",
				AllOrNone: true,
				SubRequests: []composite.CompositeSubrequest{
					*createSubrequest,
					*deleteSubrequest,
				},
			}

			result, err := composite.Composite(
				sfdcClient,
				&compositeRequest,
			)
			if err != nil {
				t.Fatal(err)
			}

			for _, it := range result.CompositeResponse {
				t.Log(it.StatusCode)
				b, err := json.Marshal(it.Body)
				if err != nil {
					t.Fatal(err)
				}

				// Convert the byte slice to a string
				jsonString := string(b)

				t.Log(jsonString)
			}

			t.Log(result)
		},
	)
}
