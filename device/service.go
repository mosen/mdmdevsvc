package device

// DeviceService provides operations to query and add devices
type DeviceService interface {
	Create()
	Update()
	Delete()
}

type deviceService struct {
	store dataStore
}

func NewService(ds *dataStore) *deviceService {
	return &deviceService{
		store: ds,
	}
}

func (d deviceService) Create() {

}

func (d deviceService) Update() {

}

func (d deviceService) Delete() {

}
