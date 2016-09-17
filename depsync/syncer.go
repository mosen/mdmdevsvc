package depsync

import (
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/micromdm/dep"
	"time"
)

type Syncer interface {
	Start(deviceChan chan dep.Device)
	Pause() error
	Resume() error
	Fetch(deviceChan chan dep.Device) (bool, error)
}

type DEPDeviceRepository interface {
	StoreDEPDevices([]dep.Device) error
}

type Cursor struct {
	Value   string
	Created *time.Time // >7 days cursor is invalid.
}

type syncer struct {
	logger               log.Logger
	client               dep.Client
	tickerChan           <-chan time.Time
	tickerChanDisabled   <-chan time.Time
	Cursor               *Cursor
	InitialFetchComplete bool
}

func NewSyncer(client dep.Client, logger log.Logger, tickerChan <-chan time.Time) Syncer {
	return &syncer{
		logger:     logger,
		client:     client,
		tickerChan: tickerChan,
	}
}

// Fetch fetches a list of DEP devices in batches. The first return value indicates that there are more devices to fetch.
func (s *syncer) Fetch(deviceChan chan dep.Device) (bool, error) {
	var deviceResponse *dep.DeviceResponse
	var err error

	if s.Cursor != nil {
		deviceResponse, err = s.client.FetchDevices(dep.Limit(100), dep.Cursor(s.Cursor.Value))
	} else {
		deviceResponse, err = s.client.FetchDevices(dep.Limit(100))
	}

	if err != nil {
		return true, err
	}

	if s.Cursor == nil {
		created := time.Now()
		s.Cursor = &Cursor{
			Value:   deviceResponse.Cursor,
			Created: &created,
		}
	}

	s.logger.Log("level", "debug", "msg", fmt.Sprintf("Fetching %d devices", len(deviceResponse.Devices)))
	for _, dev := range deviceResponse.Devices {
		deviceChan <- dev
	}

	return deviceResponse.MoreToFollow, nil
}

// Sync fetches devices since the given cursor value
func (s *syncer) Sync(deviceChan chan dep.Device) (bool, error) {
	var deviceResponse *dep.DeviceResponse
	var err error

	if s.Cursor != nil {
		deviceResponse, err = s.client.SyncDevices(s.Cursor.Value, dep.Limit(100))
	} else {
		return true, errors.New("cannot sync dep devices without a cursor value")
	}

	if err != nil {
		return true, err
	}

	if s.Cursor == nil {
		created := time.Now()
		s.Cursor = &Cursor{
			Value:   deviceResponse.Cursor,
			Created: &created,
		}
	}

	s.logger.Log("level", "debug", "msg", fmt.Sprintf("Syncing %d devices", len(deviceResponse.Devices)))
	for _, dev := range deviceResponse.Devices {
		deviceChan <- dev
	}

	return deviceResponse.MoreToFollow, nil
}

// Start starts the synchronisation schedule which runs at a configured interval.
// The first run will always download a complete list of devices. Deltas will be fetched if the process is still running
// after the first sync interval
func (s *syncer) Start(deviceChan chan dep.Device) {
	for range s.tickerChan {
		if !s.InitialFetchComplete {
			more, err := s.Fetch(deviceChan)
			if err != nil {
				s.logger.Log("level", "error", "msg", fmt.Sprintf("Fetching initial snapshot of devices from DEP: %v", err))
			} else {
				s.logger.Log("level", "debug", "msg", fmt.Sprintf("More devices after this batch: %t", more))
				if !more {
					s.InitialFetchComplete = true
				}
			}
		} else {
			s.logger.Log("level", "debug", "msg", "Synchronizing devices from DEP service")

			more, err := s.Sync(deviceChan)
			if err != nil {
				s.logger.Log("level", "warn", "msg", fmt.Sprintf("Unable to fetch devices: %s", err))
			} else {
				s.logger.Log("level", "debug", "msg", fmt.Sprintf("More devices after this sync: %t", more))
			}
		}
	}
}

// Pause pauses the synchronisation timer
func (s *syncer) Pause() error {
	s.tickerChanDisabled = s.tickerChan
	s.tickerChan = nil

	return nil
}

// Resume resumes the synchronisation timer
func (s *syncer) Resume() error {
	s.tickerChan = s.tickerChanDisabled
	s.tickerChanDisabled = nil

	return nil
}
