package main

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Route struct {
	Listen string
	SNI []string
	Dial string
}

type Config struct {
	Routes []Route
}

func (c *Config) Load(filename string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(b, c); err != nil {
		return err
	}

	return nil
}
