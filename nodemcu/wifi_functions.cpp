// WiFiFunctions.cpp
#include "wifi_functions.h"

extern const char *ssid;
extern const char *password;

void connectToWiFi()
{
  Serial.print("Connecting to WiFi");
  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED)
  {
    delay(250);
    Serial.print(".");
  }

  Serial.println();
  Serial.println("WiFi connected");
}
