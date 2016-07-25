package device

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

type dataStore struct {
	*sqlx.DB
}

// Finds a single device
func (d dataStore) Find(uuid uuid.UUID) (*Device, error) {
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

func (d dataStore) FindAll() ([]Device, error) {
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
func (d dataStore) Insert(device *Device) (string, error) {
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
		return "", err
	}

	var uuid string
	d.QueryRow(sql, args...).Scan(&uuid)

	return uuid, nil
}

func (d dataStore) Delete(uuid uuid.UUID) error {
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

func NewDatastore(db *sqlx.DB) *dataStore {
	return &dataStore{
		db,
	}
}
