package sobject

import (
	"errors"

	Req "github.com/stackasaur/goforce/shared/request"
)

type SObjectResponse struct {
	Id      string         `json:"id"`
	Errors  []Req.ApiError `json:"errors"`
	Success bool           `json:"success"`
}

var ErrUnknown = errors.New("unknown sobject error")
