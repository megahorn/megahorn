package main

import (
	"encoding/json"
	"github.com/webminal/webminal/driver"
	"io"
	"os"
)

type Config struct {
	Stdout map[string]map[string]string `json:"stdout"`
	Stderr map[string]map[string]string `json:"stderr"`
}

func newConfig() *Config {
	return &Config{
		Stdout: map[string]map[string]string{},
		Stderr: map[string]map[string]string{},
	}
}

func (c *Config) LoadFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)

	return decoder.Decode(c)
}

func (c *Config) StdoutDrviers(echo bool) []io.WriteCloser {
	drivers := []io.WriteCloser{}

	if echo {
		drivers = append(drivers, os.Stdout)
	}

	for name, config := range c.Stdout {
		stdoutDriver := driver.New(name)
		if stdoutDriver != nil {
			err := stdoutDriver.Configure(config)
			if err != nil {
				os.Stderr.WriteString(err.Error())
				os.Stderr.WriteString("\n")
			} else {
				drivers = append(drivers, stdoutDriver)
			}
		}
	}

	return drivers
}

func (c *Config) StderrDrviers(echo bool) []io.WriteCloser {
	drivers := []io.WriteCloser{}

	if echo {
		drivers = append(drivers, os.Stderr)
	}

	for name, config := range c.Stderr {
		stderrDriver := driver.New(name)
		if stderrDriver != nil {
			err := stderrDriver.Configure(config)
			if err != nil {
				os.Stderr.WriteString(err.Error())
				os.Stderr.WriteString("\n")
			} else {
				drivers = append(drivers, stderrDriver)
			}
		}
	}

	return drivers
}
