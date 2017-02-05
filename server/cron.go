package main

import (
	"github.com/nlopes/slack"
	"github.com/robfig/cron"
)

func RegisterCronJob(api *slack.Client) {
	c := cron.New()
	// gonna have to figure out timezones
	c.AddFunc("0 0 21 * * MON-FRI", func() { runCronPost(api) })
	c.Start()
}

func runCronPost(api *slack.Client) {

}
