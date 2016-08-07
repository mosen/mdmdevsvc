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
	ins := sq.Insert("devices").Columns(
		"udid",
		"serial_number",
		"imei",
		"meid",
		"product_name",
		"model",
		"asset_tag",
		"color",
		"os_version",
	).Values("udid", device.UDID).
		Values("serial_number", device.SerialNumber).
		Values("imei", device.IMEI).
		Values("meid", device.MEID).
		Values("product_name", device.ProductName).
		Values("model", device.Model).
		Values("asset_tag", device.AssetTag).
		Values("color", device.Color).
		Values("os_version", device.OSVersion).Suffix("RETURNING \"uuid\"")

	sql, args, err := ins.ToSql()
	if err != nil {
		return err
	}

	d.Logger.Log(sql)

	var uuidStr string
	d.QueryRow(sql, args...).Scan(&uuidStr)
	device.UUID, err = uuid.FromString(uuidStr)
	if err != nil {
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
