#include <ESP8266WiFi.h>
#include <PubSubClient.h>
#include "mqtt_functions.h"

extern WiFiClient espClient;
extern PubSubClient client;
extern const char *mqtt_server;
extern const char *mqtt_client_id;
extern const char *mqtt_username;
extern const char *mqtt_password;

void callback(char *topic, byte *payload, unsigned int length);

void connectToMQTT()
{
  client.setServer(mqtt_server, 1883);
  client.setCallback(callback);

  while (!client.connected())
  {
    Serial.println("Connecting to MQTT...");
    if (client.connect(mqtt_client_id, mqtt_username, mqtt_password))
    {
      Serial.println("Connected to MQTT");
    }
    else
    {
      Serial.print("Failed, rc=");
      Serial.print(client.state());
      Serial.println(" Retrying in 5 seconds...");
      delay(5000);
    }
  }
}

void mqttLoop()
{
  client.loop();
}
