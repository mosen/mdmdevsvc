package device

import (
	sq "github.com/Masterminds/squirrel"
	kitlog "github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
	"os"
)

type DeviceRepository interface {
	Find(uuid uuid.UUID) (*Device, error)
	FindAll() ([]Device, error)
	Store(device *Device) error
	Delete(uuid uuid.UUID) error
	Update(device *Device) error
}

type deviceRepository struct {
	*sqlx.DB
	kitlog.Logger
}

// Finds a single device
func (d *deviceRepository) Find(uuid uuid.UUID) (*Device, error) {
	device := sq.Select("*").From("devices").Where(sq.Eq{"uuid": uuid.String()})
	sql, args, err := device.ToSql()
	if err != nil {
		return nil, err
	}

	var result Device
	if err := d.Get(&result, sql, args...); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *deviceRepository) FindAll() ([]Device, error) {
	stmt := sq.Select("*").From("devices")
	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}

	var result []Device
	if err := d.Select(&result, sql, args...); err != nil {
		return nil, err
	}

	return result, nil
}

// Inserts a device and returns its uuid
func (d *deviceRepository) Store(device *Device) error {
	query, args, err := sq.Insert("devices").
		Columns(
		"udid",
		"name",
		"serial_number",
		"imei",
		"meid",
		"model",
		"description",
		"asset_tag",
		"has_dep",
		"color",
		"dep_profile_status",
		"dep_profile_assigned_date",
		"dep_profile_assigned_by",
		).
		Values(
		device.UDID,
		device.Name,
		device.SerialNumber,
		device.IMEI,
		device.MEID,
		device.Model,
		device.Description,
		device.AssetTag,
		device.HasDEP,
		device.Color,
		device.DepProfileStatus,
		device.DepProfileAssignedDate,
		device.DepProfileAssignedBy,
		).
		Suffix("RETURNING \"uuid\", \"created_at\", \"updated_at\"").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	if err := d.QueryRow(query, args...).Scan(&device.UUID, &device.Created, &device.Updated); err != nil {
		return err
	}

	return nil
}

func (d *deviceRepository) Update(device *Device) error {
	stmt := sq.Update("devices").SetMap(
		sq.Eq{
			"name": device.Name,
		},
	).Where(sq.Eq{"uuid": device.UUID.String()})

	sql, args, err := stmt.ToSql()
	if err != nil {
		return err
	}

	if _, err := d.Exec(sql, args...); err != nil {
		return err
	}

	return nil
}

//func (d *dataStore) Upsert(device *Device) (bool, error) {
//
//}

func (d *deviceRepository) Delete(uuid uuid.UUID) error {
	stmt := sq.Delete("devices").Where(sq.Eq{"uuid": uuid.String()})
	sql, args, err := stmt.ToSql()
	if err != nil {
		return err
	}

	if _, err := d.Exec(sql, args...); err != nil {
		return err
	}

	return nil
}

func NewRepository(db *sqlx.DB) DeviceRepository {
	return &deviceRepository{
		db,
		kitlog.NewLogfmtLogger(os.Stdout),
	}
}
