CREATE TABLE data_points (
  device_id VARCHAR NOT NULL,
  device_name VARCHAR NOT NULL,
  property VARCHAR NOT NULL,
  value DOUBLE PRECISION,
  timestamp TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (device_id, property, timestamp)
);

SELECT create_hypertable('data_points', 'timestamp');
