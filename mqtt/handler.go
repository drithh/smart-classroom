package mqtt

import (
	"encoding/json"
	"fmt"

	"github.com/drithh/smart-classroom/database"
	pahomqtt "github.com/eclipse/paho.mqtt.golang"
)

var pirTopic = "classroom/sensor/pir"
var ldrTopic = "classroom/sensor/ldr"
var dht11Topic = "classroom/sensor/dht11"

type PirSensorData struct {
	PirStatus int
}

type LdrSensorData struct {
	LdrStatus int
}

type Dht11SensorData struct {
	Temperature float32
	Humidity    float32
}

var MessagePubHandler pahomqtt.MessageHandler = func(client pahomqtt.Client, msg pahomqtt.Message) {
	db := database.GetDB()
	switch msg.Topic() {
	case pirTopic:
		marshalledData := PirSensorData{}
		err := json.Unmarshal(msg.Payload(), &marshalledData)
		if err != nil {
			fmt.Println("Error unmarshalling data: ", err)
		}
		fmt.Println("PIR sensor data received with value: ", marshalledData.PirStatus)

		status := marshalledData.PirStatus == 1

		_, err = db.Exec("INSERT INTO pir_sensor_data (presence) VALUES ($1)", status)
		if err != nil {
			fmt.Println("Error inserting data into database: ", err)
		}
	case ldrTopic:
		fmt.Println("LDR sensor data received")
	case dht11Topic:
		marshalledData := Dht11SensorData{}
		err := json.Unmarshal(msg.Payload(), &marshalledData)
		if err != nil {
			fmt.Println("Error unmarshalling data: ", err)
		}
		fmt.Println("DHT11 sensor data received with temperature: ", marshalledData.Temperature, " and humidity: ", marshalledData.Humidity)

		_, err = db.Exec("INSERT INTO dht11_sensor_data (temperature, humidity) VALUES ($1, $2)", marshalledData.Temperature, marshalledData.Humidity)
		if err != nil {
			fmt.Println("Error inserting data into database: ", err)
		}
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
