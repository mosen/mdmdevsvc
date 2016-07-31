CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS devices (
  uuid uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  udid uuid,
  name text NOT NULL DEFAULT '',
  description text NOT NULL DEFAULT '',
  serial_number text,
  imei text,
  meid text,
  os_version text NOT NULL DEFAULT '',
  manufacturer text NOT NULL DEFAULT '',
  vendor text NOT NULL DEFAULT '',
  model text NOT NULL DEFAULT '',
  color text NOT NULL DEFAULT '',
  asset_tag text NOT NULL DEFAULT '',
  warranty_expires date,
  build_version text NOT NULL DEFAULT '',
  product_name text NOT NULL DEFAULT '',

  -- DEP only
  dep_device boolean,
  dep_profile_status text,
  dep_profile_uuid text,
  dep_profile_assign_time date,
  dep_profile_push_time date,
  dep_profile_assigned_date date,
  dep_profile_assigned_by text,


  apple_mdm_token text,
  apple_mdm_topic text,
  apple_push_magic text,
  apple_mdm_enrolled boolean,
  apple_unlock_token BYTEA,

  awaiting_configuration boolean,
  last_checkin timestamp,


  created_at date NOT NULL DEFAULT now(),
  updated_at date,
  deleted_at date
);

CREATE UNIQUE INDEX IF NOT EXISTS serial_idx ON devices (serial_number);
CREATE UNIQUE INDEX IF NOT EXISTS udid_idx ON devices (udid);
CREATE UNIQUE INDEX IF NOT EXISTS imei_idx ON devices (imei);
CREATE UNIQUE INDEX IF NOT EXISTS meid_idx ON devices (meid);


