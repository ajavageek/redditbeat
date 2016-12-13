package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/nfrankel/redditbeat/config"
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
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       b.Name,
			"counter":    counter,
		}
		bt.client.PublishEvent(event)
		logp.Info("Event sent")
		counter++
	}
}

func (bt *Redditbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
