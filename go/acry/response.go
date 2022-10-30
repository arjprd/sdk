package acry

import "errors"

type response struct {
	Hash       string `json:"hash"`
	Verfied    bool   `json:"verified"`
	IsError    bool   `json:"error"`
	ErrMessage string `json:"message"`
	conIndex   uint
}

func (r *response) GetHash() (string, error) {
	if r.IsError {
		return "", errors.New(r.ErrMessage)
	}

	return r.Hash, nil
}

func (r *response) GetVerified() (bool, error) {
	if r.IsError {
		return false, errors.New(r.ErrMessage)
	}

	return r.Verfied, nil
}
