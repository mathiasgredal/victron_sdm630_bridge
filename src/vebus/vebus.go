package vebus

import (
	log "github.com/sirupsen/logrus"

	godbus "github.com/godbus/dbus"
)

// NewMeasurePowerFunc creates a function which can be called to get the DC power value
func NewMeasureDCPowerFunc() func() float64 {
	// Create a connection to the system bus
	var conn, err = godbus.SystemBus()

	if err != nil {
		log.Fatal(err)
	}

	// Get the object
	var obj = conn.Object("com.victronenergy.system", "/Dc/Vebus/Power")

	// Return a function which can be called to get the value
	return func() float64 {
		// Print the value
		var res = obj.Call("GetValue", 0)

		if res.Err != nil {
			panic(res.Err)
		}

		var value float64
		err = res.Store(&value)

		if err != nil {
			panic(err)
		}

		return value
	}
}
