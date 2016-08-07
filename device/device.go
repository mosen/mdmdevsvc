package device

import (
	"github.com/satori/go.uuid"
	"time"
)

type Device struct {
	UUID            uuid.UUID `json:"uuid,omitempty"`
	UDID            uuid.UUID `json:"udid,omitempty"`
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
	LastCheckin     time.Time `json:"last_checkin,omitempty"`

	// DEP
	HasDEP                 bool      `json:"has_dep"`
	DepProfileStatus       string    `json:"dep_profile_status,omitempty"`
	DepProfileUUID         uuid.UUID `json:"dep_profile_uuid,omitempty"`
	DepProfileAssignTime   time.Time `json:"dep_profile_assign_time,omitempty"`
	DepProfilePushTime     time.Time `json:"dep_profile_push_time,omitempty"`
	DepProfileAssignedDate time.Time `json:"dep_profile_assigned_date,omitempty"`
	DepProfileAssignedBy   string    `json:"dep_profile_assigned_by,omitempty"`

	// Apple
	AppleMDMToken              string `json:"apple_mdm_token,omitempty"`
	AppleMDMTopic              string `json:"apple_mdm_topic,omitempty"`
	ApplePushMagic             string `json:"apple_push_magic,omitempty"`
	AppleMDMEnrolled           bool   `json:"apple_mdm_enrolled"`
	AppleUnlockToken           string `json:"apple_unlock_token,omitempty"`
	AppleAwaitingConfiguration bool   `json:"apple_awaiting_configuration,omitempty"`
}
