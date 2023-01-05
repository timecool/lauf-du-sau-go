package utils

import (
	"os"
)

func CreateFolder(dirname string) error {
	_, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		errDir := os.Mkdir(dirname, 0777)
		if errDir != nil {
			return errDir
		}
	}
	return nil
}
