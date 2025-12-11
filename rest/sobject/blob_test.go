package sobject

import (
	"os"
	"testing"

	"github.com/stackasaur/goforce/auth"
	"github.com/stackasaur/goforce/client"
)

func TestBlobMethods(t *testing.T) {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	tokenEndpoint := os.Getenv("TOKEN_ENDPOINT")
	accountId := os.Getenv("ACCOUNT_ID")

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

	var recordId string
	t.Run(
		"Create Blob",
		func(t *testing.T) {
			blobCreateRequest := BlobCreateRequest{
				SObjectApiName: "Attachment",
				BinaryPartName: "Body",
				BinaryData:     []byte("testing"),
				FieldsPartName: "entity_attachment",
				Fields: map[string]any{
					"ParentId": accountId,
					"Name":     "test.txt",
				},
				FileName: "test.txt",
			}

			body, err := blobCreateRequest.GetBody()

			if err != nil {
				t.Fatal(err)
			}
			t.Log(string(body))

			recordId, err = BlobCreate(
				sfdcClient,
				&blobCreateRequest,
			)
			if err != nil {
				t.Fatal(err)
			}
		},
	)

	t.Run(
		"Update Blob",
		func(t *testing.T) {
			blobUpdateRequest := BlobUpdateRequest{
				SObjectApiName: "Attachment",
				BinaryPartName: "Body",
				BinaryData:     []byte("testing2"),
				FieldsPartName: "entity_attachment",
				RecordId:       recordId,
				Fields:         map[string]any{},
				FileName:       "test.txt",
			}

			body, err := blobUpdateRequest.GetBody()

			if err != nil {
				t.Fatal(err)
			}
			t.Log(string(body))

			err = BlobUpdate(
				sfdcClient,
				&blobUpdateRequest,
			)
			if err != nil {
				t.Fatal(err)
			}
		},
	)

	t.Run(
		"Get Blob",
		func(t *testing.T) {
			blobGetRequest := BlobGetRequest{
				SObjectApiName: "Attachment",
				RecordId:       recordId,
				BlobField:      "Body",
			}

			blob, err := BlobGet(
				sfdcClient,
				&blobGetRequest,
			)
			if err != nil {
				t.Fatal(err)
			}

			body := string(blob.Data)

			if body != "testing2" {
				t.Fatalf(
					"expected file data: 'testing2', received: '%s'",
					body,
				)
			}
		},
	)

	if len(recordId) > 0 {
		deleteSObjectRequest := DeleteSObjectRequest{
			SObjectApiName: "Attachment",
			RecordId:       recordId,
		}
		err = DeleteSObject(
			sfdcClient,
			&deleteSObjectRequest,
		)
		if err != nil {
			t.Fatal(err)
		}
	}
}
