CREATE TABLE IF NOT EXISTS sensors (
    id1 VARCHAR(64) PRIMARY KEY,       -- Unique sensor ID
    sensor_type VARCHAR(100) NOT NULL  -- Sensor type (e.g., temperature, humidity)
);