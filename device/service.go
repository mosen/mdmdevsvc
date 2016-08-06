package device

import (
	"github.com/mosen/devicestore/jsonapi"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

// DeviceService provides operations to query and add devices
type Service interface {
	PostDevice(ctx context.Context, d Device) (uuid.UUID, *jsonapi.Error)
	GetDevice(ctx context.Context, uuidStr string) (Device, *jsonapi.Error)
	PutDevice(ctx context.Context, uuidStr string, d Device) error
	PatchDevice(ctx context.Context, uuidStr string, d Device) error
	DeleteDevice(ctx context.Context, uuidStr string) error
}

type service struct {
	store Datastore
}

func NewService(ds Datastore) *service {
	return &service{
		store: ds,
	}
}

// Create a device. newly generated uuid is set on the device and returned as the first value
func (s *service) PostDevice(ctx context.Context, d *Device) (uuid.UUID, *jsonapi.Error) {
	uuidObj, err := s.store.Insert(d)
	if err != nil {
		return uuid.Nil, &jsonapi.Error{
			Status: "500",
			Title:  "Error creating device",
			Detail: "",
		}
	}

	d.UUID = uuidObj
	return uuidObj, nil
}

func (s *service) GetDevice(ctx context.Context, uuidStr string) (Device, error) {
	uuidObj, err := uuid.FromString(uuidStr)
	if err != nil {
		return nil, &jsonapi.Error{
			Status: "400",
			Title: "Malformed UUID",
			Detail: "",
		}
	}

	device, err := s.store.Find(uuidObj)
	if err != nil {
		return nil, &jsonapi.Error{
			Status: "500",
			Title: "Query Error",
			Detail: err,
		}
	}

	if device == nil {
		return nil, &jsonapi.Error{
			Status: "404",
			Title: "Device not found",
			Detail: "",
		}
	} else {
		return device, nil
	}
}
