package device

import (
	"github.com/mosen/devicestore/jsonapi"
	"github.com/satori/go.uuid"
)

// DeviceService provides operations to query and add devices
type DeviceService interface {
	Create(d *Device) (uuid.UUID, *jsonapi.Error)
	//Update()
	//Delete()
}

type deviceService struct {
	store *dataStore
}

func NewService(ds *dataStore) *deviceService {
	return &deviceService{
		store: ds,
	}
}

// Create a device. newly generated uuid is set on the device and returned as the first value
func (svc *deviceService) Create(d *Device) (uuid.UUID, *jsonapi.Error) {
	uuidStr, err := svc.store.Insert(d)
	if err != nil {
		return uuid.Nil, &jsonapi.Error{
			Status: "500",
			Title:  "Error creating device",
			Detail: "",
		}
	}

	uuidObj, err := uuid.FromString(uuidStr)
	if err != nil {
		return uuid.Nil, &jsonapi.Error{
			Status: "500",
			Title:  "Error setting new UUID on device",
			Detail: "",
		}
	}

	d.UUID = uuidObj
	return uuidObj, nil
}

//func (d deviceService) Update() {
//
//}
//
//func (d deviceService) Delete() {
//
//}
