package main

import (
	"encoding/json"
	"github.com/megahorn/megahorn/driver"
	"io"
	"os"
)

type Config struct {
	WorkingDir string                       `json"working_dir"`
	Env        map[string]string            `json:"env"`
	Stdout     map[string]map[string]string `json:"stdout"`
	Stderr     map[string]map[string]string `json:"stderr"`
}

func newConfig() *Config {
	wd, err := os.Getwd()
	if err != nil {
		wd = ""
	}

	return &Config{
		WorkingDir: wd,
		Env:        map[string]string{},
		Stdout:     map[string]map[string]string{},
		Stderr:     map[string]map[string]string{},
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

func (c *Config) Environments() []string {
	current := os.Environ()

	for key, value := range c.Env {
		current = append(current, key+"="+value)
	}

	return current
}
