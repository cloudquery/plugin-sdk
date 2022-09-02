package docs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// PrepareDirectory is a helper function that creates a directory if it does not exist, or
// removes all files in a directory if it does exist. This leaves a clean slate for documentation
// to be generated into.
func PrepareDirectory(dirname string) error {
	files, err := ioutil.ReadDir(dirname)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		// create directory if it does not exist
		return os.MkdirAll(dirname, 0744)
	}
	for _, d := range files {
		if err := os.RemoveAll(path.Join(dirname, d.Name())); err != nil {
			return fmt.Errorf("failed to remove files: %s\n", err)
		}
	}
	return nil
}
