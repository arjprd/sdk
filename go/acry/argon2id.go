package acry

import (
	"encoding/json"
)

const (
	algo_name_argon2id       = "argon2id"
	default_argon2id_time    = 1
	default_argon2id_memory  = 64 * 1024
	default_argon2id_threads = 1
	default_argon2id_keylen  = 32
)

type argon2iDParameter struct {
	Time    int    `json:"time"`
	Keylen  int    `json:"klen"`
	Memory  int    `json:"mem"`
	Threads int    `json:"th"`
	Salt    string `json:"salt"`
}

func NewArgon2iD() *argon2iDParameter {
	return &argon2iDParameter{
		Time:    default_argon2id_time,
		Keylen:  default_argon2id_keylen,
		Memory:  default_argon2id_memory,
		Threads: default_argon2id_threads,
	}
}

func (a *argon2iDParameter) SetTime(time int) {
	a.Time = time
}

func (a *argon2iDParameter) SetKeyLength(len int) {
	a.Keylen = len
}

func (a *argon2iDParameter) SetMemory(mem int) {
	a.Memory = mem
}

func (a *argon2iDParameter) SetThreads(t int) {
	a.Threads = t
}

func (a *argon2iDParameter) SetSalt(salt string) {
	a.Salt = salt
}

func (a *argon2iDParameter) Generate(client *connectionPool, password string) (string, error) {
	parameters, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	request := NewRequest(algo_name_argon2id, tcp_operation_generation, parameters)

	request.SetPassword(password)

	response, err := client.Query(request)
	if err != nil {
		return "", err
	}

	return response.GetHash()
}

func (a *argon2iDParameter) Verify(client *connectionPool, password string, hash string) (bool, error) {
	parameters, err := json.Marshal(a)
	if err != nil {
		return false, err
	}
	request := NewRequest(algo_name_argon2id, tcp_operation_verify, parameters)
	request.SetHash(hash)
	request.SetPassword(password)

	response, err := client.Query(request)
	if err != nil {
		return false, err
	}

	return response.GetVerified()
}
