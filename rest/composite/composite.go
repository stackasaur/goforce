package composite

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/stackasaur/goforce/client"
	Req "github.com/stackasaur/goforce/shared/request"
)

type CompositeRequest struct {
	Version            string
	AllOrNone          bool
	CollateSubrequests bool
	SubRequests        []CompositeSubrequest
}
type compositeRequestBody struct {
	AllOrNone          bool                  `json:"allOrNone"`
	CollateSubrequests bool                  `json:"collateSubrequests"`
	CompositeRequest   []CompositeSubrequest `json:"compositeRequest"`
}

type CompositeSubrequest struct {
	HttpHeaders map[string]string `json:"httpHeaders"`
	Method      string            `json:"method"`
	ReferenceId string            `json:"referenceId,omitempty"`
	Url         string            `json:"url"`
	Body        *json.RawMessage  `json:"body,omitempty"`
}

type CompositeSubrequestResult struct {
	Body        *json.RawMessage  `json:"body"`
	HttpHeaders map[string]string `json:"httpHeaders"`
	StatusCode  int               `json:"statusCode"`
	ReferenceId string            `json:"referenceId"`
}
type CompositeResult struct {
	CompositeResponse []CompositeSubrequestResult `json:"compositeResponse"`
}

type SubRequestOptions struct {
	Version     string
	ReferenceId string
}

func SubRequest(
	sfdcReq Req.SfdcRequest,
	options *SubRequestOptions,
) (*CompositeSubrequest, error) {
	var version string
	if options != nil && len(options.Version) > 0 {
		version = options.Version
	}
	var referenceId string
	if options != nil && len(options.ReferenceId) > 0 {
		referenceId = options.ReferenceId
	}

	headers, err := sfdcReq.GetHeaders()
	if err != nil {
		return nil, err
	}
	delete(headers, "Content-Type")

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

func (req CompositeRequest) GetMethod() (string, error) {
	return http.MethodPost, nil
}
func (req CompositeRequest) GetHeaders() (map[string]string, error) {
	return map[string]string{
		"Content-Type": "application/json",
	}, nil
}
func (req CompositeRequest) GetPath(
	version string,
) (*url.URL, error) {
	v := req.Version
	if len(v) == 0 {
		v = version
	}

	ret, err := url.Parse(fmt.Sprintf(
		"/services/data/v%s/composite",
		version,
	))
	if err != nil {
		return nil, err
	}

	return ret, nil
}
func (req CompositeRequest) GetBody() ([]byte, error) {

	body := compositeRequestBody{
		AllOrNone:          req.AllOrNone,
		CollateSubrequests: req.CollateSubrequests,
		CompositeRequest:   req.SubRequests,
	}
	return json.Marshal(body)
}

func Composite(
	sfdcClient *client.Client,
	request *CompositeRequest,
) (*CompositeResult, error) {
	httpResponse, err := sfdcClient.Send(
		request,
	)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode >= 200 && httpResponse.StatusCode < 300 {
		var ret CompositeResult
		decodeError := json.NewDecoder(httpResponse.Body).Decode(&ret)

		if decodeError != nil {
			return nil, decodeError
		}

		return &ret, nil
	}
	var errorResponse []Req.ApiError
	decodeError := json.NewDecoder(httpResponse.Body).Decode(&errorResponse)
	if decodeError != nil {
		return nil, decodeError
	}
	if len(errorResponse) > 0 {
		return nil, errorResponse[0]
	}
	return nil, ErrUnknown
}

var ErrUnknown = errors.New("unknown composite error")
