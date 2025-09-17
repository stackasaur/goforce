package sobject

import (
	"fmt"
	"net/http"
)

type BlobGetRequest struct {
	Version        string
	SObjectApiName string
	RecordId       string
	BlobField      string
}

func (req BlobGetRequest) GetMethod() string {
	return http.MethodGet
}
func (req BlobGetRequest) GetHeaders() map[string]string {
	return nil
}
func (req BlobGetRequest) GetPath(
	version string,
) string {
	v := req.Version
	if len(v) == 0 {
		v = version
	}
	ret := fmt.Sprintf(
		"/services/data/v%s/sobjects/%s/%s/%s",
		v,
		req.SObjectApiName,
		req.RecordId,
		req.BlobField,
	)

	return ret
}
func (req BlobGetRequest) GetBody() ([]byte, error) {
	return nil, nil
}
