package sobject

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type CreateSObjectRequest struct {
	Version        string
	SObjectApiName string
	Fields         any
}

func (req CreateSObjectRequest) GetMethod() (string, error) {
	return http.MethodPost, nil
}
func (req CreateSObjectRequest) GetHeaders() (map[string]string, error) {
	return nil, nil
}
func (req CreateSObjectRequest) GetPath(
	version string,
) (*url.URL, error) {
	v := req.Version
	if len(v) == 0 {
		v = version
	}

	ret, err := url.Parse(fmt.Sprintf(
		"/services/data/v%s/sobjects/%s",
		version,
		req.SObjectApiName,
	))
	if err != nil {
		return nil, err
	}

	return ret, nil
}
func (req CreateSObjectRequest) GetBody() ([]byte, error) {
	return json.Marshal(req.Fields)
}
