#include <ESP8266WiFi.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>
#include <Arduino.h>
#include <IRremoteESP8266.h>
#include <IRsend.h>
#include <ir_Panasonic.h>

#define DEVICE2

#include "config.h"
#include "wifi_functions.h"
#include "mqtt_functions.h"

const char *mqtt_topic_ldr = "classroom/sensor/ldr";
const char *mqtt_topic_ky005 = "classroom/actuator/ky005";
const char *mqtt_topic_lamp1 = "classroom/actuator/lamp1";
const char *mqtt_topic_lamp2 = "classroom/actuator/lamp2";
const char *mqtt_topic_lamp3 = "classroom/actuator/lamp3";

const int ldrPin = A0; // LDR sensor pin

const uint16_t ky005Pin = D7; // ESP8266 GPIO pin to use. Recommended: 4 (D2).
IRPanasonicAc ac(ky005Pin);   // Set the GPIO used for sending messages.

WiFiClient espClient;
PubSubClient client(espClient);

unsigned long lastMsg = 0;
const long interval = 1000; // 1 second interval

void printState()
{
  // Display the settings.
  Serial.println("Panasonic A/C remote is in the following state:");
  Serial.printf("  %s\n", ac.toString().c_str());
  // Display the encoded IR sequence.
  unsigned char *ir_code = ac.getRaw();
  Serial.print("IR Code: 0x");
  for (uint8_t i = 0; i < kPanasonicAcStateLength; i++)
    Serial.printf("%02X", ir_code[i]);
  Serial.println();
}

void setup()
{
  Serial.begin(115200);

  // Connect to Wi-Fi
  connectToWiFi();

  // Connect to EMQX MQTT broker
  connectToMQTT();

  // Subscribe to the new AC topic
  client.subscribe(mqtt_topic_ky005);

  // Subscribe to lamp topics
  client.subscribe(mqtt_topic_lamp1);
  client.subscribe(mqtt_topic_lamp2);
  client.subscribe(mqtt_topic_lamp3);

  pinMode(D0, OUTPUT); // Lamp 1
  pinMode(D1, OUTPUT); // Lamp 2
  pinMode(D2, OUTPUT); // Lamp 3

  ac.begin();
  delay(200);

  // Set up what we want to send. See ir_Panasonic.cpp for all the options.
  Serial.println("Default state of the remote.");
  printState();
  Serial.println("Setting desired state for A/C.");
  ac.setModel(kPanasonicRkr);
  ac.on();
}

void loop()
{
  unsigned long currentMillis = millis();

  // Send LDR sensor data at most once per second
  if (currentMillis - lastMsg >= interval)
  {
    // Read LDR sensor
    int ldrValue = analogRead(ldrPin);

    // Create a JSON document
    StaticJsonDocument<200> doc;
    doc["ldrValue"] = ldrValue;

    // Serialize the JSON document to a char array
    char jsonMsg[200];
    serializeJson(doc, jsonMsg);

    // Publish JSON-formatted LDR sensor data
    client.publish(mqtt_topic_ldr, jsonMsg);

    // Update last message time
    lastMsg = currentMillis;
    // Serial.println("LDR Data Successfully Sent!");
  }

  // Handle MQTT messages
  mqttLoop();
}

void processLampControl(const JsonDocument &doc, int pin)
{
  bool state = doc["led"];

  // Check the state and perform the corresponding action
  if (state)
  {
    // Turn on the lamp (assuming the specified pin)
    int brightness = doc["brightness"];

    // Check if brightness is within the valid range (0 to 255)
    if (brightness >= 0 && brightness <= 255)
    {
      // Set the brightness using analogWrite
      analogWrite(pin, brightness);
      Serial.print("Lamp brightness set to: ");
      Serial.println(brightness);
    }
    else
    {
      Serial.println("Invalid brightness value. It should be between 0 and 255.");
    }
    Serial.println("Lamp turned ON");
  }
  else
  {
    // Turn off the lamp
    digitalWrite(pin, LOW);
    Serial.println("Lamp turned OFF");
  }
}

void processKY005Control(const JsonDocument &doc, int pin)
{
  bool state = doc["status"];

  // Check the state and perform the corresponding action
  if (state)
  {
    // Turn on the KY005 device (assuming the specified pin)
    // Add your KY005-specific code here
    // digitalWrite(pin, HIGH);
    ac.on();
    ac.setFan(kPanasonicAcFanAuto);
    ac.setMode(kPanasonicAcCool);
    ac.setTemp(26);
    ac.setSwingVertical(kPanasonicAcSwingVAuto);
    ac.setSwingHorizontal(kPanasonicAcSwingHAuto);
    Serial.println("Turning ON AC");
    ac.send();
  }
  else
  {
    // digitalWrite(pin, LOW);
    ac.off();
    Serial.println("Turning OFF AC");
  }
  printState();
}

// In your callback function:
void callback(char *topic, byte *payload, unsigned int length)
{
  Serial.println("Message arrived [Topic: " + String(topic) + "]");
  Serial.print("Payload: ");
  for (int i = 0; i < length; i++)
  {
    Serial.print((char)payload[i]);
  }
  Serial.println();

  // Parse JSON payload
  StaticJsonDocument<200> doc;
  DeserializationError error = deserializeJson(doc, payload, length);

  // Check if there is an error in parsing JSON
  if (error)
  {
    Serial.print("deserializeJson() failed: ");
    Serial.println(error.c_str());
    return;
  }

  // Check if the received message is related to a lamp topic
  if (strcmp(topic, mqtt_topic_lamp1) == 0)
  {
    // Process the payload for lamp1 control
    processLampControl(doc, D0); // Specify the pin for lamp1
  }
  // Check if the received message is related to lamp2
  else if (strcmp(topic, mqtt_topic_lamp2) == 0)
  {
    // Process the payload for lamp2 control
    processLampControl(doc, D1); // Specify the pin for lamp2
  }
  // Check if the received message is related to lamp3
  else if (strcmp(topic, mqtt_topic_lamp3) == 0)
  {
    // Process the payload for lamp3 control
    processLampControl(doc, D2); // Specify the pin for lamp3
  }
  else if (strcmp(topic, mqtt_topic_ky005) == 0)
  {
    // Process the payload for ky005 control
    processKY005Control(doc, ky005Pin); // Specify the pin for ky005
  }
}
