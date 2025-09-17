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

type BlobUpdateRequest struct {
	Version        string
	SObjectApiName string
	RecordId       string
	BinaryPartName string
	BinaryData     []byte
	FieldsPartName string
	Fields         any
	FileName       string
}

func (req BlobUpdateRequest) GetMethod() (string, error) {
	return http.MethodPatch, nil
}
func (req BlobUpdateRequest) GetHeaders() (map[string]string, error) {
	return nil, nil
}
func (req BlobUpdateRequest) GetPath(
	version string,
) (*url.URL, error) {
	v := req.Version
	if len(v) == 0 {
		v = version
	}
	ret, err := url.Parse(fmt.Sprintf(
		"/services/data/v%s/sobjects/%s/%s",
		v,
		req.SObjectApiName,
		req.RecordId,
	))
	if err != nil {
		return nil, err
	}

	return ret, nil
}
func (req BlobUpdateRequest) GetBody() ([]byte, error) {
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
