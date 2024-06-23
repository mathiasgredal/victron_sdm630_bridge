package mqtt

import (
	"fmt"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	godbus "github.com/godbus/dbus"
	"sdm630-bridge.go/src/dbus"
)

/* Configuration */
type MqttConfig struct {
	Broker     string
	BrokerPort int
	Topic      string
	ClientId   string
	Username   string
	Password   string
}

/* Default configuration */
var DefaultConfig = MqttConfig{
	Broker:     "localhost",
	BrokerPort: 1883,
	Topic:      "stromzaehler/#",
	ClientId:   "sdm630-bridge",
	Username:   "user",
	Password:   "pass",
}

/* MQTT State */
type Topic struct {
	Path      godbus.ObjectPath
	Unit      string
	Regex     string
	Value     float64
	IsUpdated bool
}

/* Topics */
func NewState() []*Topic {
	return []*Topic{
		{dbus.AcL1Power, "W", ".*Power/L1$", 0, false},
		{dbus.AcL2Power, "W", ".*Power/L2$", 0, false},
		{dbus.AcL3Power, "W", ".*Power/L3$", 0, false},
		// {dbus.AcEnergyForward, "kWh", ".*/Import$"},
		// {dbus.AcEnergyReverse, "kWh", ".*/Export$"},
		{dbus.AcPower, "W", ".*/Power$", 0, false},
		// {dbus.AcL1Current, "A", ".*/Current/L1$"},
		// {dbus.AcL2Current, "A", ".*/Current/L2$"},
		// {dbus.AcL3Current, "A", ".*/Current/L3$"},
		// {dbus.AcL1Voltage, "V", ".*/Voltage/L1$"},
		// {dbus.AcL2Voltage, "V", ".*/Voltage/L2$"},
		// {dbus.AcL3Voltage, "V", ".*/Voltage/L3$"},
		// {dbus.AcL1EnergyForward, "kWh", ".*/Sum/L1$"},
		// {dbus.AcL2EnergyForward, "kWh", ".*/Sum/L2$"},
		// {dbus.AcL3EnergyForward, "kWh", ".*/Sum/L3$"},
		// {dbus.AcL1EnergyReverse, "kWh", ".*/Export/L1$"},
		// {dbus.AcL2EnergyReverse, "kWh", ".*/Export/L2$"},
		// {dbus.AcL3EnergyReverse, "kWh", ".*/Export/L3$"},
	}
}

/* Start the mqtt collector */
func StartCollector(config MqttConfig, topics []*Topic) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.Broker, config.BrokerPort))
	opts.SetClientID(config.ClientId)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)

	// Set the default publish handler
	opts.SetDefaultPublishHandler(func(_ mqtt.Client, m mqtt.Message) {
		onPublish(m, topics)
	})

	// Set the connection and connection lost handlers
	opts.OnConnect = onConnect
	opts.OnConnectionLost = onConnectionLost

	// Create a new mqtt client
	client := mqtt.NewClient(opts)

	// Connect to the broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe to the topic
	token := client.Subscribe(config.Topic, 1, nil)
	token.Wait()
}

/* Convert binary to float64 */
func bin2Float64(bin string) float64 {
	foostring := string(bin)
	result, err := strconv.ParseFloat(foostring, 64)
	if err != nil {
		panic(err)
	}
	return result
}

/* Search for string with regex */
func ContainString(searchstring string, str string) bool {
	obj, err := regexp.MatchString(searchstring, str)

	if err != nil {
		panic(err)
	}

	return obj
}

/* Called when connected to the broker */
func onConnect(client mqtt.Client) {
	log.Info("Connected to MQTT broker")
}

/* Called when connection is lost */
func onConnectionLost(client mqtt.Client, err error) {
	log.Info(fmt.Sprintf("Connect lost: %v", err))
	panic(err)
}

/* Called when a message is published */
func onPublish(msg mqtt.Message, topics []*Topic) {
	log.Debug(fmt.Sprintf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic()))

	for _, t := range topics {
		if !ContainString(t.Regex, msg.Topic()) {
			continue
		}

		t.Value = bin2Float64(string(msg.Payload()))
		t.IsUpdated = true
	}
}
