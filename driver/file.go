package driver

import (
	"errors"
	"os"
)

type FileDriver struct {
	*os.File
	path string
}

func (f *FileDriver) Configure(config map[string]string) (err error) {
	path, ok := config["path"]
	if !ok {
		return errors.New("path is required")
	}

	mode, ok := config["mode"]
	if !ok {
		mode = "append"
	}

	f.path = path

	if mode == "append" {
		f.File, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0666)
	} else {
		f.File, err = os.Create(path)
	}

	return
}
