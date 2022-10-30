package acry

import (
	"encoding/json"
)

const (
	algo_name_argon2i       = "argon2i"
	default_argon2i_time    = 1
	default_argon2i_memory  = 64 * 1024
	default_argon2i_threads = 1
	default_argon2i_keylen  = 32
)

type argon2iParameter struct {
	Time    int    `json:"time"`
	Keylen  int    `json:"klen"`
	Memory  int    `json:"mem"`
	Threads int    `json:"th"`
	Salt    string `json:"salt"`
}

func NewArgon2i() *argon2iParameter {
	return &argon2iParameter{
		Time:    default_argon2i_time,
		Keylen:  default_argon2i_keylen,
		Memory:  default_argon2i_memory,
		Threads: default_argon2i_threads,
	}
}

func (a *argon2iParameter) SetTime(time int) {
	a.Time = time
}

func (a *argon2iParameter) SetKeyLength(len int) {
	a.Keylen = len
}

func (a *argon2iParameter) SetMemory(mem int) {
	a.Memory = mem
}

func (a *argon2iParameter) SetThreads(t int) {
	a.Threads = t
}

func (a *argon2iParameter) SetSalt(salt string) {
	a.Salt = salt
}

func (a *argon2iParameter) Generate(client *connectionPool, password string) (string, error) {
	parameters, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	request := NewRequest(algo_name_argon2i, tcp_operation_generation, parameters)

	request.SetPassword(password)

	response, err := client.Query(request)
	if err != nil {
		return "", err
	}

	return response.GetHash()
}

func (a *argon2iParameter) Verify(client *connectionPool, password string, hash string) (bool, error) {
	parameters, err := json.Marshal(a)
	if err != nil {
		return false, err
	}
	request := NewRequest(algo_name_argon2i, tcp_operation_verify, parameters)
	request.SetHash(hash)
	request.SetPassword(password)

	response, err := client.Query(request)
	if err != nil {
		return false, err
	}

	return response.GetVerified()
}
