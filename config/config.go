// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period time.Duration `config:"period"`
	Subreddit string `config:"subreddit"`
}

var DefaultConfig = Config{
	Period: 15 * time.Second,
	Subreddit: "elastic",
}
