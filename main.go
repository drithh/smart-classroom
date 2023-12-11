// main.go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"github.com/drithh/smart-classroom/database"
	"github.com/drithh/smart-classroom/fiber"
	mqtt "github.com/drithh/smart-classroom/mqtt"
	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
	database.ConnectDB()
	defer database.CloseDB()

	db := database.GetDB()

	fiber.SetupFiber(db)

	broker := os.Getenv("BROKER_HOST")
	port := os.Getenv("BROKER_PORT")

	if broker == "" || port == "" {
		fmt.Println("Error: BROKER_HOST and BROKER_PORT environment variables must be set")
		os.Exit(1)
	}

	intPort, err := strconv.Atoi(port)

	if err != nil {
		fmt.Println("Error: BROKER_PORT must be an integer")
		os.Exit(1)
	}

	opts := mqtt.NewMQTTClientOptions(broker, intPort)

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
