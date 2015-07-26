package driver

import (
	"errors"
	"github.com/fluent/fluent-logger-golang/fluent"
	"strconv"
	"time"
)

type FluentDriver struct {
	fluent *fluent.Fluent
	Tag    string
}

func (f *FluentDriver) Configure(config map[string]string) (err error) {
	fluentConfig := fluent.Config{}

	host, ok := config["host"]
	if ok {
		fluentConfig.FluentHost = host
	}

	port, ok := config["port"]
	if ok {
		fluentConfig.FluentPort, err = strconv.Atoi(port)
		if err != nil {
			return
		}
	}

	timeout, ok := config["timeout"]
	if ok {
		fluentConfig.Timeout, err = time.ParseDuration(timeout)
		if err != nil {
			return
		}
	}

	tagPrefix, ok := config["prefix"]
	if ok {
		fluentConfig.TagPrefix = tagPrefix
	}

	f.Tag, ok = config["tag"]
	if !ok {
		return errors.New("tag is required")
	}

	f.fluent, err = fluent.New(fluentConfig)

	return
}

func (f *FluentDriver) Write(data []byte) (int, error) {
	err := f.fluent.Post(f.Tag, map[string]string{"output": string(data)})

	return len(data), err
}

func (f *FluentDriver) Close() error {
	return f.fluent.Close()
}
