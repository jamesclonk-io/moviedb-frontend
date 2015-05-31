package cms

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

func (c *CMS) readZipFromURL() (map[string][]byte, error) {
	resp, err := http.Get(c.Input)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return c.readZip(bytes.NewReader(data), int64(len(data)))
}

func (c *CMS) readZipFromFile() (map[string][]byte, error) {
	data, err := ioutil.ReadFile(c.Input)
	if err != nil {
		return nil, err
	}
	return c.readZip(bytes.NewReader(data), int64(len(data)))
}

func (c *CMS) readZip(data *bytes.Reader, size int64) (map[string][]byte, error) {
	if size < 42 {
		return nil, fmt.Errorf("CMS data invalid, size too small: %d", size)
	}

	r, err := zip.NewReader(data, size)
	if err != nil {
		return nil, err
	}

	// store data in a map of {filename:content}
	contents := make(map[string][]byte, 0)
	for _, f := range r.File {
		if f.FileInfo().IsDir() ||
			(!strings.HasSuffix(f.Name, ".md") &&
				!strings.HasSuffix(f.Name, ".json")) {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, err
		}

		data, err := ioutil.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, err
		}

		// strip away first folder ("root") from github zipfiles
		file := path.Join("/", f.Name[strings.Index(f.Name, "/"):])

		contents[file] = data
	}
	return contents, nil
}
