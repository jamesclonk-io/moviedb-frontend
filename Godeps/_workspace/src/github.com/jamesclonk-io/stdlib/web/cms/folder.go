package cms

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (c *CMS) readFromFolder() (map[string][]byte, error) {
	data := make(map[string][]byte)

	if err := filepath.Walk(c.Input, c.fileWalker(data)); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *CMS) fileWalker(data map[string][]byte) filepath.WalkFunc {
	return func(file string, info os.FileInfo, err error) error {
		if info.IsDir() ||
			(!strings.HasSuffix(info.Name(), ".md") &&
				!strings.HasSuffix(info.Name(), ".json")) {
			return nil
		}

		content, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		// strip away everything from path until the "root" folder
		file = path.Join("/", strings.TrimPrefix(file, c.Input))

		data[file] = content

		return nil
	}
}
