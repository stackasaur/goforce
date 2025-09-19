package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type SfdcRequest interface {
	GetMethod() (string, error)
	GetHeaders() (map[string]string, error)
	GetPath(
		string,
	) (*url.URL, error)
	GetBody() ([]byte, error)
}

type GenericRequest struct {
	Headers map[string]string
	Method  string
	Path    *url.URL
	Body    []byte
}

func (req GenericRequest) GetMethod() (string, error) {
	return req.Method, nil
}
func (req GenericRequest) GetHeaders() (map[string]string, error) {
	return req.Headers, nil
}
func (req GenericRequest) GetPath(
	_ string,
) (*url.URL, error) {
	return req.Path, nil
}
func (req GenericRequest) GetBody() ([]byte, error) {
	return req.Body, nil
}

type CompositeSubrequest struct {
	HttpHeaders map[string]string `json:"httpHeaders"`
	Method      string            `json:"method"`
	ReferenceId string            `json:"referenceId,omitempty"`
	Url         string            `json:"url"`
	Body        *json.RawMessage  `json:"body,omitempty"`
}

type SubRequestable interface {
	IntoSubRequest(string, string) (*CompositeSubrequest, error)
}

func SfdcRequestAsSubRequest(
	sfdcReq SfdcRequest,
	version string,
	referenceId string,
) (*CompositeSubrequest, error) {
	headers, err := sfdcReq.GetHeaders()
	if err != nil {
		return nil, err
	}
	method, err := sfdcReq.GetMethod()
	if err != nil {
		return nil, err
	}
	path, err := sfdcReq.GetPath(version)
	if err != nil {
		return nil, err
	}

	body, err := sfdcReq.GetBody()
	if err != nil {
		return nil, err
	}

	rawBody := json.RawMessage(body)
	ret := CompositeSubrequest{
		HttpHeaders: headers,
		Method:      method,
		ReferenceId: referenceId,
		Url:         path.String(),
		Body:        &rawBody,
	}

	return &ret, nil
}

// converts a SfdcRequest into an http request to be called by a client.
// this is used internally by the sfdc client when send it invoked.
func SfdcRequestAsHttpRequest(
	sfdcReq SfdcRequest,
	baseUrl *url.URL,
	version string,
) (*http.Request, error) {
	bodyBytes, err := sfdcReq.GetBody()
	if err != nil {
		return nil, errors.Join(
			errors.New(
				"error getting body",
			),
			err,
		)
	}
	method, err := sfdcReq.GetMethod()
	if err != nil {
		return nil, errors.Join(
			errors.New(
				"error getting method",
			),
			err,
		)
	}
	path, err := sfdcReq.GetPath(version)
	if err != nil {
		return nil, errors.Join(
			errors.New(
				"error getting path",
			),
			err,
		)
	}
	headers, err := sfdcReq.GetHeaders()
	if err != nil {
		return nil, errors.Join(
			errors.New(
				"error getting headers",
			),
			err,
		)
	}

	endpoint := baseUrl.ResolveReference(path)
	ret, err := http.NewRequest(
		method,
		endpoint.String(),
		bytes.NewReader(
			bodyBytes,
		),
	)
	if err != nil {
		return nil, errors.Join(
			errors.New(
				"error building request",
			),
			err,
		)
	}

	for key, value := range headers {
		ret.Header.Add(
			key,
			value,
		)
	}

	return ret, nil

}

// a custom error type to denote an actual error received from a salesforce api.
// a general http request error such as a timeout would not generate this error.
type ApiError struct {
	Message   string   `json:"message"`
	ErrorCode string   `json:"errorCode"`
	Fields    []string `json:"fields"`
}

func (r ApiError) Error() string {
	return fmt.Sprintf("%s: %s", r.ErrorCode, r.Message)
}
