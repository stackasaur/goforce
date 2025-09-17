package sobject

import (
	"fmt"
	"net/http"
	"net/url"
)

type BlobGetRequest struct {
	Version        string
	SObjectApiName string
	RecordId       string
	BlobField      string
}

func (req BlobGetRequest) GetMethod() (string, error) {
	return http.MethodGet, nil
}
func (req BlobGetRequest) GetHeaders() (map[string]string, error) {
	return nil, nil
}
func (req BlobGetRequest) GetPath(
	version string,
) (*url.URL, error) {
	v := req.Version
	if len(v) == 0 {
		v = version
	}
	ret, err := url.Parse(fmt.Sprintf(
		"/services/data/v%s/sobjects/%s/%s/%s",
		v,
		req.SObjectApiName,
		req.RecordId,
		req.BlobField,
	))
	if err != nil {
		return nil, err
	}

	return ret, nil
}
func (req BlobGetRequest) GetBody() ([]byte, error) {
	return nil, nil
}
