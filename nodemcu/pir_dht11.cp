#include <ESP8266WiFi.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>
#include <DHT.h>

#define DEVICE1

#include "config.h"
#include "wifi_functions.h"
#include "mqtt_functions.h"

const char *mqtt_topic_pir = "classroom/sensor/pir";
const char *mqtt_topic_dht11 = "classroom/sensor/dht11";

const int pirPin = D5; // PIR sensor pin
const int dhtPin = D7; // DHT11 sensor pin

DHT dht(dhtPin, DHT11);

WiFiClient espClient;
PubSubClient client(espClient);

unsigned long lastMsg = 0;
const long interval = 1000;

void setup()
{
  Serial.begin(115200);

  connectToWiFi();
  connectToMQTT();

  dht.begin();
}

void loop()
{
  unsigned long currentMillis = millis();

  if (currentMillis - lastMsg >= interval)
  {
    int pirValue = digitalRead(pirPin);
    float humidity = dht.readHumidity();
    float temperature = dht.readTemperature();

    if (!isnan(humidity) && !isnan(temperature))
    {
      // Publish PIR sensor data
      StaticJsonDocument<200> pirDoc;
      pirDoc["pirStatus"] = pirValue == 1 ? true : false;
      char pirJsonMsg[200];
      serializeJson(pirDoc, pirJsonMsg);
      client.publish(mqtt_topic_pir, pirJsonMsg);

      // Publish DHT11 sensor data
      StaticJsonDocument<200> dhtDoc;
      dhtDoc["humidity"] = humidity;
      dhtDoc["temperature"] = temperature;
      char dhtJsonMsg[200];
      serializeJson(dhtDoc, dhtJsonMsg);
      client.publish(mqtt_topic_dht11, dhtJsonMsg);

      lastMsg = currentMillis;
    }
    else
    {
      Serial.println("Failed to read from DHT11 sensor!");
    }
  }

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
