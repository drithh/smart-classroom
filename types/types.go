package types

import "time"

type PirSensorData struct {
	Timestamp time.Time `db:"timestamp"`
	PirStatus bool      `db:"presence"`
}

type LdrSensorData struct {
	Brightness int `db:"light_intensity"`
}

type Dht11SensorData struct {
	Temperature float32 `db:"temperature"`
	Humidity    float32 `db:"humidity"`
}

type Led struct {
	Led        bool `json:"led"`
	Brightness int  `json:"brightness"`
}

type Ac struct {
	Ac          bool `json:"ac"`
	FanSpeed    int  `json:"fan_speed"`
	Temperature int  `json:"temperature"`
}

type DeviceSetting struct {
	DeviceId     string `db:"device_id"`
	SettingName  string `db:"setting_name"`
	SettingValue string `db:"setting_value"`
}

type Device struct {
	DeviceId string `db:"device_id"`
	Status   bool   `db:"status"`
	Setting  []DeviceSetting
}

type Setting struct {
	Pir   PirSensorData
	Ldr   LdrSensorData
	Dht11 Dht11SensorData
	Ac    Device
	Lamp1 Device
	Lamp2 Device
	Lamp3 Device
}
