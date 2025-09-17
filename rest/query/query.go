package query

import (
	"fmt"
	"net/http"
	"net/url"
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
