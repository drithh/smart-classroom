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

CREATE TABLE device_settings (
    id SERIAL PRIMARY KEY,
    device_id VARCHAR(255),
    setting_name VARCHAR(255),
    setting_value VARCHAR(255)
);

-- AC Settings
INSERT INTO device_settings (device_id, setting_name, setting_value) VALUES
('ac1', 'temperature', '23'),
('ac1', 'fan_speed', 'medium'),
('ac1', 'swing', 'on');

-- Lamp Settings
INSERT INTO device_settings (device_id, setting_name, setting_value) VALUES
('lamp1', 'brightness', '50');
