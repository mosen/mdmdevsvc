package depsync

import (
	"github.com/micromdm/dep"
	"time"
	"github.com/go-kit/kit/log"
)

type Syncer interface {
	Start()
}

type syncer struct {
	logger log.Logger
	client dep.Client
	tickerChan <-chan time.Time
}

func NewSyncer(client dep.Client, logger log.Logger, tickerChan <-chan time.Time) Syncer {
	return &syncer{
		logger,
		client,
		tickerChan,
	}
}

func (s *syncer) Start() {
	for range s.tickerChan {
		s.logger.Log("level", "debug", "msg", "synchronizing devices from DEP service")
	}
}