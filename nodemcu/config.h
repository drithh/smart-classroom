// Config.h
#ifndef CONFIG_H
#define CONFIG_H

// Wi-Fi configuration
const char *ssid = "SW Home";
const char *password = "informatika";

// MQTT configuration
const char *mqtt_server = "192.168.1.8";
// Device-specific MQTT credentials
#ifdef DEVICE1
const char *mqtt_username = "pir";
const char *mqtt_password = "pir_password";
const char *mqtt_client_id = "pir";
#endif

#ifdef DEVICE2
const char *mqtt_username = "dht11_ldr";
const char *mqtt_password = "dht11_ldr_password";
const char *mqtt_client_id = "dht11_ldr";
#endif

#endif // CONFIG_H
