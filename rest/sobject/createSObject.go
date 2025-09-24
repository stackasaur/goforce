package sobject

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/stackasaur/goforce/client"
	Req "github.com/stackasaur/goforce/shared/request"
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
	return map[string]string{
		"Content-Type": "application/json",
	}, nil
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

func CreateSObject(
	sfdcClient *client.Client,
	request *CreateSObjectRequest,
) (string, error) {
	httpResponse, err := sfdcClient.Send(
		request,
	)
	if err != nil {
		return "", err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == 201 {
		var ret SObjectResponse
		decodeError := json.NewDecoder(httpResponse.Body).Decode(&ret)

		if decodeError != nil {
			return "", decodeError
		}

		if ret.Success {
			return ret.Id, nil
		} else {
			return "", ret.Errors[0]
		}
	}
	var errorResponse []Req.ApiError
	decodeError := json.NewDecoder(httpResponse.Body).Decode(&errorResponse)
	if decodeError != nil {
		return "", decodeError
	}
	if len(errorResponse) > 0 {
		return "", errorResponse[0]
	}
	return "", ErrUnknown
}
