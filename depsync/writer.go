package depsync

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/micromdm/dep"
	"github.com/mosen/devicestore/device"
	//"github.com/satori/go.uuid"
	"time"
)

type Writer interface {
	Start(deviceChan <-chan dep.Device)
	Write(*dep.Device) error
}

type writer struct {
	repository device.DeviceRepository
	logger     log.Logger
}

func NewWriter(repository device.DeviceRepository, logger log.Logger) Writer {
	return &writer{
		repository: repository,
		logger:     logger,
	}
}

func (w *writer) Start(deviceChan <-chan dep.Device) {
	for dev := range deviceChan {
		//var opType string = dev.OpType
		//var opDate time.Time = dev.OpDate

		switch dev.OpType {
		case "added":
			w.logger.Log("level", "debug", "msg", "writing added device to database")
			if err := w.Write(&dev); err != nil {
				w.logger.Log("level", "error", "msg", fmt.Sprintf("Failed to write DEP device: %s", err))
			}
		case "modified":
		case "deleted":
		default:
			w.logger.Log("level", "debug", "msg", "writing fetched device to database")
			if err := w.Write(&dev); err != nil {
				w.logger.Log("level", "error", "msg", fmt.Sprintf("Failed to write DEP device: %s", err))
			}
		}

	}
}

func (w *writer) Write(dev *dep.Device) error {
	//depProfileUUID, err := uuid.FromString(dev.ProfileUUID)
	//if err != nil {
	//	return err
	//}

	var created, updated time.Time = time.Now(), time.Now()

	device := &device.Device{
		SerialNumber: dev.SerialNumber,
		Model:        dev.Model,
		Description:  dev.Description,
		Color:        dev.Color,
		AssetTag:     dev.AssetTag,

		DepProfileStatus: dev.ProfileStatus,
		//DepProfileUUID: depProfileUUID,
		DepProfileAssignTime:   dev.ProfileAssignTime,
		DepProfilePushTime:     dev.ProfilePushTime,
		DepProfileAssignedDate: dev.DeviceAssignedDate,
		DepProfileAssignedBy:   dev.DeviceAssignedBy,

		Created: &created,
		Updated: &updated,
		HasDEP:  true,
	}

	if err := w.repository.Store(device); err != nil {
		return err
	}

	return nil
}
