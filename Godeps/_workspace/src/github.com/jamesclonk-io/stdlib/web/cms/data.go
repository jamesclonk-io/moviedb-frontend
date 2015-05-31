package cms

import (
	"encoding/json"
	"html/template"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/russross/blackfriday"
)

func (c *CMS) checkData(refresh bool) error {
	// reset data if not set
	if c.Data == nil {
		c.Data = &CMSData{}
	}

	// refresh either every 24 hours, or if refresh parameter set to true
	if time.Since(c.Data.Timestamp).Hours() >= 24 || refresh {
		// refresh cms content data
		if err := c.refreshData(); err != nil {
			return err
		}

		// update navigation
		c.frontend.SetNavigation(c.GetNavigation())
	}
	return nil
}

func (c *CMS) refreshData() (err error) {
	if err := c.getData(); err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err,
			"file":  c.Input,
		}).Error("Could not refresh data")
		return err
	}
	return nil
}

func (c *CMS) getData() error {
	var data map[string][]byte
	var err error

	if strings.HasSuffix(c.Input, ".zip") {
		if strings.HasPrefix(c.Input, "http") {
			// read zip content from url
			data, err = c.readZipFromURL()
		} else {
			// read zip content from local file
			data, err = c.readZipFromFile()
		}
	} else {
		// read zip content from local folder
		data, err = c.readFromFolder()
	}
	if err != nil {
		return err
	}

	c.Data = &CMSData{
		Content:   make([]*CMSContent, 0),
		Timestamp: time.Now(),
	}

	// go through all files
	for file, bytes := range data {
		basename := filepath.Base(file)

		if strings.HasSuffix(basename, ".md") {
			html := blackfriday.MarkdownCommon(bytes) // generate HTML from markdown
			content := &CMSContent{
				Name:     path.Base(file),
				Basename: strings.TrimSuffix(basename, filepath.Ext(basename)),
				Path:     path.Dir(file),
				Content:  template.HTML(html),
			}
			c.Data.Content = append(c.Data.Content, content)

		} else if basename == "navigation.json" {
			var nav cmsNavigation
			if err := json.Unmarshal(bytes, &nav); err != nil {
				return err
			}
			c.Data.navigation = &nav

		} else if basename == "configuration.json" {
			var config cmsConfiguration
			if err := json.Unmarshal(bytes, &config); err != nil {
				return err
			}
			c.Data.configuration = &config
		}
	}
	return nil
}
