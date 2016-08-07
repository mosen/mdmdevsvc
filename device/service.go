package device

import (
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"errors"
)

// DeviceService provides operations to query and add devices
type Service interface {
	PostDevice(ctx context.Context, d Device) (uuid.UUID, error)
	GetDevice(ctx context.Context, uuidStr string) (Device, error)
	PutDevice(ctx context.Context, uuidStr string, d Device) error
	PatchDevice(ctx context.Context, uuidStr string, d Device) error
	DeleteDevice(ctx context.Context, uuidStr string) error
}

var (
	ErrMalformedUUID = errors.New("malformed UUID")
	ErrNotFound	= errors.New("device not found")
	ErrQueryError	= errors.New("error performing query")
)

type service struct {
	store DeviceRepository
}

func NewService(ds DeviceRepository) Service {
	return &service{
		store: ds,
	}
}

// Create a device. newly generated uuid is set on the device and returned as the first value
func (s *service) PostDevice(ctx context.Context, d Device) (uuid.UUID, error) {
	if err := s.store.Store(&d); err != nil {
		return uuid.Nil, ErrQueryError
	}

	return d.UUID, nil
}

func (s *service) GetDevice(ctx context.Context, uuidStr string) (Device, error) {
	uuidObj, err := uuid.FromString(uuidStr)
	if err != nil {
		return Device{}, ErrMalformedUUID
	}

	device, err := s.store.Find(uuidObj)
	if err != nil {
		return Device{}, ErrQueryError
	}

	if device == nil {
		return Device{}, ErrNotFound
	} else {
		return *device, nil
	}
}

func (s *service) PutDevice(ctx context.Context, uuidStr string, d Device) error {
	return errors.New("Not Implemented")
}

func (s *service) PatchDevice(ctx context.Context, uuidStr string, d Device) error {
	return errors.New("Not Implemented")
}


func (s *service) DeleteDevice(ctx context.Context, uuidStr string) error {
	uuidObj, err := uuid.FromString(uuidStr)
	if err != nil {
		return ErrMalformedUUID
	}

	if err := s.store.Delete(uuidObj); err != nil {
		return err
	} else {
		return nil
	}
}