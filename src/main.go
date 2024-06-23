package main

import (
	"fmt"
	"slices"
	"time"

	godbus "github.com/godbus/dbus"
	log "github.com/sirupsen/logrus"
	"sdm630-bridge.go/src/dbus"
	"sdm630-bridge.go/src/mqtt"
	"sdm630-bridge.go/src/vebus"
)

func main() {
	// Set the log level to debug
	log.SetLevel(log.DebugLevel)

	// Create a new state
	var state = dbus.New()

	// Start the mqtt collector
	var topics = mqtt.NewState()
	mqtt.StartCollector(mqtt.DefaultConfig, topics)

	// Create the vebus collector
	var affectedTopics = []godbus.ObjectPath{dbus.AcPower, dbus.AcL1Power}
	var measureDCPower = vebus.NewMeasureDCPowerFunc()
	var lastDcPower = measureDCPower()

	// Control loop
	for {
		time.Sleep(500 * time.Millisecond)

		var currentDcPower = measureDCPower()

		// Update the DC power
		for _, topic := range topics {
			for _, affectedTopic := range affectedTopics {
				if topic.Path != affectedTopic {
					continue
				}

				var correctedPower = (topic.Value - lastDcPower) + currentDcPower
				state.Update(topic.Path, dbus.MakeValueItem(fmt.Sprintf("%.2f", correctedPower)+topic.Unit, correctedPower))
			}
		}

		// Update the state
		for _, topic := range topics {
			// Check if the topic is updated
			if !topic.IsUpdated {
				continue
			}

			// Reset the flag
			topic.IsUpdated = false

			// Update the DC power if the AC power is updated
			if slices.Contains(affectedTopics, topic.Path) {
				lastDcPower = measureDCPower()
				continue
			}

			// Update the state
			state.Update(topic.Path, dbus.MakeValueItem(fmt.Sprintf("%.2f", topic.Value)+topic.Unit, topic.Value))
		}
	}
}
