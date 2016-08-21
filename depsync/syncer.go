package depsync

import (
	"github.com/micromdm/dep"
	"time"
	"github.com/go-kit/kit/log"
	"fmt"
)

type Syncer interface {
	Start(deviceChan chan dep.Device)
	Fetch(deviceChan chan dep.Device) error
}

type Cursor struct {
	Value string
	Created *time.Time // >7 days cursor is invalid.
}

type syncer struct {
	logger log.Logger
	client dep.Client
	tickerChan <-chan time.Time
	Cursor *Cursor
}

func NewSyncer(client dep.Client, logger log.Logger, tickerChan <-chan time.Time) Syncer {
	return &syncer{
		logger: logger,
		client: client,
		tickerChan: tickerChan,
	}
}

func (s *syncer) Fetch(deviceChan chan dep.Device) error {
	fetchDevices, err := s.client.FetchDevices(dep.Limit(100))
	if err != nil {
		return err
	}

	created := time.Now()
	s.Cursor = &Cursor{
		Value: fetchDevices.Cursor,
		Created: &created,
	}

	for _, dev := range fetchDevices.Devices {
		deviceChan <- dev
	}

	return nil
}

func (s *syncer) Start(deviceChan chan dep.Device) {
	if err := s.Fetch(deviceChan); err != nil {
		s.logger.Log("level", "error", "msg", "Fetching initial snapshot of devices from DEP")
	}

	for range s.tickerChan {
		s.logger.Log("level", "debug", "msg", "Synchronizing devices from DEP service")

		response, err := s.client.SyncDevices(s.Cursor.Value)
		if err != nil {
			s.logger.Log("level", "warn", "msg", fmt.Sprintf("Unable to fetch devices: %s", err))
		}

		for _, dev := range response.Devices {
			deviceChan <- dev
		}
	}
}