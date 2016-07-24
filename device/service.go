package device

// DeviceService provides operations to query and add devices
type DeviceService interface {
	Create()
	Update()
	Delete()
}


type deviceService struct{}

func NewService() *deviceService {
	return &deviceService{}
}

func (d deviceService) Create() {

}
