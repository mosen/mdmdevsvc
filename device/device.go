package device

import (
	"github.com/satori/go.uuid"
	"time"
)

type Device struct {
	UUID            uuid.UUID `json:"uuid,omitempty"`
	UDID            string    `json:"udid,omitempty"`
	SerialNumber    string    `json:"serial_number,omitempty"`
	IMEI            string    `json:"imei,omitempty"`
	MEID            string    `json:"meid,omitempty"`
	Created         time.Time `json:"created_at,omitempty"`
	Updated         time.Time `json:"updated_at,omitempty"`
	Deleted         time.Time `json:"deleted_at,omitempty"`
	Name            string    `json:"name,omitempty"`
	Description     string    `json:"description,omitempty"`
	Manufacturer    string    `json:"manufacturer,omitempty"`
	Vendor          string    `json:"vendor,omitempty"`
	ProductName     string    `json:"product_name,omitempty"`
	Model           string    `json:"model,omitempty"`
	AssetTag        string    `json:"asset_tag,omitempty"`
	WarrantyExpires time.Time `json:"warranty_expires,omitempty"`
	Color           string    `json:"color,omitempty"`
	OSVersion       string    `json:"os_version,omitempty"`
}
