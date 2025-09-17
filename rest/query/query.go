package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/stackasaur/goforce/client"
	Req "github.com/stackasaur/goforce/shared/request"
)

type QueryRequest struct {
	Version string
	Query   string
}

func (req QueryRequest) GetMethod() (string, error) {
	return http.MethodGet, nil
}
func (req QueryRequest) GetHeaders() (map[string]string, error) {
	return nil, nil
}
func (req QueryRequest) GetPath(
	version string,
) (*url.URL, error) {
	v := req.Version
	if len(v) == 0 {
		v = version
	}
	ret, err := url.Parse(fmt.Sprintf(
		"/services/data/v%s/query",
		v,
	))
	if err != nil {
		return nil, err
	}
	q := ret.Query()
	q.Add(
		"q",
		req.Query,
	)
	ret.RawQuery = q.Encode()

	return ret, nil
}
func (req QueryRequest) GetBody() ([]byte, error) {
	return nil, nil
}

type QueryResponse[T any] struct {
	TotalSize      int    `json:"totalSize"`
	Done           bool   `json:"done"`
	NextRecordsUrl string `json:"nextRecordsUrl"`
	Records        []T
}

func Query[T any](
	sfdcClient *client.Client,
	request *QueryRequest,
) ([]T, error) {
	httpResponse, err := sfdcClient.Send(
		request,
	)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == 200 {
		var queryResponse QueryResponse[T]
		decodeError := json.NewDecoder(httpResponse.Body).Decode(&queryResponse)

		if decodeError != nil {
			return nil, errors.Join(
				decodeError,
				errors.New("error execting query"),
			)
		}
		return queryResponse.Records, nil
	}

	var errorResponse []Req.ApiError
	decodeError := json.NewDecoder(httpResponse.Body).Decode(&errorResponse)
	if decodeError != nil {
		return nil, errors.Join(
			decodeError,
			errors.New("error decoding response"),
		)
	}
	if len(errorResponse) > 0 {
		return nil, errors.Join(
			errorResponse[0],
			errors.New("sfdc query error"),
		)
	}
	return nil, errors.New("unexpected query error")

}
