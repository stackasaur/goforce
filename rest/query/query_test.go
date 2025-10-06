package query

import (
	"errors"
	"os"
	"testing"

	"github.com/stackasaur/goforce/auth"
	"github.com/stackasaur/goforce/client"
	Req "github.com/stackasaur/goforce/shared/request"
)

type Account struct {
	Id   string
	Name string
}
type Contact struct {
	Id   string
	Name string
}

func TestQueryRequest(t *testing.T) {
	queryRequest := QueryRequest{
		Version: "60.0",
		Query:   "SELECT Id, Name FROM Account",
	}

	actualUrl, err := queryRequest.GetPath("60.0")
	if err != nil {
		t.Fatal(err)
	}

	expectedQueryString := "q=SELECT+Id%2C+Name+FROM+Account"

	actualQueryString := actualUrl.RawQuery

	if expectedQueryString != actualQueryString {
		t.Fatalf(
			"expected %v, actual %v",
			expectedQueryString,
			actualQueryString,
		)

	}
}

func TestQuery(t *testing.T) {
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

	queryRequest := QueryRequest{
		Version: "60.0",
		Query:   "SELECT Id, Name FROM Account LIMIT 1",
	}

	accountResp, err := Query[Account](
		sfdcClient,
		&queryRequest,
	)
	if err != nil {
		t.Fatal(err)
	}

	accounts := accountResp.Records

	if len(accounts) != 1 {
		t.Fatalf(
			"expected 1 account, received %d",
			len(accounts),
		)
	}

	t.Log(accounts)

	badQueryRequest := QueryRequest{
		Version: "60.0",
		Query:   "SELECT Id, Name, FROM Account LIMIT 1",
	}

	_, err = Query[Account](
		sfdcClient,
		&badQueryRequest,
	)
	if err == nil {
		t.Fatal(
			"expected err",
		)
	}

	var apiError Req.ApiError
	if errors.As(err, &apiError) {
		if apiError.ErrorCode != "MALFORMED_QUERY" {
			t.Fatalf(
				"expected MALFORMED_QUERY error, received: %v",
				apiError,
			)
		}
	} else {
		t.Fatalf(
			"expected err to be of type ApiError: %v",
			err,
		)
	}

}

func TestQueryMore(t *testing.T) {
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

	queryRequest := QueryRequest{
		Version: "60.0",
		Query:   "SELECT Id, Name FROM Contact",
		QueryOptions: QueryOptions{
			BatchSize: 200,
		},
	}

	contactResp, err := Query[Contact](
		sfdcClient,
		&queryRequest,
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf(
		"totalSize: %d, records: %d",
		contactResp.TotalSize,
		len(contactResp.Records),
	)
	if contactResp.Done {
		t.Fatal("shouldn't be done, check the linked org.")
	}

	contactNextResp, err := contactResp.QueryMore(
		sfdcClient,
		QueryOptions{},
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf(
		"totalSize: %d, records: %d",
		contactNextResp.TotalSize,
		len(contactNextResp.Records),
	)

}
