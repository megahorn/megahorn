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
	switch strings.ToLower(name) {
	case "file":
		return &FileDriver{}
	case "pusher":
		return &PusherDriver{}
	}

	return nil
}
