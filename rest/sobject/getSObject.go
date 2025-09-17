package sobject

import (
	"fmt"
	"net/http"
	"net/url"
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

func (req GetSObjectRequest) GetMethod() (string, error) {
	return http.MethodGet, nil
}
func (req GetSObjectRequest) GetHeaders() (map[string]string, error) {
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
	return headers, nil
}
func (req GetSObjectRequest) GetPath(
	version string,
) (*url.URL, error) {
	v := req.Version
	if len(v) == 0 {
		v = version
	}
	ret, err := url.Parse(fmt.Sprintf(
		"/services/data/v%s/sobjects/%s/%s/",
		v,
		req.SObjectApiName,
		req.RecordId,
	))
	if err != nil {
		return nil, err
	}
	if len(req.Fields) > 0 {
		q := ret.Query()
		q.Add(
			"fields",
			req.Fields,
		)
		ret.RawQuery = q.Encode()
	}

	return ret, nil
}
func (req GetSObjectRequest) GetBody() ([]byte, error) {
	return nil, nil
}
