package mqtt

import (
	"fmt"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
)

var MessagePubHandler pahomqtt.MessageHandler = func(client pahomqtt.Client, msg pahomqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var ConnectHandler pahomqtt.OnConnectHandler = func(client pahomqtt.Client) {
	fmt.Println("Connected")

	// Subscribe to the desired topic here
	topic := "classroom/sensor/pir"
	if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Printf("Subscribed to topic: %s\n", topic)
}

var ConnectLostHandler pahomqtt.ConnectionLostHandler = func(client pahomqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}
