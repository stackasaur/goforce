package sobject

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/stackasaur/goforce/client"
	Req "github.com/stackasaur/goforce/shared/request"
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

type Blob struct {
	Data          []byte
	ContentType   string
	ContentLength int64
}

func BlobGet(
	sfdcClient *client.Client,
	request *BlobGetRequest,
) (*Blob, error) {
	httpResponse, err := sfdcClient.Send(
		request,
	)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode == 200 {
		blob := Blob{
			ContentType: httpResponse.Header.Get("Content-Type"),
		}
		blob.Data, err = io.ReadAll(httpResponse.Body)
		if err != nil {
			return nil, err
		}
		blob.ContentLength = int64(len(blob.Data))

		return &blob, nil
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
