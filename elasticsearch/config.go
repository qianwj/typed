package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8"
	"net/http"
)

type ConfigBuilder struct {
	config elasticsearch.Config
}

func NewBuilder() *ConfigBuilder {
	return &ConfigBuilder{}
}

func (c *ConfigBuilder) Address(addrs []string) *ConfigBuilder {
	c.config.Addresses = addrs
	return c
}

func (c *ConfigBuilder) Username(username string) *ConfigBuilder {
	c.config.Username = username
	return c
}

func (c *ConfigBuilder) Password(password string) *ConfigBuilder {
	c.config.Password = password
	return c
}

func (c *ConfigBuilder) CloudID(cloudId string) *ConfigBuilder {
	c.config.CloudID = cloudId
	return c
}

func (c *ConfigBuilder) APIKey(apiKey string) *ConfigBuilder {
	c.config.APIKey = apiKey
	return c
}

func (c *ConfigBuilder) ServiceToken(token string) *ConfigBuilder {
	c.config.ServiceToken = token
	return c
}

func (c *ConfigBuilder) DisableRetry() *ConfigBuilder {
	c.config.DisableRetry = true
	return c
}

func (c *ConfigBuilder) Retry(maxRetries int) *ConfigBuilder {
	c.config.MaxRetries = maxRetries
	return c
}

func (c *ConfigBuilder) RetryOnStatus(status []int) *ConfigBuilder {
	c.config.RetryOnStatus = status
	return c
}

func (c *ConfigBuilder) RetryOnError(retryOnError func(*http.Request, error) bool) *ConfigBuilder {
	c.config.RetryOnError = retryOnError
	return c
}

func (c *ConfigBuilder) build() elasticsearch.Config {
	return c.config
}
