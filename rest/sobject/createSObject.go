package sobject

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateSObjectRequest struct {
	Version        string
	SObjectApiName string
	Fields         interface{}
}

func (req CreateSObjectRequest) GetMethod() string {
	return http.MethodPost
}
func (req CreateSObjectRequest) GetHeaders() map[string]string {
	return nil
}
func (req CreateSObjectRequest) GetPath(
	version string,
) string {
	v := req.Version
	if len(v) == 0 {
		v = version
	}

	ret := fmt.Sprintf(
		"/services/data/v%s/sobjects/%s",
		version,
		req.SObjectApiName,
	)

	return ret
}
func (req CreateSObjectRequest) GetBody() ([]byte, error) {
	return json.Marshal(req.Fields)
}
