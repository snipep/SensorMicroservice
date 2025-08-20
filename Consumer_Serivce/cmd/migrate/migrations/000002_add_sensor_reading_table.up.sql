CREATE TABLE IF NOT EXISTS sensor_readings (
    id BIGINT AUTO_INCREMENT PRIMARY KEY, -- Unique row ID for readings
    id1 VARCHAR(64) NOT NULL,             -- FK to sensors.id1
    id2 INT NOT NULL,                     -- Data entry point ID
    sensor_value FLOAT NOT NULL,          -- Sensor measurement
    timestamp DATETIME NOT NULL,          -- When data was recorded

    CONSTRAINT fk_sensor FOREIGN KEY (id1)
        REFERENCES sensors(id1)
        ON DELETE CASCADE                 -- Delete readings if sensor is deleted
);