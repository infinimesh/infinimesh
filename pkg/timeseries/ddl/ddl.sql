CREATE TABLE data_points_test (
  device_id VARCHAR NOT NULL,
  message_id BIGINT NOT NULL,
  property VARCHAR NOT NULL,
  value DOUBLE PRECISION,
  timestamp TIMESTAMPTZ NOT NULL,
  message_length DOUBLE PRECISION,
  PRIMARY KEY (device_id, property, timestamp)
);

