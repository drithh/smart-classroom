CREATE TABLE pir_sensor_data (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP default CURRENT_TIMESTAMP,
    presence BOOLEAN
);

CREATE TABLE dht11_sensor_data (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP default CURRENT_TIMESTAMP,
    temperature DECIMAL(5, 2),
    humidity DECIMAL(5, 2)
);

CREATE TABLE ldr_sensor_data (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP default CURRENT_TIMESTAMP,
    light_intensity INTEGER
);

CREATE TABLE devices (
    id SERIAL PRIMARY KEY,
    device_id VARCHAR(255) UNIQUE,
    status boolean
);

CREATE TABLE device_settings (
    id SERIAL PRIMARY KEY,
    device_id VARCHAR(255) NOT NULL REFERENCES devices(device_id),
    setting_name VARCHAR(255),
    setting_value VARCHAR(255)
);

INSERT INTO devices (device_id, status) VALUES
('ac1', true),
('lamp1', true),
('lamp2', true),
('lamp3', true);

-- AC Settings
INSERT INTO device_settings (device_id, setting_name, setting_value) VALUES
('ac1', 'temperature', '23'),
('ac1', 'fan_speed', 'medium'),
('ac1', 'swing', 'on');

-- Lamp Settings
INSERT INTO device_settings (device_id, setting_name, setting_value) VALUES
('lamp1', 'brightness', '50'),
('lamp2', 'brightness', '50'),
('lamp3', 'brightness', '50');
