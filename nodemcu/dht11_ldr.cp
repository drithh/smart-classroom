#include <ESP8266WiFi.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>
#include <DHT.h> // Include the DHT library

#define DEVICE2

#include "config.h"
#include "wifi_functions.h"
#include "mqtt_functions.h"

const char *mqtt_topic_dht11 = "classroom/sensor/dht11";
const char *mqtt_topic_ac = "classroom/input/ac"; // New topic for AC

const int dhtPin = D5; // DHT11 sensor pin
DHT dht(dhtPin, DHT11);

WiFiClient espClient;
PubSubClient client(espClient);

unsigned long lastMsg = 0;
const long interval = 1000; // 1 second interval

void setup()
{
  Serial.begin(115200);

  // Connect to Wi-Fi
  connectToWiFi();

  // Connect to EMQX MQTT broker
  connectToMQTT();

  // Initialize DHT sensor
  dht.begin();

  // Subscribe to the new AC topic
  client.subscribe(mqtt_topic_ac);
}

void loop()
{
  unsigned long currentMillis = millis();

  // Send DHT11 sensor data at most once per second
  if (currentMillis - lastMsg >= interval)
  {
    // Read DHT11 sensor
    float humidity = dht.readHumidity();
    float temperature = dht.readTemperature();

    // Check if readings are valid
    if (!isnan(humidity) && !isnan(temperature))
    {
      // Create a JSON document
      StaticJsonDocument<200> doc;
      doc["humidity"] = humidity;
      doc["temperature"] = temperature;

      // Serialize the JSON document to a char array
      char jsonMsg[200];
      serializeJson(doc, jsonMsg);

      // Publish JSON-formatted DHT11 sensor data
      client.publish(mqtt_topic_dht11, jsonMsg);

      // Update last message time
      lastMsg = currentMillis;
      Serial.println("DHT11 Data Successfully Sent!");
    }
    else
    {
      Serial.println("Failed to read from DHT11 sensor!");
    }
  }

  // Handle MQTT messages
  mqttLoop();
}

void callback(char *topic, byte *payload, unsigned int length)
{
  Serial.println("Message arrived [Topic: " + String(topic) + "]");
  Serial.print("Payload: ");
  for (int i = 0; i < length; i++)
  {
    Serial.print((char)payload[i]);
  }
  Serial.println();

  // Check if the received message is related to the AC topic
  if (strcmp(topic, mqtt_topic_ac) == 0)
  {
    // Process the payload for AC control
    // Example: Implement logic to control the AC based on the received payload
    Serial.println("Received AC control message: " + String((char *)payload));
  }
}
