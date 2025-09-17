package sobject

import (
	"fmt"
	"net/http"
	"time"
)

type GetSObjectRequest struct {
	Version           string
	SObjectApiName    string
	RecordId          string
	Fields            string
	IfMatch           string
	IfNoneMatch       string
	IfModifiedSince   time.Time
	IfUnmodifiedSince time.Time
}

func (req GetSObjectRequest) GetMethod() string {
	return http.MethodGet
}
func (req GetSObjectRequest) GetHeaders() map[string]string {
	headers := map[string]string{}

	if len(req.IfMatch) > 0 {
		headers["If-Match"] = req.IfMatch
	}
	if len(req.IfNoneMatch) > 0 {
		headers["If-None-Match"] = req.IfNoneMatch
	}
	if !req.IfModifiedSince.IsZero() {
		headers["If-Modified-Since"] = req.IfModifiedSince.Format(
			time.RFC1123,
		)
	}
	if !req.IfUnmodifiedSince.IsZero() {
		headers["If-Unmodified-Since"] = req.IfUnmodifiedSince.Format(
			time.RFC1123,
		)
	}
	return headers
}
func (req GetSObjectRequest) GetPath(
	version string,
) string {
	v := req.Version
	if len(v) == 0 {
		v = version
	}
	ret := fmt.Sprintf(
		"/services/data/v%s/sobjects/%s/%s/",
		v,
		req.SObjectApiName,
		req.RecordId,
	)
	if len(req.Fields) > 0 {
		ret = fmt.Sprintf(
			"%s?fields=%s",
			ret,
			req.Fields,
		)
	}

	return ret
}
func (req GetSObjectRequest) GetBody() ([]byte, error) {
	return nil, nil
}
