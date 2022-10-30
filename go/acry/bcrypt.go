package acry

import (
	"encoding/json"
	"log"
)

const (
	default_bcrypt_cost = 10
	algo_name_bcrypt    = "bcrypt"
)

type bcryptParameter struct {
	Cost int `json:"cost"`
}

func NewBcrypt() *bcryptParameter {
	return &bcryptParameter{
		Cost: default_bcrypt_cost,
	}
}

func (b *bcryptParameter) SetCost(cost int) {
	b.Cost = cost
}

func (b *bcryptParameter) Generate(client *connectionPool, password string) (string, error) {
	parameters, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	request := NewRequest(algo_name_bcrypt, tcp_operation_generation, parameters)

	request.SetPassword(password)

	response, err := client.Query(request)
	if err != nil {
		return "", err
	}
	log.Println(client.inuse, client.poolsize)
	return response.GetHash()
}

func (b *bcryptParameter) Verify(client *connectionPool, password string, hash string) (bool, error) {
	parameters, err := json.Marshal(b)
	if err != nil {
		return false, err
	}
	request := NewRequest(algo_name_bcrypt, tcp_operation_verify, parameters)
	request.SetHash(hash)
	request.SetPassword(password)

	response, err := client.Query(request)
	if err != nil {
		return false, err
	}

	return response.GetVerified()
}
