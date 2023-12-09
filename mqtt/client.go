package mqtt

import (
	"fmt"

	pahmqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewMQTTClientOptions(broker string, port int) *pahmqtt.ClientOptions {
	opts := pahmqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("emqxs")
	opts.SetPassword("public")
	opts.SetDefaultPublishHandler(MessagePubHandler)
	opts.OnConnect = ConnectHandler
	opts.OnConnectionLost = ConnectLostHandler

	return opts
}
