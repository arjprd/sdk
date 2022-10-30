package acry

import "fmt"

const (
	default_required_connection = 1
	default_maximum_connection  = 10
	default_retry_interval      = 100
)

type acryConfig struct {
	required      uint
	maximum       uint
	retryInterval uint
	host          string
}

func NewAcryConfig() *acryConfig {
	return &acryConfig{
		required:      default_required_connection,
		maximum:       default_maximum_connection,
		retryInterval: default_retry_interval,
	}
}

func (c *acryConfig) SetRequiredConnectionPool(poolSize uint) {
	c.required = poolSize
}

func (c *acryConfig) GetRequiredConnectionPool() uint {
	if c.required > 0 {
		return c.required
	}
	return default_required_connection
}

func (c *acryConfig) SetMaximumConnectionPool(poolSize uint) {
	c.maximum = poolSize
}

func (c *acryConfig) GetMaximumConnectionPool() uint {
	if c.maximum > c.GetRequiredConnectionPool() {
		return c.maximum
	}
	return c.GetRequiredConnectionPool()
}

func (c *acryConfig) SetRetryInterval(interval uint) {
	c.required = interval
}

func (c *acryConfig) GetRetryInterval() uint {
	if c.retryInterval > 0 {
		return c.retryInterval
	}
	return default_retry_interval
}

func (c *acryConfig) SetHost(address string, port uint) {
	c.host = fmt.Sprintf("%s:%d", address, port)
}

func (c *acryConfig) GetHost() string {
	return c.host
}
