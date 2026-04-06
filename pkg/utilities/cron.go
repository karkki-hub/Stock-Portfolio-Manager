package utilities

import (
	"log"

	"github.com/robfig/cron/v3"
)

type CronManager struct {
	cron *cron.Cron
}

func NewCronManager() *CronManager {
	return &CronManager{cron: cron.New()}
}

func (c *CronManager) Start() {
	c.cron.Start()
	log.Println("cron started")
}

func (c *CronManager) AddJob(spec string, job func()) {
	_, err := c.cron.AddFunc(spec, job)
	if err != nil {
		log.Printf("Error adding cron job: %v", err)
	}
}
