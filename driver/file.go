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

	f.path = path
	f.File, err = os.Create(path)

	return
}
