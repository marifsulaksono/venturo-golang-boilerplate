package helpers

import (
	"simple-crud-rnd/config"

	"github.com/robfig/cron/v3"
)

type CronJob struct {
	Cron *cron.Cron
	Cfg  *config.Config
}

func NewCronJobInstance(cfg *config.Config) *CronJob {
	return &CronJob{
		Cron: cron.New(),
		Cfg:  cfg,
	}
}

func (c *CronJob) Start() {
	c.Cron.Start()
}

func (c *CronJob) Stop() {
	c.Cron.Stop()
}
