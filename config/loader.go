package config

import (
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"strings"
)

var loader = koanf.New(".")

func Load(envPrefix string) error {
	err := loader.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil)
	fmt.Printf("config: %+v\r\n", loader.All())
	return err
}

func Unmarshal[T any](prefix string) (T, error) {
	var t T
	if err := loader.Unmarshal(prefix, t); err != nil {
		return t, err
	}
	return t, nil
}
