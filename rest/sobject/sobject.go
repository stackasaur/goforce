package sobject

import "errors"

// type SObjectResponse struct {
// 	Id      string     `json:"id"`
// 	Errors  []ApiError `json:"errors"`
// 	Success bool       `json:"success"`
// }

var ErrUnknown = errors.New("unknown sobject error")
