// main.go
package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/drithh/smart-classroom/fiber"
	mqtt "github.com/drithh/smart-classroom/mqtt"
	pahomqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	fiber.SetupFiber()

	var broker = "localhost"
	var port = 1883

	opts := mqtt.NewMQTTClientOptions(broker, port)

	client := pahomqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Wait for a signal to interrupt the program
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	// Graceful shutdown
	fmt.Println("Disconnecting from MQTT broker...")
	client.Disconnect(250)
	fmt.Println("Disconnected")
	os.Exit(0)

}
