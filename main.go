package main

import (
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

func main() {
	var conn, err = dbus.SystemBus()

	if err != nil {
		log.Fatal(err)
	}

	// Get the object
	var obj = conn.Object("com.victronenergy.system", "/Dc/Vebus/Power")

	for {
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

		fmt.Println("Value: ", value)

		time.Sleep(1 * time.Second)
	}
}
