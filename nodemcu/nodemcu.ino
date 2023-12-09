#include <ESP8266WiFi.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>

#define DEVICE1

#include "config.h"
#include "wifi_functions.h"
#include "mqtt_functions.h"

const char *mqtt_topic_pir = "classroom/sensor/pir";
const char *mqtt_topic_ac = "classroom/input/ac";

const int pirPin = D5; // PIR sensor pin
bool lastPirValue = false;

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
}

void loop()
{
  unsigned long currentMillis = millis();

  // Send message at most once per second
  if (currentMillis - lastMsg >= interval)
  {
    // Read PIR sensor
    int pirValue = digitalRead(pirPin);
    if (pirValue != lastPirValue)
    {
      Serial.println("PIR sensor status changed");
      StaticJsonDocument<200> doc;
      doc["pirStatus"] = pirValue;
      char jsonMsg[200];
      serializeJson(doc, jsonMsg);
      client.publish(mqtt_topic_pir, jsonMsg);
      lastPirValue = pirValue;
    }
    // Create a JSON document
    StaticJsonDocument<200> doc;
    doc["pirStatus"] = pirValue;

    // Serialize the JSON document to a char array
    char jsonMsg[200];
    serializeJson(doc, jsonMsg);

    // Publish JSON-formatted PIR sensor status
    client.publish(mqtt_topic_pir, jsonMsg);

    // Update last message time
    lastMsg = currentMillis;
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
}
