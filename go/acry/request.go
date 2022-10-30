package acry

import "encoding/json"

const (
	tcp_operation_generation = "generate"
	tcp_operation_verify     = "verify"
)

type request struct {
	Algorithm  string          `json:"algo"`
	Operation  string          `json:"op"`
	Password   string          `json:"pass"`
	Hash       string          `json:"hash"`
	Parameters json.RawMessage `json:"param"`
}

func NewRequest(algorithm string, operation string, parameters json.RawMessage) *request {
	return &request{
		Algorithm:  algorithm,
		Operation:  operation,
		Parameters: parameters,
	}
}

func (r *request) SetPassword(password string) {
	r.Password = password
}

func (r *request) SetHash(hash string) {
	r.Hash = hash
}
