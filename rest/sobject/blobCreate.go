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

func (req BlobCreateRequest) GetMethod() (string, error) {
	return http.MethodPost, nil
}
func (req BlobCreateRequest) GetHeaders() (map[string]string, error) {
	return nil, nil
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

	h := make(textproto.MIMEHeader)
	h.Set(
		"Content-Disposition",
		fmt.Sprintf(
			`form-data; name="%s"; filename=""`,
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
