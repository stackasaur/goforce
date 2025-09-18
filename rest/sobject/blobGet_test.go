package sobject

import (
	"os"
	"testing"

	"github.com/stackasaur/goforce/auth"
	"github.com/stackasaur/goforce/client"
)

func TestBlobGet(t *testing.T) {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	tokenEndpoint := os.Getenv("TOKEN_ENDPOINT")
	recordId := os.Getenv("ACCOUNT_ID")

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

	blobCreateRequest := BlobCreateRequest{
		SObjectApiName: "Attachment",
		BinaryPartName: "Body",
		BinaryData:     []byte("testing"),
		FieldsPartName: "entity_attachment",
		Fields: map[string]any{
			"ParentId": recordId,
			"Name":     "test.txt",
		},
		FileName: "test.txt",
	}

	recordId, err = BlobCreate(
		*sfdcClient,
		blobCreateRequest,
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(recordId)

	blobGetRequest := BlobGetRequest{
		SObjectApiName: "Attachment",
		RecordId:       recordId,
		BlobField:      "Body",
	}

	blob, err := BlobGet(
		*sfdcClient,
		blobGetRequest,
	)
	if err != nil {
		t.Fatal(err)
	}

	body := string(blob.Data)

	if body != "testing" {
		t.Fatalf(
			"expected file data: 'testing', received: '%s'",
			body,
		)
	}

	deleteSObjectRequest := DeleteSObjectRequest{
		SObjectApiName: "Attachment",
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
