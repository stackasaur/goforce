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

type QueryOptions struct {
	BatchSize int
	QueryAll  bool
}
type QueryRequest struct {
	Version      string
	Query        string
	QueryOptions QueryOptions
}

func (req QueryRequest) GetMethod() (string, error) {
	return http.MethodGet, nil
}
func (req QueryRequest) GetHeaders() (map[string]string, error) {
	ret := map[string]string{}
	if req.QueryOptions.BatchSize != 0 {
		ret["Sforce-Query-Options"] = fmt.Sprintf(
			"batchSize=%d",
			req.QueryOptions.BatchSize,
		)
	}
	return ret, nil
}
func (req QueryRequest) GetPath(
	version string,
) (*url.URL, error) {
	v := req.Version
	if len(v) == 0 {
		v = version
	}
	var queryPath string
	if req.QueryOptions.QueryAll {
		queryPath = "queryAll"
	} else {
		queryPath = "query"
	}

	ret, err := url.Parse(fmt.Sprintf(
		"/services/data/v%s/%s",
		v,
		queryPath,
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

	queryResponse, err := handleQueryRequest[T](
		sfdcClient,
		request,
	)
	if err != nil {
		return nil, err
	}
	records := queryResponse.Records

	// get all results
	for res := queryResponse; !res.Done && len(res.NextRecordsUrl) > 0; {
		endpoint, err := url.Parse(res.NextRecordsUrl)
		if err != nil {
			break
		}
		req := Req.GenericRequest{
			Method: http.MethodGet,
			Path:   endpoint,
		}

		res, err = handleQueryRequest[T](
			sfdcClient,
			req,
		)

		if err != nil {
			break
		}

		records = append(records, res.Records...)
	}
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
			return nil, err
		}
		return queryResponse.Records, nil
	}

	var errorResponse []Req.ApiError
	decodeError := json.NewDecoder(httpResponse.Body).Decode(&errorResponse)
	if decodeError != nil {
		return nil, err
	}
	if len(errorResponse) > 0 {
		return nil, errorResponse[0]
	}
	return nil, ErrUnknown

}

func handleQueryRequest[T any](
	sfdcClient *client.Client,
	request Req.SfdcRequest,
) (*QueryResponse[T], error) {
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
			return nil, err
		}
		return &queryResponse, nil
	}

	var errorResponse []Req.ApiError
	decodeError := json.NewDecoder(httpResponse.Body).Decode(&errorResponse)
	if decodeError != nil {
		return nil, err
	}
	if len(errorResponse) > 0 {
		return nil, errorResponse[0]
	}
	return nil, ErrUnknown
}

var ErrUnknown = errors.New("unknown query error")
