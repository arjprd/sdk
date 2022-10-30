package acry

import (
	"encoding/json"
)

const (
	algo_name_pbkdf2         = "pbkdf2"
	default_pbkdf2_iter      = 1
	default_pbkdf2_keylen    = 32
	default_pbkdf2_hash_func = "sha1"

	PBKDF2_SHA1   = "sha1"
	PBKDF2_MD5    = "md5"
	PBKDF2_SHA256 = "sha256"
	PBKDF2_SHA512 = "sha512"
)

type pbkdf2Parameter struct {
	Iter     int    `json:"iter"`
	Keylen   int    `json:"klen"`
	HashFunc string `json:"hf"`
	Salt     string `json:"salt"`
}

func NewPbkdf2() *pbkdf2Parameter {
	return &pbkdf2Parameter{
		Iter:     default_pbkdf2_iter,
		Keylen:   default_pbkdf2_keylen,
		HashFunc: default_pbkdf2_hash_func,
	}
}

func (a *pbkdf2Parameter) SetIter(iteration int) {
	a.Iter = iteration
}

func (a *pbkdf2Parameter) SetKeyLength(len int) {
	a.Keylen = len
}

func (a *pbkdf2Parameter) SetHashFunction(hashFunc string) {
	a.HashFunc = hashFunc
}

func (a *pbkdf2Parameter) SetSalt(salt string) {
	a.Salt = salt
}

func (a *pbkdf2Parameter) Generate(client *connectionPool, password string) (string, error) {
	parameters, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	request := NewRequest(algo_name_pbkdf2, tcp_operation_generation, parameters)

	request.SetPassword(password)

	response, err := client.Query(request)
	if err != nil {
		return "", err
	}

	return response.GetHash()
}

func (a *pbkdf2Parameter) Verify(client *connectionPool, password string, hash string) (bool, error) {
	parameters, err := json.Marshal(a)
	if err != nil {
		return false, err
	}
	request := NewRequest(algo_name_pbkdf2, tcp_operation_verify, parameters)
	request.SetHash(hash)
	request.SetPassword(password)

	response, err := client.Query(request)
	if err != nil {
		return false, err
	}

	return response.GetVerified()
}
