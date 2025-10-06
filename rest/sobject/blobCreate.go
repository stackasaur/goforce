package sobject

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"time"

	"github.com/stackasaur/goforce/client"
	Req "github.com/stackasaur/goforce/shared/request"
)

type BlobCreateRequest struct {
	Version        string
	SObjectApiName string
	BinaryPartName string
	BinaryData     []byte
	FieldsPartName string
	Fields         any
	FileName       string
}

var boundary string = fmt.Sprintf(
	"goforce%x",
	time.Now().UnixMilli(),
)

func (req BlobCreateRequest) GetMethod() (string, error) {
	return http.MethodPost, nil
}
func (req BlobCreateRequest) GetHeaders() (map[string]string, error) {
	return map[string]string{
		"Content-Type": fmt.Sprintf(
			"multipart/form-data; boundary=%s",
			boundary,
		),
	}, nil
}
func (req BlobCreateRequest) GetPath(
	version string,
) (*url.URL, error) {
	v := req.Version
	if len(v) == 0 {
		v = version
	}
	ret, err := url.Parse(fmt.Sprintf(
		"/services/data/v%s/sobjects/%s",
		v,
		req.SObjectApiName,
	))
	if err != nil {
		return nil, err
	}

	return ret, nil
}
func (req BlobCreateRequest) GetBody() ([]byte, error) {
	fieldData, jsonErr := json.Marshal(req.Fields)

	if jsonErr != nil {
		return nil, errors.Join(
			jsonErr,
			errors.New("error marshalling fields"),
		)
	}
	requestBody := bytes.Buffer{}

	multipartWriter := multipart.NewWriter(&requestBody)
	multipartWriter.SetBoundary(boundary)

	h := make(textproto.MIMEHeader)
	h.Set(
		"Content-Disposition",
		fmt.Sprintf(
			`form-data; name="%s"`,
			req.FieldsPartName,
		),
	)
	h.Set(
		"Content-Type",
		"application/json",
	)
	fieldWriter, writerErr := multipartWriter.CreatePart(h)

	if writerErr != nil {
		return nil, errors.Join(
			writerErr,
			errors.New("error creating form file writer"),
		)
	}
	_, writerErr = fieldWriter.Write(fieldData)
	if writerErr != nil {
		return nil, errors.Join(
			writerErr,
			errors.New("error writing data to file"),
		)
	}

	contentWriter, writerErr := multipartWriter.CreateFormFile(
		req.BinaryPartName,
		req.FileName,
	)
	if writerErr != nil {
		return nil, errors.Join(
			writerErr,
			errors.New("error creating form file writer"),
		)
	}
	_, writerErr = contentWriter.Write(req.BinaryData)
	if writerErr != nil {
		return nil, errors.Join(
			writerErr,
			errors.New("error writing data to file"),
		)
	}

	multipartWriter.Close()
	return requestBody.Bytes(), nil
}

func BlobCreate(
	sfdcClient *client.Client,
	request *BlobCreateRequest,
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
