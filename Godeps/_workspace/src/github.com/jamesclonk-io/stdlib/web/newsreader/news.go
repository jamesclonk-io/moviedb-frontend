package newsreader

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jamesclonk-io/stdlib/logger"
	"github.com/jamesclonk-io/stdlib/web"
)

type NewsReader struct {
	Title     string
	Config    *NewsConfig
	Feeds     []Feed
	mutex     *sync.Mutex
	log       *logrus.Logger
	timestamp time.Time
}

type NewsConfig struct {
	CacheDuration time.Duration
	Feeds         []string
}

type configJson struct {
	Newsreader struct {
		CacheDuration int64    `json:"cache_duration"` // in seconds
		Feeds         []string `json:"feeds"`
	} `json:"newsreader"`
}

func NewReader(frontend *web.Frontend, configuration map[string]interface{}) (*NewsReader, error) {
	data, err := json.Marshal(configuration)
	if err != nil {
		return nil, err
	}

	var c configJson
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	cacheDuration, err := time.ParseDuration(fmt.Sprintf("%ds", c.Newsreader.CacheDuration))
	if err != nil {
		return nil, err
	}

	config := &NewsConfig{cacheDuration, c.Newsreader.Feeds}
	mutex := &sync.Mutex{}
	log := logger.GetLogger()

	return &NewsReader{frontend.Title, config, nil, mutex, log, time.Unix(0, 0)}, nil
}
