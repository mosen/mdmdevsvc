package device

import (
	"github.com/satori/go.uuid"
)

type Device struct {
	UUID uuid.UUID `json:"uuid,omitempty"`
	UDID string `json:"udid,omitempty"`
	SerialNumber string `json:"serial_number,omitempty"`
	IMEI string `json:"imei,omitempty"`
	MEID string `json:"meid,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	Model string `json:"model,omitempty"`
	AssetTag string `json:"asset_tag,omitempty"`
	Color string `json:"color,omitempty"`
	OSVersion string `json:"os_version,omitempty"`
}
