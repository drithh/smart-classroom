package fiber

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/drithh/smart-classroom/types"
	"github.com/gofiber/template/html/v2"
	"github.com/jmoiron/sqlx"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
)

func getSettingValue(ac types.Device, settingName string) string {
	for _, s := range ac.Setting {
		if s.SettingName == settingName {
			return s.SettingValue
		}
	}
	return "" // Return an empty string if the setting is not found
}

func setupQuery(outputColumns string, columns string, table string) string {

	query := fmt.Sprintf("SELECT %s FROM (SELECT id, %s FROM %s ORDER BY id DESC LIMIT 200) AS subquery ORDER BY id ASC", outputColumns, columns, table)
	fmt.Println(query)
	return query
}

func SetupFiber(db *sqlx.DB, mqtt pahomqtt.Client) {
	engine := html.New("./fiber/templates", ".html")

	// Reload the templates on each render, good for development
	engine.Reload(true) // Optional. Default: false

	// Debug will print each template that is parsed, good for debugging
	engine.Debug(true) // Optional. Default: false

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Define your Fiber routes and handlers here
	// For example, a simple route to handle HTTP requests
	app.Get("/", func(c *fiber.Ctx) error {
		type Chart struct {
			Label []string `db:"label"`
			Value []string `db:"value"`
		}
		type Page struct {
			types.Setting
			Motion      Chart
			Brightness  Chart
			Humidity    Chart
			Temperature Chart
		}
		setting := Page{}

		err := db.Get(&setting.Pir, "SELECT timestamp, presence FROM pir_sensor_data WHERE presence = true ORDER BY id DESC LIMIT 1")
		if err != nil {
			fmt.Println("Error getting pir sensor data: ", err)
		}

		err = db.Get(&setting.Ldr, "SELECT light_intensity FROM ldr_sensor_data ORDER BY id DESC LIMIT 1")
		if err != nil {
			fmt.Println("Error getting ldr sensor data: ", err)
		}

		err = db.Get(&setting.Dht11, "SELECT temperature, humidity FROM dht11_sensor_data ORDER BY id DESC LIMIT 1")
		if err != nil {
			fmt.Println("Error getting dht11 sensor data: ", err)
		}

		err = db.Get(&setting.Ac, "SELECT device_id, status FROM devices WHERE device_id = 'ac1'")
		if err != nil {
			fmt.Println("Error getting ac device data: ", err)
		}

		err = db.Select(&setting.Ac.Setting, "SELECT device_id, setting_name, setting_value FROM device_settings WHERE device_id = 'ac1'")
		if err != nil {
			fmt.Println("Error getting ac device settings data: ", err)
		}

		err = db.Get(&setting.Lamp1, "SELECT device_id, status FROM devices WHERE device_id = 'lamp1'")
		if err != nil {
			fmt.Println("Error getting lamp1 device data: ", err)
		}

		err = db.Select(&setting.Lamp1.Setting, "SELECT device_id, setting_name, setting_value FROM device_settings WHERE device_id = 'lamp1'")
		if err != nil {
			fmt.Println("Error getting lamp1 device settings data: ", err)
		}

		err = db.Get(&setting.Lamp2, "SELECT device_id, status FROM devices WHERE device_id = 'lamp2'")
		if err != nil {
			fmt.Println("Error getting lamp2 device data: ", err)
		}

		err = db.Select(&setting.Lamp2.Setting, "SELECT device_id, setting_name, setting_value FROM device_settings WHERE device_id = 'lamp2'")
		if err != nil {
			fmt.Println("Error getting lamp2 device settings data: ", err)
		}

		err = db.Get(&setting.Lamp3, "SELECT device_id, status FROM devices WHERE device_id = 'lamp3'")
		if err != nil {
			fmt.Println("Error getting lamp3 device data: ", err)
		}

		err = db.Select(&setting.Lamp3.Setting, "SELECT device_id, setting_name, setting_value FROM device_settings WHERE device_id = 'lamp3'")
		if err != nil {
			fmt.Println("Error getting lamp3 device settings data: ", err)
		}
		// setting.Lamp1.Setting[0].SettingName
		// convert ldr to percentage and invert
		setting.Ldr.Brightness = 100 - (setting.Ldr.Brightness * 100 / 1023)
		engine.AddFunc("GetSettingValue", func(device types.Device, settingName string) string {
			return getSettingValue(device, settingName)
		})
		engine.AddFunc("GetTimeSince", func(timestamp time.Time) string {
			minutes := time.Since(timestamp).Minutes()
			if minutes < 1 {
				return "Just now"
			} else if minutes < 60 {
				return fmt.Sprintf("%d minutes ago", int(minutes))
			} else if minutes < 1440 {
				return fmt.Sprintf("%d hours ago", int(minutes/60))
			} else {
				return fmt.Sprintf("%d days ago", int(minutes/1440))
			}
		})
		// chart
		err = db.Select(&setting.Motion.Value, setupQuery("presence", "presence", "pir_sensor_data"))
		if err != nil {
			fmt.Println("Error getting motion chart data: ", err)
		}

		err = db.Select(&setting.Motion.Label, setupQuery("label", "TO_CHAR(timestamp, 'HH24:MI') AS label", "pir_sensor_data"))
		if err != nil {
			fmt.Println("Error getting motion chart data: ", err)
		}

		err = db.Select(&setting.Brightness.Value, setupQuery("value", "100 - light_intensity / 10.24 as value", "ldr_sensor_data"))
		if err != nil {
			fmt.Println("Error getting brightness chart data: ", err)
		}

		err = db.Select(&setting.Brightness.Label, setupQuery("label", "TO_CHAR(timestamp, 'HH24:MI') AS label", "ldr_sensor_data"))
		if err != nil {
			fmt.Println("Error getting brightness chart data: ", err)
		}

		err = db.Select(&setting.Humidity.Value, setupQuery("humidity", "humidity", "dht11_sensor_data"))
		if err != nil {
			fmt.Println("Error getting humidity chart data: ", err)
		}

		err = db.Select(&setting.Humidity.Label, setupQuery("label", "TO_CHAR(timestamp, 'HH24:MI') AS label", "dht11_sensor_data"))
		if err != nil {
			fmt.Println("Error getting humidity chart data: ", err)
		}

		err = db.Select(&setting.Temperature.Value, setupQuery("temperature", "temperature", "dht11_sensor_data"))
		if err != nil {
			fmt.Println("Error getting temperature chart data: ", err)
		}

		err = db.Select(&setting.Temperature.Label, setupQuery("label", "TO_CHAR(timestamp, 'HH24:MI') AS label", "dht11_sensor_data"))
		if err != nil {
			fmt.Println("Error getting temperature chart data: ", err)
		}

		return c.Render("index", setting)
	})

	app.Post("/devices/lamp", func(c *fiber.Ctx) error {
		// get form data
		type Lamp struct {
			DeviceId   string `form:"device_id"`
			Status     bool   `form:"status"`
			Brightness int    `form:"brightness"`
		}

		form := new(Lamp)
		if err := c.BodyParser(form); err != nil {
			fmt.Println("Error parsing form data: ", err)
		}
		fmt.Println(form.DeviceId, form.Status, form.Brightness)
		_, err := db.Exec("UPDATE devices SET status = $1 WHERE device_id = $2", form.Status, form.DeviceId)
		if err != nil {
			fmt.Println("Error updating data into database: ", err)
		}
		strBrightness := fmt.Sprintf("%d", form.Brightness)

		_, err = db.Exec("UPDATE device_settings SET setting_value = $1 WHERE device_id = $2 AND setting_name = $3", strBrightness, form.DeviceId, "brightness")
		if err != nil {
			fmt.Println("Error updating data into database: ", err)
		}

		topic := fmt.Sprintf("classroom/actuator/%s", form.DeviceId)

		// make it json
		led := types.Led{
			Led:        form.Status,
			Brightness: form.Brightness,
		}

		// marshal to json
		ledJson, err := json.Marshal(led)

		if err != nil {
			fmt.Println("Error marshalling led data: ", err)
		}

		token := mqtt.Publish(topic, 1, false, ledJson)

		token.Wait()

		return nil

	})

	app.Post("/devices/ac", func(c *fiber.Ctx) error {
		type Ac struct {
			DeviceId    string `form:"device_id"`
			Status      bool   `form:"status"`
			Temperature int    `form:"temperature"`
			FanSpeed    int    `form:"fan_speed"`
			Swing       bool   `form:"swing"`
		}

		form := new(Ac)

		if err := c.BodyParser(form); err != nil {
			fmt.Println("Error parsing form data: ", err)
		}

		fmt.Println(form.DeviceId, form.Status, form.Temperature, form.FanSpeed, form.Swing)

		_, err := db.Exec("UPDATE devices SET status = $1 WHERE device_id = $2", form.Status, form.DeviceId)
		if err != nil {
			fmt.Println("Error updating data into database: ", err)
		}

		strTemperature := fmt.Sprintf("%d", form.Temperature)
		strFanSpeed := fmt.Sprintf("%d", form.FanSpeed)
		strSwing := "off"
		if form.Swing {
			strSwing = "on"
		}

		_, err = db.Exec("UPDATE device_settings SET setting_value = $1 WHERE device_id = $2 AND setting_name = $3", strTemperature, form.DeviceId, "temperature")
		if err != nil {
			fmt.Println("Error updating data into database: ", err)
		}

		_, err = db.Exec("UPDATE device_settings SET setting_value = $1 WHERE device_id = $2 AND setting_name = $3", strFanSpeed, form.DeviceId, "fan_speed")
		if err != nil {
			fmt.Println("Error updating data into database: ", err)
		}

		_, err = db.Exec("UPDATE device_settings SET setting_value = $1 WHERE device_id = $2 AND setting_name = $3", strSwing, form.DeviceId, "swing")
		if err != nil {
			fmt.Println("Error updating data into database: ", err)
		}

		topic := fmt.Sprintf("classroom/actuator/%s", form.DeviceId)

		// make it json
		ac := types.Ac{
			Status:      form.Status,
			Temperature: form.Temperature,
			FanSpeed:    form.FanSpeed,
			Swing:       form.Swing,
		}

		// marshal to json
		acJson, err := json.Marshal(ac)

		if err != nil {
			fmt.Println("Error marshalling led data: ", err)
		}

		token := mqtt.Publish(topic, 1, false, acJson)

		token.Wait()

		if err != nil {
			fmt.Println("Error marshalling led data: ", err)
		}

		return nil
	})

	// Start the Fiber web server
	go func() {
		err := app.Listen(":3000")
		if err != nil {
			fmt.Printf("Error starting server: %v", err)
		}
	}()

}
