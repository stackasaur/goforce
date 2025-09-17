package sobject

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

type BlobUpdateRequest struct {
	Version        string
	SObjectApiName string
	RecordId       string
	BinaryPartName string
	BinaryData     []byte
	FieldsPartName string
	Fields         interface{}
	FileName       string
}

func (req BlobUpdateRequest) GetMethod() string {
	return http.MethodPatch
}
func (req BlobUpdateRequest) GetHeaders() map[string]string {
	return nil
}
func (req BlobUpdateRequest) GetPath(
	version string,
) string {
	v := req.Version
	if len(v) == 0 {
		v = version
	}
	ret := fmt.Sprintf(
		"/services/data/v%s/sobjects/%s/%s",
		v,
		req.SObjectApiName,
		req.RecordId,
	)

	return ret
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
