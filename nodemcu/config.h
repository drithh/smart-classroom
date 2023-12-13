// Config.h
#ifndef CONFIG_H
#define CONFIG_H

// Wi-Fi configuration
const char *ssid = "SW Home";
const char *password = "informatika";

// MQTT configuration
const char *mqtt_server = "157.245.154.110";
// Device-specific MQTT credentials
#ifdef DEVICE1
const char *mqtt_username = "pir_dht11";
const char *mqtt_password = "pir_dht11_password";
const char *mqtt_client_id = "pir_dht11";
#endif

#ifdef DEVICE2
const char *mqtt_username = "ldr";
const char *mqtt_password = "ldr_password";
const char *mqtt_client_id = "ldr";
#endif

#endif // CONFIG_H
