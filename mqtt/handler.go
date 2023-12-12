package mqtt

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/drithh/smart-classroom/database"
	"github.com/drithh/smart-classroom/types"
	pahomqtt "github.com/eclipse/paho.mqtt.golang"
)

var pirTopic = "classroom/sensor/pir"
var ldrTopic = "classroom/sensor/ldr"
var dht11Topic = "classroom/sensor/dht11"

var setting = types.Setting{
	Ac: types.Device{
		Status: true,
	},
	Lamp1: types.Device{
		Status: true,
	},
	Lamp2: types.Device{
		Status: true,
	},
	Lamp3: types.Device{
		Status: true,
	},
}

var lastUpdatedAcTime int64 = 0

var MessagePubHandler pahomqtt.MessageHandler = func(client pahomqtt.Client, msg pahomqtt.Message) {
	db := database.GetDB()
	switch msg.Topic() {
	case pirTopic:
		marshalledData := types.PirSensorData{}
		err := json.Unmarshal(msg.Payload(), &marshalledData)
		if err != nil {
			fmt.Println("Error unmarshalling data: ", err)
		}
		fmt.Println("PIR sensor data received with value: ", marshalledData.PirStatus)

		_, err = db.Exec("INSERT INTO pir_sensor_data (presence) VALUES ($1)", marshalledData.PirStatus)
		if err != nil {
			fmt.Println("Error inserting data into database: ", err)
		}

		// check if the value is different from the previous one
		if setting.Pir.PirStatus != marshalledData.PirStatus {
			setting.Pir = marshalledData
			fmt.Println("PIR sensor data changed to: ", marshalledData.PirStatus)
			if !setting.Pir.PirStatus {
				if !setting.Lamp1.Status && !setting.Lamp2.Status && !setting.Lamp3.Status {
					fmt.Println("Every lamp is off, returning")

				} else {
					fmt.Println("Turning off every lamp")
					lamps := []types.Device{}
					err = db.Select(&lamps, "SELECT * FROM devices WHERE device_id similar to 'lamp%'")

					if err != nil {
						fmt.Println("Error selecting data from database: ", err)
					}

					for _, lamp := range lamps {
						topic := fmt.Sprintf("classroom/actuator/%s", lamp.DeviceId)

						// make it json
						led := types.Led{
							Led:        false,
							Brightness: 0,
						}

						// marshal to json
						ledJson, err := json.Marshal(led)

						if err != nil {
							fmt.Println("Error marshalling led data: ", err)
						}

						token := client.Publish(topic, 1, false, ledJson)
						token.Wait()

						// update lamp setting value
						switch lamp.DeviceId {
						case "lamp1":
							setting.Lamp1.Status = false
						case "lamp2":
							setting.Lamp2.Status = false
						case "lamp3":
							setting.Lamp3.Status = false
						}
						_, err = db.Exec("UPDATE devices SET status = $1 WHERE device_id = $2", false, lamp.DeviceId)

						if err != nil {
							fmt.Println("Error updating data into database: ", err)
						}
					}
				}

				if !setting.Ac.Status {
					fmt.Println("AC is off, returning")
				} else {
					fmt.Println("Turning off AC")
					topic := "classroom/actuator/ky005"

					// make it json
					ac := types.Ac{
						Ac:          false,
						FanSpeed:    0,
						Temperature: 24,
					}

					// marshal to json
					acJson, err := json.Marshal(ac)

					if err != nil {
						fmt.Println("Error marshalling led data: ", err)
					}

					token := client.Publish(topic, 1, false, acJson)
					token.Wait()

					// update lamp setting value
					setting.Ac.Status = false
					_, err = db.Exec("UPDATE devices SET status = $1 WHERE device_id = $2", false, "ac")

					if err != nil {
						fmt.Println("Error updating data into database: ", err)
					}
				}
			}
		}
	case ldrTopic:
		marshalledData := types.LdrSensorData{}
		err := json.Unmarshal(msg.Payload(), &marshalledData)
		if err != nil {
			fmt.Println("Error unmarshalling data: ", err)
		}
		fmt.Println("LDR sensor data received with value: ", marshalledData.Brightness)

		_, err = db.Exec("INSERT INTO ldr_sensor_data (light_intensity) VALUES ($1)", marshalledData.Brightness)

		if err != nil {
			fmt.Println("Error inserting data into database: ", err)
		}

		// if settings.pir is 0 then return
		if !setting.Pir.PirStatus {
			return
		}

		// select lamp
		lamps := []types.DeviceSetting{}
		err = db.Select(&lamps, "SELECT * FROM device_settings WHERE device_id similar to 'lamp%'")

		if err != nil {
			fmt.Println("Error selecting data from database: ", err)
		}

		// map brightness value from 0-1023 to 0-100
		brightness := 100 - int(float32(marshalledData.Brightness)/1023*100)

		for _, lamp := range lamps {
			topic := fmt.Sprintf("classroom/actuator/%s", lamp.DeviceId)

			// average brightness value from lamp setting and expected brightness value
			// convert lamp.SettingValue to int
			lampBrightness, err := strconv.Atoi(lamp.SettingValue)

			if err != nil {
				fmt.Println("Error converting lamp setting value to int: ", err)
			}

			// if lampbrightness and brightness differs abouot 20
			if lampBrightness-brightness > 15 {
				// update lamp setting value

				expectedBrightness := lampBrightness * 255 / 100
				fmt.Println("Brightness value for ", lamp.DeviceId, " is: ", expectedBrightness, "with lamp setting value: ", lampBrightness, "and ldr value: ", brightness)

				// make it json
				led := types.Led{
					Led:        true,
					Brightness: expectedBrightness,
				}

				// marshal to json
				ledJson, err := json.Marshal(led)

				if err != nil {
					fmt.Println("Error marshalling led data: ", err)
				}

				token := client.Publish(topic, 1, false, ledJson)
				token.Wait()

				switch lamp.DeviceId {
				case "lamp1":
					setting.Lamp1.Status = true
				case "lamp2":
					setting.Lamp2.Status = true
				case "lamp3":
					setting.Lamp3.Status = true
				}
				_, err = db.Exec("UPDATE devices SET status = $1 WHERE device_id = $2", true, lamp.DeviceId)

				if err != nil {
					fmt.Println("Error updating data into database: ", err)
				}
			} else {
				fmt.Println("Brightness value for ", lamp.DeviceId, " is: ", brightness, "with lamp setting value: ", lampBrightness, "and ldr value: ", brightness)
				fmt.Println("No need to update lamp setting value")
			}
		}

	case dht11Topic:
		marshalledData := types.Dht11SensorData{}
		err := json.Unmarshal(msg.Payload(), &marshalledData)
		if err != nil {
			fmt.Println("Error unmarshalling data: ", err)
		}
		fmt.Println("DHT11 sensor data received with temperature: ", marshalledData.Temperature, " and humidity: ", marshalledData.Humidity)

		_, err = db.Exec("INSERT INTO dht11_sensor_data (temperature, humidity) VALUES ($1, $2)", marshalledData.Temperature, marshalledData.Humidity)
		if err != nil {
			fmt.Println("Error inserting data into database: ", err)
		}

		currentTime := time.Now().Unix()
		// if settings.pir is 0 then return
		if !setting.Pir.PirStatus || currentTime-lastUpdatedAcTime < 10 {
			return
		}

		// select ac
		acs := []types.DeviceSetting{}
		err = db.Select(&acs, "SELECT * FROM device_settings WHERE device_id similar to 'ac%'")
		if err != nil {
			fmt.Println("Error selecting data from database: ", err)
		}

		type ACSetting struct {
			Ac          bool `json:"ac"`
			Temperature int  `json:"temperature"`
			FanSpeed    int  `json:"fan_speed"`
			Swing       bool `json:"swing"`
		}

		acSetting := ACSetting{}

		for _, ac := range acs {
			switch ac.SettingName {
			case "temperature":
				acSetting.Temperature, err = strconv.Atoi(ac.SettingValue)
				if err != nil {
					fmt.Println("Error converting ac setting value to int: ", err)
				}
			case "fan_speed":
				acSetting.FanSpeed, err = strconv.Atoi(ac.SettingValue)
				if err != nil {
					fmt.Println("Error converting ac setting value to int: ", err)
				}
			case "swing":
				switch ac.SettingValue {
				case "on":
					acSetting.Swing = true
				case "off":
					acSetting.Swing = false
				}
			}
		}

		// publish ac setting
		topic := "classroom/actuator/ky005"

		// make it json
		ac := types.Ac{
			Ac:          true,
			FanSpeed:    acSetting.FanSpeed,
			Temperature: acSetting.Temperature,
		}

		// marshal to json
		acJson, err := json.Marshal(ac)

		if err != nil {
			fmt.Println("Error marshalling led data: ", err)
		}

		token := client.Publish(topic, 1, false, acJson)
		token.Wait()

		// update ac setting value
		setting.Ac.Status = true
		_, err = db.Exec("UPDATE devices SET status = $1 WHERE device_id = $2", true, "ac")

		if err != nil {
			fmt.Println("Error updating data into database: ", err)
		}

		// update lastUpdatedAcTime to current time
		lastUpdatedAcTime = time.Now().Unix()
		fmt.Println("AC setting updated to: ", lastUpdatedAcTime)

	default:
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	}
}

var ConnectHandler pahomqtt.OnConnectHandler = func(client pahomqtt.Client) {
	fmt.Println("Connected")

	// Subscribe to the desired topic here
	if token := client.Subscribe(pirTopic, 1, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Printf("Subscribed to topic: %s\n", pirTopic)

	if token := client.Subscribe(ldrTopic, 1, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Printf("Subscribed to topic: %s\n", ldrTopic)

	if token := client.Subscribe(dht11Topic, 1, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Printf("Subscribed to topic: %s\n", dht11Topic)
}

var ConnectLostHandler pahomqtt.ConnectionLostHandler = func(client pahomqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}
