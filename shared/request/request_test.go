package request

import (
	"net/http"
	"net/url"
	"testing"
)

func TestSfdcRequestAsHttpRequest(t *testing.T) {

	relativeUrl, _ := url.Parse(
		"/foobar",
	)
	sfdcRequest := GenericRequest{
		Headers: map[string]string{
			"testing": "testing",
		},
		Path:   relativeUrl,
		Method: http.MethodGet,
		Body:   nil,
	}

	baseUrl, _ := url.Parse(
		"https://example.com",
	)

	httpRequest, err := SfdcRequestAsHttpRequest(
		sfdcRequest,
		baseUrl,
		"60.0",
	)

	if err != nil {
		t.Fatal(err)
	}
	expectedUrl := "https://example.com/foobar"
	actualUrl := httpRequest.URL.String()

	if expectedUrl != actualUrl {
		t.Fatalf(
			"expected %v, actual %v",
			expectedUrl,
			actualUrl,
		)
	}
}
