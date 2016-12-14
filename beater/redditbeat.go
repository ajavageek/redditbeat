package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/nfrankel/redditbeat/config"
	"net/http"
	"io/ioutil"
)

type Redditbeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	configuration := config.DefaultConfig
	if err := cfg.Unpack(&configuration); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Redditbeat{
		done: make(chan struct{}),
		config: configuration,
	}
	return bt, nil
}

func (bt *Redditbeat) Run(b *beat.Beat) error {
	logp.Info("redditbeat is running! Hit CTRL-C to stop it.")
	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	reddit := "https://www.reddit.com/r/" + bt.config.Subreddit + "/.json"
	client := &http.Client {}
	logp.Info("URL configured to " + reddit)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}
		req, reqErr := http.NewRequest("GET", reddit, nil)
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
		if (reqErr != nil) {
			panic(reqErr)
		}
		resp, getErr := client.Do(req)
		if (getErr != nil) {
			panic(getErr)
		}
		status := resp.Status
		logp.Info("HTTP status code is " + status)
		body, readErr := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if (readErr != nil) {
			panic(readErr)
		}
		message := string(body)
		logp.Info(message)
		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       b.Name,
			"message":    message,
		}
		bt.client.PublishEvent(event)
		logp.Info("Event sent")
	}
}

func (bt *Redditbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
