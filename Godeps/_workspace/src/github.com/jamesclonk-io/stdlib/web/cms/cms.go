package cms

import (
	"html/template"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/env"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/jamesclonk-io/stdlib/web"
)

type CMS struct {
	Title    string
	frontend *web.Frontend
	Data     *CMSData
	Input    string
	mutex    *sync.Mutex
	log      *logrus.Logger
}

type CMSData struct {
	configuration *cmsConfiguration
	navigation    *cmsNavigation
	Content       []*CMSContent
	Timestamp     time.Time
}

type cmsConfiguration struct {
	Configuration map[string]interface{} `json:"configuration"`
}

type cmsNavigation struct {
	Navigation web.Navigation `json:"navigation"`
}

type CMSContent struct {
	Name     string
	Basename string
	Path     string
	Content  template.HTML
}

func NewCMS(frontend *web.Frontend) (*CMS, error) {
	data := &CMSData{}
	input := env.Get("JCIO_CMS_DATA", "https://github.com/jamesclonk-io/content/archive/master.zip")
	mutex := &sync.Mutex{}
	log := logger.GetLogger()

	cms := &CMS{frontend.Title, frontend, data, input, mutex, log}
	if err := cms.checkData(true); err != nil {
		return nil, err
	}
	return cms, nil
}

func (c *CMS) GetConfiguration() map[string]interface{} {
	return c.Data.configuration.Configuration
}

func (c *CMS) GetNavigation() web.Navigation {
	return c.Data.navigation.Navigation
}
