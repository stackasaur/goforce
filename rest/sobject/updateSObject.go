package sobject

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/stackasaur/goforce/client"
	Req "github.com/stackasaur/goforce/shared/request"
)

type UpdateSObjectRequest struct {
	Version           string
	SObjectApiName    string
	RecordId          string
	Fields            interface{}
	IfMatch           string
	IfNoneMatch       string
	IfModifiedSince   time.Time
	IfUnmodifiedSince time.Time
}

func (req UpdateSObjectRequest) GetMethod() (string, error) {
	return http.MethodPatch, nil
}
func (req UpdateSObjectRequest) GetHeaders() (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}

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
func (req UpdateSObjectRequest) GetPath(
	version string,
) (*url.URL, error) {
	v := req.Version
	if len(v) == 0 {
		v = version
	}

	ret, err := url.Parse(fmt.Sprintf(
		"/services/data/v%s/sobjects/%s/%s/",
		version,
		req.SObjectApiName,
		req.RecordId,
	))
	if err != nil {
		return nil, err
	}

	return ret, nil
}
func (req UpdateSObjectRequest) GetBody() ([]byte, error) {
	return json.Marshal(req.Fields)
}

func UpdateSObject(
	sfdcClient *client.Client,
	request *UpdateSObjectRequest,
) error {
	httpResponse, err := sfdcClient.Send(
		request,
	)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode >= 200 && httpResponse.StatusCode < 300 {
		return nil
	}
	var errorResponse []Req.ApiError
	decodeError := json.NewDecoder(httpResponse.Body).Decode(&errorResponse)
	if decodeError != nil {
		return decodeError
	}
	if len(errorResponse) > 0 {
		return errorResponse[0]
	}
	return ErrUnknown
}
