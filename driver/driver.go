package driver

import (
	"strings"
)

type Driver interface {
	Configure(map[string]string) error
	Write([]byte) (int, error)
	Close() error
}

func New(name string) Driver {
	driverName := strings.ToLower(name)

	if n := strings.IndexRune(driverName, '.'); n > 0 {
		driverName = driverName[0:n]
	}

	switch strings.ToLower(driverName) {
	case "file":
		return &FileDriver{}
	case "redis":
		return &RedisDriver{}
	case "fluent":
		return &FluentDriver{}
	case "pusher":
		return &PusherDriver{}
	}

	return nil
}
