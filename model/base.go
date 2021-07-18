package model

import "errors"

type BaseResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Err     string `json:"error,omitempty"`
}

func (b *BaseResponse) SetError(err error) {
	if err != nil {
		b.Err = err.Error()
	}
}

func (b *BaseResponse) Error() error {
	if b.Err != "" {
		return errors.New(b.Err)
	}
	return nil
}
