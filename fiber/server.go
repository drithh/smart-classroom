package fiber

import (
	"fmt"
	"time"

	"github.com/drithh/smart-classroom/types"
	"github.com/gofiber/template/html/v2"
	"github.com/jmoiron/sqlx"

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

func SetupFiber(db *sqlx.DB) {
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
		setting := types.Setting{}

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

		// rende html
		// return c.JSON(setting)
		return c.Render("index", setting)
		// return json setting

	})

	app.Get("/charts/temperature", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})

	app.Get("/charts/humidity", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})

	app.Get("/charts/brightness", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
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
