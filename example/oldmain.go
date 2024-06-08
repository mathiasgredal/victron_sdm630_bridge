// // package main

// // import (
// // 	"flag"
// // 	"fmt"
// // 	"time"

// // 	mqtt "github.com/eclipse/paho.mqtt.golang"
// // 	"github.com/godbus/dbus"
// // )

// // /* Configuration */
// // var (
// // 	broker     = "192.168.1.119"
// // 	brokerPort = 1883
// // 	topic      = "stromzaehler/#"
// // 	clientId   = "sdm630-bridge"
// // 	username   = "user"
// // 	password   = "pass"
// // )

// // /* Dbus Configuration */
// // type dbusValue struct {
// // 	raw       float64
// // 	pretty    string
// // 	signature dbus.Signature
// // }

// // func (d dbusValue) toVariant() dbus.Variant {
// // 	return dbus.MakeVariantWithSignature(d.raw, d.signature)
// // }

// // var dbusValues = map[string]dbusValue{}

// // func main() {
// // 	// Parse command line arguments
// // 	flag.StringVar(&broker, "broker", broker, "MQTT broker address")
// // 	flag.IntVar(&brokerPort, "port", brokerPort, "MQTT broker port")
// // 	flag.StringVar(&topic, "topic", topic, "MQTT topic prefix")
// // 	flag.StringVar(&clientId, "client-id", clientId, "MQTT client id")
// // 	flag.StringVar(&username, "username", username, "MQTT username")
// // 	flag.StringVar(&password, "password", password, "MQTT password")
// // 	flag.Parse()

// // 	// MQTT Subscription
// // 	opts := mqtt.NewClientOptions()
// // 	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, brokerPort))
// // 	opts.SetClientID(clientId)
// // 	opts.SetUsername(username)
// // 	opts.SetPassword(password)
// // 	opts.SetDefaultPublishHandler(messagePubHandler)
// // 	// opts.OnConnect = connectHandler
// // 	// opts.OnConnectionLost = connectLostHandler
// // 	client := mqtt.NewClient(opts)
// // 	if token := client.Connect(); token.Wait() && token.Error() != nil {
// // 		panic(token.Error())
// // 	}
// // 	sub(client)

// // 	time.Sleep(100 * time.Hour)

// // 	// // Connect to dbus
// // 	// var conn, err = dbus.SystemBus()
// // 	// defer conn.Close()

// // 	// if err != nil {
// // 	// 	panic(err)
// // 	// }

// // 	// // Get the object
// // 	// var obj = conn.Object("com.victronenergy.system", "/Dc/Vebus/Power")

// // 	// // Print the value
// // 	// for {

// // 	// 	var res = obj.Call("GetValue", 0)

// // 	// 	if res.Err != nil {
// // 	// 		panic(res.Err)
// // 	// 	}

// // 	// 	var value float64
// // 	// 	err = res.Store(&value)

// // 	// 	if err != nil {
// // 	// 		panic(err)
// // 	// 	}

// // 	// 	fmt.Println("Value: ", value)

// // 	// 	time.Sleep(1 * time.Millisecond)
// // 	// }
// // }

// // func sub(client mqtt.Client) {
// // 	topic := topic
// // 	token := client.Subscribe(topic, 1, nil)
// // 	token.Wait()
// // 	fmt.Println("Subscribed to topic: " + topic)
// // }

// // var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
// // 	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
// // }

// // // Every time we get a message from mqqt:
// // // 	- Update the dc power from last read
// // // 	- Update the ac power from last read

// // // Every 100 ms: Read the DC power
// // // - Subtract the current DC power from the last DC power to get the correction
// // // - Update the AC power with the correction

// // // package main

// // // import (
// // // 	"flag"
// // // 	"fmt"
// // // 	"os"
// // // 	"regexp"
// // // 	"strconv"
// // // 	"strings"
// // // 	"time"

// // // 	mqtt "github.com/eclipse/paho.mqtt.golang"
// // // 	"github.com/godbus/dbus/introspect"
// // // 	"github.com/godbus/dbus/v5"
// // // 	log "github.com/sirupsen/logrus"
// // // )

// // // /* Configuration */
// // // var (
// // // 	broker     = "192.168.1.119"
// // // 	brokerPort = 1883
// // // 	topic      = "stromzaehler/#"
// // // 	clientId   = "sdm630-bridge"
// // // 	username   = "user"
// // // 	password   = "pass"
// // // )

// // // var P1 float64 = 0.00
// // // var P2 float64 = 0.00
// // // var P3 float64 = 0.00

// // // var conn, err = dbus.SystemBus()

// // // const intro = `
// // // <node>
// // //    <interface name="com.victronenergy.BusItem">
// // //     <signal name="PropertiesChanged">
// // //       <arg type="a{sv}" name="properties" />
// // //     </signal>
// // //     <method name="SetValue">
// // //       <arg direction="in"  type="v" name="value" />
// // //       <arg direction="out" type="i" />
// // //     </method>
// // //     <method name="GetText">
// // //       <arg direction="out" type="s" />
// // //     </method>
// // //     <method name="GetValue">
// // //       <arg direction="out" type="v" />
// // //     </method>
// // //     </interface>` + introspect.IntrospectDataString + `</node> `

// // // type objectpath string

// // // var victronValues = map[int]map[objectpath]dbus.Variant{
// // // 	// 0: This will be used to store the VALUE variant
// // // 	0: map[objectpath]dbus.Variant{},
// // // 	// 1: This will be used to store the STRING variant
// // // 	1: map[objectpath]dbus.Variant{},
// // // }

// // // func (f objectpath) GetValue() (dbus.Variant, *dbus.Error) {
// // // 	log.Debug("GetValue() called for ", f)
// // // 	log.Debug("...returning ", victronValues[0][f])
// // // 	return victronValues[0][f], nil
// // // }
// // // func (f objectpath) GetText() (string, *dbus.Error) {
// // // 	log.Debug("GetText() called for ", f)
// // // 	log.Debug("...returning ", victronValues[1][f])
// // // 	// Why does this end up ""SOMEVAL"" ... trim it I guess
// // // 	return strings.Trim(victronValues[1][f].String(), "\""), nil
// // // }

// // // func init() {
// // // 	lvl, ok := os.LookupEnv("LOG_LEVEL")
// // // 	if !ok {
// // // 		lvl = "info"
// // // 	}

// // // 	ll, err := log.ParseLevel(lvl)
// // // 	if err != nil {
// // // 		ll = log.DebugLevel
// // // 	}

// // // 	log.SetLevel(ll)
// // // }

// // // func main() {
// // // 	// Parse command line arguments
// // // 	flag.StringVar(&broker, "broker", broker, "MQTT broker address")
// // // 	flag.IntVar(&brokerPort, "port", brokerPort, "MQTT broker port")
// // // 	flag.StringVar(&topic, "topic", topic, "MQTT topic prefix")
// // // 	flag.StringVar(&clientId, "client-id", clientId, "MQTT client id")
// // // 	flag.StringVar(&username, "username", username, "MQTT username")
// // // 	flag.StringVar(&password, "password", password, "MQTT password")
// // // 	flag.Parse()

// // // 	// Need to implement following paths:
// // // 	// https://github.com/victronenergy/venus/wiki/dbus#grid-meter
// // // 	// also in system.py
// // // 	victronValues[0]["/Connected"] = dbus.MakeVariant(1)
// // // 	victronValues[1]["/Connected"] = dbus.MakeVariant("1")

// // // 	victronValues[0]["/CustomName"] = dbus.MakeVariant("Grid meter")
// // // 	victronValues[1]["/CustomName"] = dbus.MakeVariant("Grid meter")

// // // 	victronValues[0]["/DeviceInstance"] = dbus.MakeVariant(30)
// // // 	victronValues[1]["/DeviceInstance"] = dbus.MakeVariant("30")

// // // 	// also in system.py
// // // 	victronValues[0]["/DeviceType"] = dbus.MakeVariant(71)
// // // 	victronValues[1]["/DeviceType"] = dbus.MakeVariant("71")

// // // 	victronValues[0]["/ErrorCode"] = dbus.MakeVariantWithSignature(0, dbus.SignatureOf(123))
// // // 	victronValues[1]["/ErrorCode"] = dbus.MakeVariant("0")

// // // 	victronValues[0]["/FirmwareVersion"] = dbus.MakeVariant(2)
// // // 	victronValues[1]["/FirmwareVersion"] = dbus.MakeVariant("2")

// // // 	// also in system.py
// // // 	victronValues[0]["/Mgmt/Connection"] = dbus.MakeVariant("/dev/ttyUSB0")
// // // 	victronValues[1]["/Mgmt/Connection"] = dbus.MakeVariant("/dev/ttyUSB0")

// // // 	victronValues[0]["/Mgmt/ProcessName"] = dbus.MakeVariant("/opt/color-control/dbus-cgwacs/dbus-cgwacs")
// // // 	victronValues[1]["/Mgmt/ProcessName"] = dbus.MakeVariant("/opt/color-control/dbus-cgwacs/dbus-cgwacs")

// // // 	victronValues[0]["/Mgmt/ProcessVersion"] = dbus.MakeVariant("1.8.0")
// // // 	victronValues[1]["/Mgmt/ProcessVersion"] = dbus.MakeVariant("1.8.0")

// // // 	victronValues[0]["/Position"] = dbus.MakeVariantWithSignature(0, dbus.SignatureOf(123))
// // // 	victronValues[1]["/Position"] = dbus.MakeVariant("0")

// // // 	// also in system.py
// // // 	victronValues[0]["/ProductId"] = dbus.MakeVariant(45058)
// // // 	victronValues[1]["/ProductId"] = dbus.MakeVariant("45058")

// // // 	// also in system.py
// // // 	victronValues[0]["/ProductName"] = dbus.MakeVariant("Grid meter")
// // // 	victronValues[1]["/ProductName"] = dbus.MakeVariant("Grid meter")

// // // 	victronValues[0]["/Serial"] = dbus.MakeVariant("BP98305081235")
// // // 	victronValues[1]["/Serial"] = dbus.MakeVariant("BP98305081235")

// // // 	// Provide some initial values... note that the values must be a valid formt otherwise dbus_systemcalc.py exits like this:
// // // 	// @400000005ecc11bf3782b374   File "/opt/victronenergy/dbus-systemcalc-py/dbus_systemcalc.py", line 386, in _handletimertick
// // // 	// @400000005ecc11bf37aa251c     self._updatevalues()
// // // 	// @400000005ecc11bf380e74cc   File "/opt/victronenergy/dbus-systemcalc-py/dbus_systemcalc.py", line 678, in _updatevalues
// // // 	// @400000005ecc11bf383ab4ec     c = _safeadd(c, p, pvpower)
// // // 	// @400000005ecc11bf386c9674   File "/opt/victronenergy/dbus-systemcalc-py/sc_utils.py", line 13, in safeadd
// // // 	// @400000005ecc11bf387b28ec     return sum(values) if values else None
// // // 	// @400000005ecc11bf38b2bb7c TypeError: unsupported operand type(s) for +: 'int' and 'unicode'
// // // 	//
// // // 	victronValues[0]["/Ac/L1/Power"] = dbus.MakeVariant(0.0)
// // // 	victronValues[1]["/Ac/L1/Power"] = dbus.MakeVariant("0 W")
// // // 	victronValues[0]["/Ac/L2/Power"] = dbus.MakeVariant(0.0)
// // // 	victronValues[1]["/Ac/L2/Power"] = dbus.MakeVariant("0 W")
// // // 	victronValues[0]["/Ac/L3/Power"] = dbus.MakeVariant(0.0)
// // // 	victronValues[1]["/Ac/L3/Power"] = dbus.MakeVariant("0 W")

// // // 	victronValues[0]["/Ac/L1/Voltage"] = dbus.MakeVariant(230)
// // // 	victronValues[1]["/Ac/L1/Voltage"] = dbus.MakeVariant("230 V")
// // // 	victronValues[0]["/Ac/L2/Voltage"] = dbus.MakeVariant(230)
// // // 	victronValues[1]["/Ac/L2/Voltage"] = dbus.MakeVariant("230 V")
// // // 	victronValues[0]["/Ac/L3/Voltage"] = dbus.MakeVariant(230)
// // // 	victronValues[1]["/Ac/L3/Voltage"] = dbus.MakeVariant("230 V")

// // // 	victronValues[0]["/Ac/L1/Current"] = dbus.MakeVariant(0.0)
// // // 	victronValues[1]["/Ac/L1/Current"] = dbus.MakeVariant("0 A")
// // // 	victronValues[0]["/Ac/L2/Current"] = dbus.MakeVariant(0.0)
// // // 	victronValues[1]["/Ac/L2/Current"] = dbus.MakeVariant("0 A")
// // // 	victronValues[0]["/Ac/L3/Current"] = dbus.MakeVariant(0.0)
// // // 	victronValues[1]["/Ac/L3/Current"] = dbus.MakeVariant("0 A")

// // // 	victronValues[0]["/Ac/L1/Energy/Forward"] = dbus.MakeVariant(0.0)
// // // 	victronValues[1]["/Ac/L1/Energy/Forward"] = dbus.MakeVariant("0 kWh")
// // // 	victronValues[0]["/Ac/L2/Energy/Forward"] = dbus.MakeVariant(0.0)
// // // 	victronValues[1]["/Ac/L2/Energy/Forward"] = dbus.MakeVariant("0 kWh")
// // // 	victronValues[0]["/Ac/L3/Energy/Forward"] = dbus.MakeVariant(0.0)
// // // 	victronValues[1]["/Ac/L3/Energy/Forward"] = dbus.MakeVariant("0 kWh")

// // // 	victronValues[0]["/Ac/L1/Energy/Reverse"] = dbus.MakeVariant(0.0)
// // // 	victronValues[1]["/Ac/L1/Energy/Reverse"] = dbus.MakeVariant("0 kWh")
// // // 	victronValues[0]["/Ac/L2/Energy/Reverse"] = dbus.MakeVariant(0.0)
// // // 	victronValues[1]["/Ac/L2/Energy/Reverse"] = dbus.MakeVariant("0 kWh")
// // // 	victronValues[0]["/Ac/L3/Energy/Reverse"] = dbus.MakeVariant(0.0)
// // // 	victronValues[1]["/Ac/L3/Energy/Reverse"] = dbus.MakeVariant("0 kWh")

// // // 	basicPaths := []dbus.ObjectPath{
// // // 		"/Connected",
// // // 		"/CustomName",
// // // 		"/DeviceInstance",
// // // 		"/DeviceType",
// // // 		"/ErrorCode",
// // // 		"/FirmwareVersion",
// // // 		"/Mgmt/Connection",
// // // 		"/Mgmt/ProcessName",
// // // 		"/Mgmt/ProcessVersion",
// // // 		"/Position",
// // // 		"/ProductId",
// // // 		"/ProductName",
// // // 		"/Serial",
// // // 	}

// // // 	updatingPaths := []dbus.ObjectPath{
// // // 		"/Ac/L1/Power",
// // // 		"/Ac/L2/Power",
// // // 		"/Ac/L3/Power",
// // // 		"/Ac/L1/Voltage",
// // // 		"/Ac/L2/Voltage",
// // // 		"/Ac/L3/Voltage",
// // // 		"/Ac/L1/Current",
// // // 		"/Ac/L2/Current",
// // // 		"/Ac/L3/Current",
// // // 		"/Ac/L1/Energy/Forward",
// // // 		"/Ac/L2/Energy/Forward",
// // // 		"/Ac/L3/Energy/Forward",
// // // 		"/Ac/L1/Energy/Reverse",
// // // 		"/Ac/L2/Energy/Reverse",
// // // 		"/Ac/L3/Energy/Reverse",
// // // 	}

// // // 	defer conn.Close()

// // // 	// Some of the victron stuff requires it be called grid.cgwacs... using the only known valid value (from the simulator)
// // // 	// This can _probably_ be changed as long as it matches com.victronenergy.grid.cgwacs_*
// // // 	reply, err := conn.RequestName("com.victronenergy.grid.cgwacs_ttyUSB0_di30_mb1",
// // // 		dbus.NameFlagDoNotQueue)
// // // 	if err != nil {
// // // 		log.Panic("Something went horribly wrong in the dbus connection")
// // // 		panic(err)
// // // 	}

// // // 	if reply != dbus.RequestNameReplyPrimaryOwner {
// // // 		log.Panic("name cgwacs_ttyUSB0_di30_mb1 already taken on dbus.")
// // // 		os.Exit(1)
// // // 	}

// // // 	for i, s := range basicPaths {
// // // 		log.Debug("Registering dbus basic path #", i, ": ", s)
// // // 		conn.Export(objectpath(s), s, "com.victronenergy.BusItem")
// // // 		conn.Export(introspect.Introspectable(intro), s, "org.freedesktop.DBus.Introspectable")
// // // 	}

// // // 	for i, s := range updatingPaths {
// // // 		log.Debug("Registering dbus update path #", i, ": ", s)
// // // 		conn.Export(objectpath(s), s, "com.victronenergy.BusItem")
// // // 		conn.Export(introspect.Introspectable(intro), s, "org.freedesktop.DBus.Introspectable")
// // // 	}

// // // 	log.Info("Successfully connected to dbus and registered as a meter... Commencing reading of the SDM630 meter")

// // // 	// MQTT Subscripte
// // // 	opts := mqtt.NewClientOptions()
// // // 	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, brokerPort))
// // // 	opts.SetClientID(clientId)
// // // 	opts.SetUsername(username)
// // // 	opts.SetPassword(password)
// // // 	opts.SetDefaultPublishHandler(messagePubHandler)
// // // 	opts.OnConnect = connectHandler
// // // 	opts.OnConnectionLost = connectLostHandler
// // // 	client := mqtt.NewClient(opts)
// // // 	if token := client.Connect(); token.Wait() && token.Error() != nil {
// // // 		panic(token.Error())
// // // 	}
// // // 	sub(client)
// // // 	// Infinite loop
// // // 	for true {
// // // 		// Every 100ms, we want to read the Dc/0/Power
// // // 		// fmt.Println("Infinite Loop entered")
// // // 		time.Sleep(time.Second)
// // // 	}

// // // 	// This is a forever loop^^
// // // 	panic("Error: We terminated.... how did we ever get here?")
// // // }

// // // /* MQTT Subscribe Function */
// // // func sub(client mqtt.Client) {
// // // 	topic := topic
// // // 	token := client.Subscribe(topic, 1, nil)
// // // 	token.Wait()
// // // 	log.Info("Subscribed to topic: " + topic)
// // // }

// // // /* MQTT Publish Function */
// // // func publish(client mqtt.Client) {
// // // 	num := 10
// // // 	for i := 0; i < num; i++ {
// // // 		text := fmt.Sprintf("Message %d", i)
// // // 		token := client.Publish("topic/test", 0, false, text)
// // // 		token.Wait()
// // // 		time.Sleep(time.Second)
// // // 	}
// // // }

// // // /* Write dbus Values to Victron handler */
// // // func updateVariant(value float64, unit string, path string) {
// // // 	emit := make(map[string]dbus.Variant)
// // // 	emit["Text"] = dbus.MakeVariant(fmt.Sprintf("%.2f", value) + unit)
// // // 	emit["Value"] = dbus.MakeVariant(float64(value))
// // // 	victronValues[0][objectpath(path)] = emit["Value"]
// // // 	victronValues[1][objectpath(path)] = emit["Text"]
// // // 	conn.Emit(dbus.ObjectPath(path), "com.victronenergy.BusItem.PropertiesChanged", emit)
// // // }

// // // /* Convert binary to float64 */
// // // func bin2Float64(bin string) float64 {
// // // 	foostring := string(bin)
// // // 	result, err := strconv.ParseFloat(foostring, 64)
// // // 	if err != nil {
// // // 		panic(err)
// // // 	}
// // // 	return result
// // // }

// // // /* Called if connection is established */
// // // var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
// // // 	log.Info(fmt.Sprintf("Connected to broker %s:%d", broker, brokerPort))
// // // }

// // // /* Called if connection is lost  */
// // // var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
// // // 	log.Info(fmt.Sprintf("Connect lost: %v", err))
// // // 	os.Exit(1)
// // // }

// // // /* Search for string with regex */
// // // func ContainString(searchstring string, str string) bool {
// // // 	var obj bool

// // // 	obj, err = regexp.MatchString(searchstring, str)

// // // 	if err != nil {
// // // 		panic(err)
// // // 	}

// // // 	return obj
// // // }

// // // /* MQTT Subscribe Handler */
// // // var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

// // // 	log.Debug(fmt.Sprintf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic()))

// // // 	// Power L1
// // // 	if ContainString(".*Power/L1$", msg.Topic()) {
// // // 		P1 = bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(P1), "W", "/Ac/L1/Power")
// // // 	}

// // // 	// Power L2
// // // 	if ContainString(".*Power/L2$", msg.Topic()) {
// // // 		P2 = bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(P2), "W", "/Ac/L2/Power")
// // // 	}

// // // 	// Power L3
// // // 	if ContainString(".*Power/L3$", msg.Topic()) {
// // // 		P3 = bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(P3), "W", "/Ac/L3/Power")
// // // 	}

// // // 	// /Ac/Energy/Forward     <- kWh  - bought energy (total of all phases)
// // // 	if ContainString(".*/Import$", msg.Topic()) {
// // // 		IP := bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(IP), "kWh", "/Ac/Energy/Forward")
// // // 		log.Debug(fmt.Sprintf("Import Power: %.3f kWh", IP))
// // // 	}

// // // 	// /Ac/Energy/Reverse     <- kWh  - sold energy (total of all phases)
// // // 	if ContainString(".*/Export$", msg.Topic()) {
// // // 		EP := bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(EP), "kWh", "/Ac/Energy/Reverse")
// // // 		log.Debug(fmt.Sprintf("Export Power: %.3f kWh", EP))
// // // 	}

// // // 	// /Ac/Power              <- W    - total of all phases, real power
// // // 	if ContainString(".*/Power$", msg.Topic()) {
// // // 		TP := bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(TP), "W", "/Ac/Power")
// // // 		log.Debug(fmt.Sprintf("Total Power: %.3f W", TP))
// // // 	}

// // // 	// /Ac/L1/Current         <- A AC
// // // 	if ContainString(".*/Current/L1$", msg.Topic()) {
// // // 		CL1 := bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(CL1), "A", "/Ac/L1/Current")
// // // 		log.Debug(fmt.Sprintf("Current L1: %.3f A", CL1))
// // // 	}

// // // 	// /Ac/L2/Current         <- A AC
// // // 	if ContainString(".*/Current/L2$", msg.Topic()) {
// // // 		CL2 := bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(CL2), "A", "/Ac/L2/Current")
// // // 		log.Debug(fmt.Sprintf("Current L2: %.3f A", CL2))
// // // 	}

// // // 	// /Ac/L3/Current         <- A AC
// // // 	if ContainString(".*/Current/L3$", msg.Topic()) {
// // // 		CL3 := bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(CL3), "A", "/Ac/L3/Current")
// // // 		log.Debug(fmt.Sprintf("Current L3: %.3f A", CL3))
// // // 	}

// // // 	// /Ac/L1/Voltage <- V AC
// // // 	if ContainString(".*/Voltage/L1$", msg.Topic()) {
// // // 		VL1 := bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(VL1), "V", "/Ac/L1/Voltage")
// // // 		log.Debug(fmt.Sprintf("Voltage L1: %.3f V", VL1))
// // // 	}

// // // 	// /Ac/L2/Voltage <- V AC
// // // 	if ContainString(".*/Voltage/L2$", msg.Topic()) {
// // // 		VL2 := bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(VL2), "V", "/Ac/L2/Voltage")
// // // 		log.Debug(fmt.Sprintf("Voltage L2: %.3f V", VL2))
// // // 	}

// // // 	// /Ac/L3/Voltage <- V AC
// // // 	if ContainString(".*/Voltage/L3$", msg.Topic()) {
// // // 		VL3 := bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(VL3), "V", "/Ac/L3/Voltage")
// // // 		log.Debug(fmt.Sprintf("Voltage L3: %.3f V", VL3))
// // // 	}

// // // 	// /Ac/L1/Energy/Forward  <- kWh  - bought
// // // 	if ContainString(".*/Sum/L1$", msg.Topic()) {
// // // 		SL1 := bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(SL1), "kWh", "/Ac/L1/Energy/Forward")
// // // 		log.Debug(fmt.Sprintf("Energy Forward L1: %.3f kWh", SL1))
// // // 	}

// // // 	// /Ac/L2/Energy/Forward  <- kWh  - bought
// // // 	if ContainString(".*/Sum/L2$", msg.Topic()) {
// // // 		SL2 := bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(SL2), "kWh", "/Ac/L2/Energy/Forward")
// // // 		log.Debug(fmt.Sprintf("Energy Forward L2: %.3f kWh", SL2))
// // // 	}
// // // 	// /Ac/L3/Energy/Forward  <- kWh  - bought
// // // 	if ContainString(".*/Sum/L3$", msg.Topic()) {
// // // 		SL3 := bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(SL3), "kWh", "/Ac/L3/Energy/Forward")
// // // 		log.Debug(fmt.Sprintf("Energy Forward L3: %.3f kWh", SL3))
// // // 	}

// // // 	// /Ac/L1/Energy/Reverse  <- kWh  - bought
// // // 	if ContainString(".*/Export/L1$", msg.Topic()) {
// // // 		EL1 := bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(EL1), "kWh", "/Ac/L1/Energy/Reverse")
// // // 		log.Debug(fmt.Sprintf("Energy Reverse L1: %.3f kWh", EL1))
// // // 	}

// // // 	// /Ac/L2/Energy/Reverse  <- kWh  - bought
// // // 	if ContainString(".*/Export/L2$", msg.Topic()) {
// // // 		EL2 := bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(EL2), "kWh", "/Ac/L2/Energy/Reverse")
// // // 		log.Debug(fmt.Sprintf("Energy Reverse L2: %.3f kWh", EL2))
// // // 	}
// // // 	// /Ac/L3/Energy/Reverse  <- kWh  - bought
// // // 	if ContainString(".*/Export/L3$", msg.Topic()) {
// // // 		EL3 := bin2Float64(string(msg.Payload()))
// // // 		updateVariant(float64(EL3), "kWh", "/Ac/L3/Energy/Reverse")
// // // 		log.Debug(fmt.Sprintf("Energy Reverse L3: %.3f kWh", EL3))
// // // 	}
// // // }

// package main

// import (
// 	"flag"
// 	"fmt"
// 	"os"
// 	"regexp"
// 	"strconv"
// 	"time"

// 	mqtt "github.com/eclipse/paho.mqtt.golang"
// 	"github.com/godbus/dbus/introspect"
// 	"github.com/godbus/dbus/v5"
// 	log "github.com/sirupsen/logrus"
// )

// /* Configuration */
// var (
// 	broker     = "192.168.1.119"
// 	brokerPort = 1883
// 	topic      = "stromzaehler/#"
// 	clientId   = "sdm630-bridge"
// 	username   = "user"
// 	password   = "pass"
// )

// var P1 float64 = 0.00
// var P2 float64 = 0.00
// var P3 float64 = 0.00
// var conn, err = dbus.SystemBus()

// const intro = `
// <node>
//    <interface name="com.victronenergy.BusItem">
//     <signal name="PropertiesChanged">
//       <arg type="a{sv}" name="properties" />
//     </signal>
//     <method name="SetValue">
//       <arg direction="in"  type="v" name="value" />
//       <arg direction="out" type="i" />
//     </method>
//     <method name="GetText">
//       <arg direction="out" type="s" />
//     </method>
//     <method name="GetValue">
//       <arg direction="out" type="v" />
//     </method>
//     </interface>` + introspect.IntrospectDataString + `</node> `

// type objectpath string

// type dbusValue struct {
// 	value     *float64
// 	text      string
// 	signature *dbus.Signature
// }

// func MakeDbusNumericValue(text string, value float64, signature *dbus.Signature) dbusValue {
// 	return dbusValue{
// 		value:     &value,
// 		text:      text,
// 		signature: signature,
// 	}
// }

// func MakeDbusTextValue(text string, signature *dbus.Signature) dbusValue {
// 	return dbusValue{
// 		value:     nil,
// 		text:      text,
// 		signature: signature,
// 	}
// }

// var victronValues = map[objectpath]dbusValue{}

// func (f objectpath) GetValue() (dbus.Variant, *dbus.Error) {
// 	log.Debug("GetValue() called for ", f)
// 	log.Debug("...returning ", victronValues[f].value)
// 	if victronValues[f].signature != nil {
// 		return dbus.MakeVariantWithSignature(victronValues[f].value, *victronValues[f].signature), nil
// 	}
// 	return dbus.MakeVariant(victronValues[f].value), nil
// }

// func (f objectpath) GetText() (string, *dbus.Error) {
// 	log.Debug("GetText() called for ", f)
// 	log.Debug("...returning ", victronValues[f].text)
// 	return victronValues[f].text, nil
// }

// func init() {
// 	lvl, ok := os.LookupEnv("LOG_LEVEL")
// 	if !ok {
// 		lvl = "info"
// 	}

// 	ll, err := log.ParseLevel(lvl)
// 	if err != nil {
// 		ll = log.DebugLevel
// 	}

// 	log.SetLevel(ll)
// }

// func main() {
// 	// Parse command line arguments
// 	flag.StringVar(&broker, "broker", broker, "MQTT broker address")
// 	flag.IntVar(&brokerPort, "port", brokerPort, "MQTT broker port")
// 	flag.StringVar(&topic, "topic", topic, "MQTT topic prefix")
// 	flag.StringVar(&clientId, "client-id", clientId, "MQTT client id")
// 	flag.StringVar(&username, "username", username, "MQTT username")
// 	flag.StringVar(&password, "password", password, "MQTT password")
// 	flag.Parse()

// 	// Need to implement following paths:
// 	// https://github.com/victronenergy/venus/wiki/dbus#grid-meter
// 	// also in system.py
// 	victronValues["/Connected"] = MakeDbusNumericValue("1", 1, nil)
// 	victronValues["/CustomName"] = MakeDbusTextValue("Grid meter", nil)
// 	victronValues["/DeviceInstance"] = MakeDbusNumericValue("30", 30, nil)

// 	// also in system.py
// 	victronValues["/DeviceType"] = MakeDbusNumericValue("71", 71, nil)
// 	var signature = dbus.SignatureOf(123)
// 	victronValues["/ErrorCode"] = MakeDbusNumericValue("0", 0, &signature)
// 	victronValues["/FirmwareVersion"] = MakeDbusNumericValue("2", 2, nil)

// 	// also in system.py
// 	victronValues["/Mgmt/Connection"] = MakeDbusTextValue("/dev/ttyUSB0", nil)
// 	victronValues["/Mgmt/ProcessName"] = MakeDbusTextValue("/opt/color-control/dbus-cgwacs/dbus-cgwacs", nil)
// 	victronValues["/Mgmt/ProcessVersion"] = MakeDbusTextValue("1.8.0", nil)
// 	victronValues["/Position"] = MakeDbusNumericValue("0", 0, &signature)

// 	// also in system.py
// 	victronValues["/ProductId"] = MakeDbusNumericValue("45058", 45058, nil)
// 	victronValues["/ProductName"] = MakeDbusTextValue("Grid meter", nil)
// 	victronValues["/Serial"] = MakeDbusTextValue("BP98305081235", nil)

// 	// Provide some initial values... note that the values must be a valid
// 	victronValues["/Ac/L1/Power"] = MakeDbusNumericValue("0 W", 0.0, nil)
// 	victronValues["/Ac/L2/Power"] = MakeDbusNumericValue("0 W", 0.0, nil)
// 	victronValues["/Ac/L3/Power"] = MakeDbusNumericValue("0 W", 0.0, nil)

// 	victronValues["/Ac/L1/Voltage"] = MakeDbusNumericValue("230 V", 230, nil)
// 	victronValues["/Ac/L2/Voltage"] = MakeDbusNumericValue("230 V", 230, nil)
// 	victronValues["/Ac/L3/Voltage"] = MakeDbusNumericValue("230 V", 230, nil)

// 	victronValues["/Ac/L1/Current"] = MakeDbusNumericValue("0 A", 0.0, nil)
// 	victronValues["/Ac/L2/Current"] = MakeDbusNumericValue("0 A", 0.0, nil)
// 	victronValues["/Ac/L3/Current"] = MakeDbusNumericValue("0 A", 0.0, nil)

// 	victronValues["/Ac/L1/Energy/Forward"] = MakeDbusNumericValue("0 kWh", 0.0, nil)
// 	victronValues["/Ac/L2/Energy/Forward"] = MakeDbusNumericValue("0 kWh", 0.0, nil)
// 	victronValues["/Ac/L3/Energy/Forward"] = MakeDbusNumericValue("0 kWh", 0.0, nil)

// 	victronValues["/Ac/L1/Energy/Reverse"] = MakeDbusNumericValue("0 kWh", 0.0, nil)
// 	victronValues["/Ac/L2/Energy/Reverse"] = MakeDbusNumericValue("0 kWh", 0.0, nil)
// 	victronValues["/Ac/L3/Energy/Reverse"] = MakeDbusNumericValue("0 kWh", 0.0, nil)

// 	defer conn.Close()

// 	// Some of the victron stuff requires it be called grid.cgwacs... using the only known valid value (from the simulator)
// 	// This can _probably_ be changed as long as it matches com.victronenergy.grid.cgwacs_*
// 	reply, err := conn.RequestName("com.victronenergy.grid.cgwacs_ttyUSB0_di30_mb1",
// 		dbus.NameFlagDoNotQueue)
// 	if err != nil {
// 		log.Panic("Something went horribly wrong in the dbus connection")
// 		panic(err)
// 	}

// 	if reply != dbus.RequestNameReplyPrimaryOwner {
// 		log.Panic("name cgwacs_ttyUSB0_di30_mb1 already taken on dbus.")
// 		os.Exit(1)
// 	}

// 	for path := range victronValues {
// 		conn.Export(path, dbus.ObjectPath(path), "com.victronenergy.BusItem")
// 		conn.Export(introspect.Introspectable(intro), dbus.ObjectPath(path), "org.freedesktop.DBus.Introspectable")
// 	}

// 	log.Info("Successfully connected to dbus and registered as a meter... Commencing reading of the SDM630 meter")

// 	// MQTT Subscripte
// 	opts := mqtt.NewClientOptions()
// 	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, brokerPort))
// 	opts.SetClientID(clientId)
// 	opts.SetUsername(username)
// 	opts.SetPassword(password)
// 	opts.SetDefaultPublishHandler(messagePubHandler)
// 	opts.OnConnect = connectHandler
// 	opts.OnConnectionLost = connectLostHandler
// 	client := mqtt.NewClient(opts)
// 	if token := client.Connect(); token.Wait() && token.Error() != nil {
// 		panic(token.Error())
// 	}
// 	sub(client)
// 	// Infinite loop
// 	for true {
// 		// fmt.Println("Infinite Loop entered")
// 		time.Sleep(time.Second)
// 	}

// 	// This is a forever loop^^
// 	panic("Error: We terminated.... how did we ever get here?")
// }

// /* MQTT Subscribe Function */
// func sub(client mqtt.Client) {
// 	topic := topic
// 	token := client.Subscribe(topic, 1, nil)
// 	token.Wait()
// 	log.Info("Subscribed to topic: " + topic)
// }

// /* MQTT Publish Function */
// func publish(client mqtt.Client) {
// 	num := 10
// 	for i := 0; i < num; i++ {
// 		text := fmt.Sprintf("Message %d", i)
// 		token := client.Publish("topic/test", 0, false, text)
// 		token.Wait()
// 		time.Sleep(time.Second)
// 	}
// }

// /* Write dbus Values to Victron handler */
// func updateVariant(value float64, unit string, path string) {
// 	emit := make(map[string]dbus.Variant)
// 	emit["Text"] = dbus.MakeVariant(fmt.Sprintf("%.2f", value) + unit)
// 	emit["Value"] = dbus.MakeVariant(float64(value))
// 	var val = victronValues[objectpath(path)]
// 	val.value = &value
// 	val.text = fmt.Sprintf("%.2f", value) + unit
// 	conn.Emit(dbus.ObjectPath(path), "com.victronenergy.BusItem.PropertiesChanged", emit)
// }

// /* Convert binary to float64 */
// func bin2Float64(bin string) float64 {
// 	foostring := string(bin)
// 	result, err := strconv.ParseFloat(foostring, 64)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return result
// }

// /* Called if connection is established */
// var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
// 	log.Info(fmt.Sprintf("Connected to broker %s:%d", broker, brokerPort))
// }

// /* Called if connection is lost  */
// var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
// 	log.Info(fmt.Sprintf("Connect lost: %v", err))
// 	os.Exit(1)
// }

// /* Search for string with regex */
// func ContainString(searchstring string, str string) bool {
// 	var obj bool

// 	obj, err = regexp.MatchString(searchstring, str)

// 	if err != nil {
// 		panic(err)
// 	}

// 	return obj
// }

// /* MQTT Subscribe Handler */
// var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

// 	log.Debug(fmt.Sprintf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic()))

// 	// Power L1
// 	if ContainString(".*Power/L1$", msg.Topic()) {
// 		P1 = bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(P1), "W", "/Ac/L1/Power")
// 	}

// 	// Power L2
// 	if ContainString(".*Power/L2$", msg.Topic()) {
// 		P2 = bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(P2), "W", "/Ac/L2/Power")
// 	}

// 	// Power L3
// 	if ContainString(".*Power/L3$", msg.Topic()) {
// 		P3 = bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(P3), "W", "/Ac/L3/Power")
// 	}

// 	// /Ac/Energy/Forward     <- kWh  - bought energy (total of all phases)
// 	if ContainString(".*/Import$", msg.Topic()) {
// 		IP := bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(IP), "kWh", "/Ac/Energy/Forward")
// 		log.Debug(fmt.Sprintf("Import Power: %.3f kWh", IP))
// 	}

// 	// /Ac/Energy/Reverse     <- kWh  - sold energy (total of all phases)
// 	if ContainString(".*/Export$", msg.Topic()) {
// 		EP := bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(EP), "kWh", "/Ac/Energy/Reverse")
// 		log.Debug(fmt.Sprintf("Export Power: %.3f kWh", EP))
// 	}

// 	// /Ac/Power              <- W    - total of all phases, real power
// 	if ContainString(".*/Power$", msg.Topic()) {
// 		TP := bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(TP), "W", "/Ac/Power")
// 		log.Debug(fmt.Sprintf("Total Power: %.3f W", TP))
// 	}

// 	// /Ac/L1/Current         <- A AC
// 	if ContainString(".*/Current/L1$", msg.Topic()) {
// 		CL1 := bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(CL1), "A", "/Ac/L1/Current")
// 		log.Debug(fmt.Sprintf("Current L1: %.3f A", CL1))
// 	}

// 	// /Ac/L2/Current         <- A AC
// 	if ContainString(".*/Current/L2$", msg.Topic()) {
// 		CL2 := bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(CL2), "A", "/Ac/L2/Current")
// 		log.Debug(fmt.Sprintf("Current L2: %.3f A", CL2))
// 	}

// 	// /Ac/L3/Current         <- A AC
// 	if ContainString(".*/Current/L3$", msg.Topic()) {
// 		CL3 := bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(CL3), "A", "/Ac/L3/Current")
// 		log.Debug(fmt.Sprintf("Current L3: %.3f A", CL3))
// 	}

// 	// /Ac/L1/Voltage <- V AC
// 	if ContainString(".*/Voltage/L1$", msg.Topic()) {
// 		VL1 := bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(VL1), "V", "/Ac/L1/Voltage")
// 		log.Debug(fmt.Sprintf("Voltage L1: %.3f V", VL1))
// 	}

// 	// /Ac/L2/Voltage <- V AC
// 	if ContainString(".*/Voltage/L2$", msg.Topic()) {
// 		VL2 := bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(VL2), "V", "/Ac/L2/Voltage")
// 		log.Debug(fmt.Sprintf("Voltage L2: %.3f V", VL2))
// 	}

// 	// /Ac/L3/Voltage <- V AC
// 	if ContainString(".*/Voltage/L3$", msg.Topic()) {
// 		VL3 := bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(VL3), "V", "/Ac/L3/Voltage")
// 		log.Debug(fmt.Sprintf("Voltage L3: %.3f V", VL3))
// 	}

// 	// /Ac/L1/Energy/Forward  <- kWh  - bought
// 	if ContainString(".*/Sum/L1$", msg.Topic()) {
// 		SL1 := bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(SL1), "kWh", "/Ac/L1/Energy/Forward")
// 		log.Debug(fmt.Sprintf("Energy Forward L1: %.3f kWh", SL1))
// 	}

// 	// /Ac/L2/Energy/Forward  <- kWh  - bought
// 	if ContainString(".*/Sum/L2$", msg.Topic()) {
// 		SL2 := bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(SL2), "kWh", "/Ac/L2/Energy/Forward")
// 		log.Debug(fmt.Sprintf("Energy Forward L2: %.3f kWh", SL2))
// 	}
// 	// /Ac/L3/Energy/Forward  <- kWh  - bought
// 	if ContainString(".*/Sum/L3$", msg.Topic()) {
// 		SL3 := bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(SL3), "kWh", "/Ac/L3/Energy/Forward")
// 		log.Debug(fmt.Sprintf("Energy Forward L3: %.3f kWh", SL3))
// 	}

// 	// /Ac/L1/Energy/Reverse  <- kWh  - bought
// 	if ContainString(".*/Export/L1$", msg.Topic()) {
// 		EL1 := bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(EL1), "kWh", "/Ac/L1/Energy/Reverse")
// 		log.Debug(fmt.Sprintf("Energy Reverse L1: %.3f kWh", EL1))
// 	}

// 	// /Ac/L2/Energy/Reverse  <- kWh  - bought
// 	if ContainString(".*/Export/L2$", msg.Topic()) {
// 		EL2 := bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(EL2), "kWh", "/Ac/L2/Energy/Reverse")
// 		log.Debug(fmt.Sprintf("Energy Reverse L2: %.3f kWh", EL2))
// 	}
// 	// /Ac/L3/Energy/Reverse  <- kWh  - bought
// 	if ContainString(".*/Export/L3$", msg.Topic()) {
// 		EL3 := bin2Float64(string(msg.Payload()))
// 		updateVariant(float64(EL3), "kWh", "/Ac/L3/Energy/Reverse")
// 		log.Debug(fmt.Sprintf("Energy Reverse L3: %.3f kWh", EL3))
// 	}
// }

package oldmain

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

/* Configuration */
var (
	broker     = "192.168.1.119"
	brokerPort = 1883
	topic      = "stromzaehler/#"
	clientId   = "sdm630-bridge"
	username   = "user"
	password   = "pass"
)

var P1 float64 = 0.00
var P2 float64 = 0.00
var P3 float64 = 0.00
var conn, err = dbus.SystemBus()

const intro = `
<node>
   <interface name="com.victronenergy.BusItem">
    <signal name="PropertiesChanged">
      <arg type="a{sv}" name="properties" />
    </signal>
    <method name="SetValue">
      <arg direction="in"  type="v" name="value" />
      <arg direction="out" type="i" />
    </method>
    <method name="GetText">
      <arg direction="out" type="s" />
    </method>
    <method name="GetValue">
      <arg direction="out" type="v" />
    </method>
    </interface>` + introspect.IntrospectDataString + `</node> `

type objectpath string

var victronValues = map[int]map[objectpath]dbus.Variant{
	// 0: This will be used to store the VALUE variant
	0: map[objectpath]dbus.Variant{},
	// 1: This will be used to store the STRING variant
	1: map[objectpath]dbus.Variant{},
}

func (f objectpath) GetValue() (dbus.Variant, *dbus.Error) {
	log.Debug("GetValue() called for ", f)
	log.Debug("...returning ", victronValues[0][f])
	return victronValues[0][f], nil
}
func (f objectpath) GetText() (string, *dbus.Error) {
	log.Debug("GetText() called for ", f)
	log.Debug("...returning ", victronValues[1][f])
	// Why does this end up ""SOMEVAL"" ... trim it I guess
	return strings.Trim(victronValues[1][f].String(), "\""), nil
}

func init() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		lvl = "info"
	}

	ll, err := log.ParseLevel(lvl)
	if err != nil {
		ll = log.DebugLevel
	}

	log.SetLevel(ll)
}

func main() {
	// Parse command line arguments
	flag.StringVar(&broker, "broker", broker, "MQTT broker address")
	flag.IntVar(&brokerPort, "port", brokerPort, "MQTT broker port")
	flag.StringVar(&topic, "topic", topic, "MQTT topic prefix")
	flag.StringVar(&clientId, "client-id", clientId, "MQTT client id")
	flag.StringVar(&username, "username", username, "MQTT username")
	flag.StringVar(&password, "password", password, "MQTT password")
	flag.Parse()

	// Need to implement following paths:
	// https://github.com/victronenergy/venus/wiki/dbus#grid-meter
	// also in system.py
	victronValues[0]["/Connected"] = dbus.MakeVariant(1)
	victronValues[1]["/Connected"] = dbus.MakeVariant("1")

	victronValues[0]["/CustomName"] = dbus.MakeVariant("Grid meter")
	victronValues[1]["/CustomName"] = dbus.MakeVariant("Grid meter")

	victronValues[0]["/DeviceInstance"] = dbus.MakeVariant(30)
	victronValues[1]["/DeviceInstance"] = dbus.MakeVariant("30")

	// also in system.py
	victronValues[0]["/DeviceType"] = dbus.MakeVariant(71)
	victronValues[1]["/DeviceType"] = dbus.MakeVariant("71")

	victronValues[0]["/ErrorCode"] = dbus.MakeVariantWithSignature(0, dbus.SignatureOf(123))
	victronValues[1]["/ErrorCode"] = dbus.MakeVariant("0")

	victronValues[0]["/FirmwareVersion"] = dbus.MakeVariant(2)
	victronValues[1]["/FirmwareVersion"] = dbus.MakeVariant("2")

	// also in system.py
	victronValues[0]["/Mgmt/Connection"] = dbus.MakeVariant("/dev/ttyUSB0")
	victronValues[1]["/Mgmt/Connection"] = dbus.MakeVariant("/dev/ttyUSB0")

	victronValues[0]["/Mgmt/ProcessName"] = dbus.MakeVariant("/opt/color-control/dbus-cgwacs/dbus-cgwacs")
	victronValues[1]["/Mgmt/ProcessName"] = dbus.MakeVariant("/opt/color-control/dbus-cgwacs/dbus-cgwacs")

	victronValues[0]["/Mgmt/ProcessVersion"] = dbus.MakeVariant("1.8.0")
	victronValues[1]["/Mgmt/ProcessVersion"] = dbus.MakeVariant("1.8.0")

	victronValues[0]["/Position"] = dbus.MakeVariantWithSignature(0, dbus.SignatureOf(123))
	victronValues[1]["/Position"] = dbus.MakeVariant("0")

	// also in system.py
	victronValues[0]["/ProductId"] = dbus.MakeVariant(45058)
	victronValues[1]["/ProductId"] = dbus.MakeVariant("45058")

	// also in system.py
	victronValues[0]["/ProductName"] = dbus.MakeVariant("Grid meter")
	victronValues[1]["/ProductName"] = dbus.MakeVariant("Grid meter")

	victronValues[0]["/Serial"] = dbus.MakeVariant("BP98305081235")
	victronValues[1]["/Serial"] = dbus.MakeVariant("BP98305081235")

	// Provide some initial values... note that the values must be a valid formt otherwise dbus_systemcalc.py exits like this:
	// @400000005ecc11bf3782b374   File "/opt/victronenergy/dbus-systemcalc-py/dbus_systemcalc.py", line 386, in _handletimertick
	// @400000005ecc11bf37aa251c     self._updatevalues()
	// @400000005ecc11bf380e74cc   File "/opt/victronenergy/dbus-systemcalc-py/dbus_systemcalc.py", line 678, in _updatevalues
	// @400000005ecc11bf383ab4ec     c = _safeadd(c, p, pvpower)
	// @400000005ecc11bf386c9674   File "/opt/victronenergy/dbus-systemcalc-py/sc_utils.py", line 13, in safeadd
	// @400000005ecc11bf387b28ec     return sum(values) if values else None
	// @400000005ecc11bf38b2bb7c TypeError: unsupported operand type(s) for +: 'int' and 'unicode'
	//
	victronValues[0]["/Ac/L1/Power"] = dbus.MakeVariant(0.0)
	victronValues[1]["/Ac/L1/Power"] = dbus.MakeVariant("0 W")
	victronValues[0]["/Ac/L2/Power"] = dbus.MakeVariant(0.0)
	victronValues[1]["/Ac/L2/Power"] = dbus.MakeVariant("0 W")
	victronValues[0]["/Ac/L3/Power"] = dbus.MakeVariant(0.0)
	victronValues[1]["/Ac/L3/Power"] = dbus.MakeVariant("0 W")

	victronValues[0]["/Ac/L1/Voltage"] = dbus.MakeVariant(230)
	victronValues[1]["/Ac/L1/Voltage"] = dbus.MakeVariant("230 V")
	victronValues[0]["/Ac/L2/Voltage"] = dbus.MakeVariant(230)
	victronValues[1]["/Ac/L2/Voltage"] = dbus.MakeVariant("230 V")
	victronValues[0]["/Ac/L3/Voltage"] = dbus.MakeVariant(230)
	victronValues[1]["/Ac/L3/Voltage"] = dbus.MakeVariant("230 V")

	victronValues[0]["/Ac/L1/Current"] = dbus.MakeVariant(0.0)
	victronValues[1]["/Ac/L1/Current"] = dbus.MakeVariant("0 A")
	victronValues[0]["/Ac/L2/Current"] = dbus.MakeVariant(0.0)
	victronValues[1]["/Ac/L2/Current"] = dbus.MakeVariant("0 A")
	victronValues[0]["/Ac/L3/Current"] = dbus.MakeVariant(0.0)
	victronValues[1]["/Ac/L3/Current"] = dbus.MakeVariant("0 A")

	victronValues[0]["/Ac/L1/Energy/Forward"] = dbus.MakeVariant(0.0)
	victronValues[1]["/Ac/L1/Energy/Forward"] = dbus.MakeVariant("0 kWh")
	victronValues[0]["/Ac/L2/Energy/Forward"] = dbus.MakeVariant(0.0)
	victronValues[1]["/Ac/L2/Energy/Forward"] = dbus.MakeVariant("0 kWh")
	victronValues[0]["/Ac/L3/Energy/Forward"] = dbus.MakeVariant(0.0)
	victronValues[1]["/Ac/L3/Energy/Forward"] = dbus.MakeVariant("0 kWh")

	victronValues[0]["/Ac/L1/Energy/Reverse"] = dbus.MakeVariant(0.0)
	victronValues[1]["/Ac/L1/Energy/Reverse"] = dbus.MakeVariant("0 kWh")
	victronValues[0]["/Ac/L2/Energy/Reverse"] = dbus.MakeVariant(0.0)
	victronValues[1]["/Ac/L2/Energy/Reverse"] = dbus.MakeVariant("0 kWh")
	victronValues[0]["/Ac/L3/Energy/Reverse"] = dbus.MakeVariant(0.0)
	victronValues[1]["/Ac/L3/Energy/Reverse"] = dbus.MakeVariant("0 kWh")

	basicPaths := []dbus.ObjectPath{
		"/Connected",
		"/CustomName",
		"/DeviceInstance",
		"/DeviceType",
		"/ErrorCode",
		"/FirmwareVersion",
		"/Mgmt/Connection",
		"/Mgmt/ProcessName",
		"/Mgmt/ProcessVersion",
		"/Position",
		"/ProductId",
		"/ProductName",
		"/Serial",
	}

	updatingPaths := []dbus.ObjectPath{
		"/Ac/L1/Power",
		"/Ac/L2/Power",
		"/Ac/L3/Power",
		"/Ac/L1/Voltage",
		"/Ac/L2/Voltage",
		"/Ac/L3/Voltage",
		"/Ac/L1/Current",
		"/Ac/L2/Current",
		"/Ac/L3/Current",
		"/Ac/L1/Energy/Forward",
		"/Ac/L2/Energy/Forward",
		"/Ac/L3/Energy/Forward",
		"/Ac/L1/Energy/Reverse",
		"/Ac/L2/Energy/Reverse",
		"/Ac/L3/Energy/Reverse",
	}

	defer conn.Close()

	// Some of the victron stuff requires it be called grid.cgwacs... using the only known valid value (from the simulator)
	// This can _probably_ be changed as long as it matches com.victronenergy.grid.cgwacs_*
	reply, err := conn.RequestName("com.victronenergy.grid.cgwacs_ttyUSB0_di30_mb1",
		dbus.NameFlagDoNotQueue)
	if err != nil {
		log.Panic("Something went horribly wrong in the dbus connection")
		panic(err)
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		log.Panic("name cgwacs_ttyUSB0_di30_mb1 already taken on dbus.")
		os.Exit(1)
	}

	for i, s := range basicPaths {
		log.Debug("Registering dbus basic path #", i, ": ", s)
		conn.Export(objectpath(s), s, "com.victronenergy.BusItem")
		conn.Export(introspect.Introspectable(intro), s, "org.freedesktop.DBus.Introspectable")
	}

	for i, s := range updatingPaths {
		log.Debug("Registering dbus update path #", i, ": ", s)
		conn.Export(objectpath(s), s, "com.victronenergy.BusItem")
		conn.Export(introspect.Introspectable(intro), s, "org.freedesktop.DBus.Introspectable")
	}

	log.Info("Successfully connected to dbus and registered as a meter... Commencing reading of the SDM630 meter")

	// MQTT Subscription
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, brokerPort))
	opts.SetClientID(clientId)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	sub(client)
	// Infinite loop
	for true {
		// fmt.Println("Infinite Loop entered")
		time.Sleep(time.Second)
	}

	// This is a forever loop^^
	panic("Error: We terminated.... how did we ever get here?")
}

/* MQTT Subscribe Function */
func sub(client mqtt.Client) {
	topic := topic
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	log.Info("Subscribed to topic: " + topic)
}

/* MQTT Publish Function */
func publish(client mqtt.Client) {
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := client.Publish("topic/test", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}

/* Write dbus Values to Victron handler */
func updateVariant(value float64, unit string, path string) {
	emit := make(map[string]dbus.Variant)
	emit["Text"] = dbus.MakeVariant(fmt.Sprintf("%.2f", value) + unit)
	emit["Value"] = dbus.MakeVariant(float64(value))
	victronValues[0][objectpath(path)] = emit["Value"]
	victronValues[1][objectpath(path)] = emit["Text"]
	conn.Emit(dbus.ObjectPath(path), "com.victronenergy.BusItem.PropertiesChanged", emit)
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

/* Called if connection is established */
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Info(fmt.Sprintf("Connected to broker %s:%d", broker, brokerPort))
}

/* Called if connection is lost  */
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Info(fmt.Sprintf("Connect lost: %v", err))
	os.Exit(1)
}

/* Search for string with regex */
func ContainString(searchstring string, str string) bool {
	var obj bool

	obj, err = regexp.MatchString(searchstring, str)

	if err != nil {
		panic(err)
	}

	return obj
}

/* MQTT Subscribe Handler */
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

	log.Debug(fmt.Sprintf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic()))

	// Power L1
	if ContainString(".*Power/L1$", msg.Topic()) {
		P1 = bin2Float64(string(msg.Payload()))
		updateVariant(float64(P1), "W", "/Ac/L1/Power")
	}

	// Power L2
	if ContainString(".*Power/L2$", msg.Topic()) {
		P2 = bin2Float64(string(msg.Payload()))
		updateVariant(float64(P2), "W", "/Ac/L2/Power")
	}

	// Power L3
	if ContainString(".*Power/L3$", msg.Topic()) {
		P3 = bin2Float64(string(msg.Payload()))
		updateVariant(float64(P3), "W", "/Ac/L3/Power")
	}

	// /Ac/Energy/Forward     <- kWh  - bought energy (total of all phases)
	if ContainString(".*/Import$", msg.Topic()) {
		IP := bin2Float64(string(msg.Payload()))
		updateVariant(float64(IP), "kWh", "/Ac/Energy/Forward")
		log.Debug(fmt.Sprintf("Import Power: %.3f kWh", IP))
	}

	// /Ac/Energy/Reverse     <- kWh  - sold energy (total of all phases)
	if ContainString(".*/Export$", msg.Topic()) {
		EP := bin2Float64(string(msg.Payload()))
		updateVariant(float64(EP), "kWh", "/Ac/Energy/Reverse")
		log.Debug(fmt.Sprintf("Export Power: %.3f kWh", EP))
	}

	// /Ac/Power              <- W    - total of all phases, real power
	if ContainString(".*/Power$", msg.Topic()) {
		TP := bin2Float64(string(msg.Payload()))
		updateVariant(float64(TP), "W", "/Ac/Power")
		log.Debug(fmt.Sprintf("Total Power: %.3f W", TP))
	}

	// /Ac/L1/Current         <- A AC
	if ContainString(".*/Current/L1$", msg.Topic()) {
		CL1 := bin2Float64(string(msg.Payload()))
		updateVariant(float64(CL1), "A", "/Ac/L1/Current")
		log.Debug(fmt.Sprintf("Current L1: %.3f A", CL1))
	}

	// /Ac/L2/Current         <- A AC
	if ContainString(".*/Current/L2$", msg.Topic()) {
		CL2 := bin2Float64(string(msg.Payload()))
		updateVariant(float64(CL2), "A", "/Ac/L2/Current")
		log.Debug(fmt.Sprintf("Current L2: %.3f A", CL2))
	}

	// /Ac/L3/Current         <- A AC
	if ContainString(".*/Current/L3$", msg.Topic()) {
		CL3 := bin2Float64(string(msg.Payload()))
		updateVariant(float64(CL3), "A", "/Ac/L3/Current")
		log.Debug(fmt.Sprintf("Current L3: %.3f A", CL3))
	}

	// /Ac/L1/Voltage <- V AC
	if ContainString(".*/Voltage/L1$", msg.Topic()) {
		VL1 := bin2Float64(string(msg.Payload()))
		updateVariant(float64(VL1), "V", "/Ac/L1/Voltage")
		log.Debug(fmt.Sprintf("Voltage L1: %.3f V", VL1))
	}

	// /Ac/L2/Voltage <- V AC
	if ContainString(".*/Voltage/L2$", msg.Topic()) {
		VL2 := bin2Float64(string(msg.Payload()))
		updateVariant(float64(VL2), "V", "/Ac/L2/Voltage")
		log.Debug(fmt.Sprintf("Voltage L2: %.3f V", VL2))
	}

	// /Ac/L3/Voltage <- V AC
	if ContainString(".*/Voltage/L3$", msg.Topic()) {
		VL3 := bin2Float64(string(msg.Payload()))
		updateVariant(float64(VL3), "V", "/Ac/L3/Voltage")
		log.Debug(fmt.Sprintf("Voltage L3: %.3f V", VL3))
	}

	// /Ac/L1/Energy/Forward  <- kWh  - bought
	if ContainString(".*/Sum/L1$", msg.Topic()) {
		SL1 := bin2Float64(string(msg.Payload()))
		updateVariant(float64(SL1), "kWh", "/Ac/L1/Energy/Forward")
		log.Debug(fmt.Sprintf("Energy Forward L1: %.3f kWh", SL1))
	}

	// /Ac/L2/Energy/Forward  <- kWh  - bought
	if ContainString(".*/Sum/L2$", msg.Topic()) {
		SL2 := bin2Float64(string(msg.Payload()))
		updateVariant(float64(SL2), "kWh", "/Ac/L2/Energy/Forward")
		log.Debug(fmt.Sprintf("Energy Forward L2: %.3f kWh", SL2))
	}
	// /Ac/L3/Energy/Forward  <- kWh  - bought
	if ContainString(".*/Sum/L3$", msg.Topic()) {
		SL3 := bin2Float64(string(msg.Payload()))
		updateVariant(float64(SL3), "kWh", "/Ac/L3/Energy/Forward")
		log.Debug(fmt.Sprintf("Energy Forward L3: %.3f kWh", SL3))
	}

	// /Ac/L1/Energy/Reverse  <- kWh  - bought
	if ContainString(".*/Export/L1$", msg.Topic()) {
		EL1 := bin2Float64(string(msg.Payload()))
		updateVariant(float64(EL1), "kWh", "/Ac/L1/Energy/Reverse")
		log.Debug(fmt.Sprintf("Energy Reverse L1: %.3f kWh", EL1))
	}

	// /Ac/L2/Energy/Reverse  <- kWh  - bought
	if ContainString(".*/Export/L2$", msg.Topic()) {
		EL2 := bin2Float64(string(msg.Payload()))
		updateVariant(float64(EL2), "kWh", "/Ac/L2/Energy/Reverse")
		log.Debug(fmt.Sprintf("Energy Reverse L2: %.3f kWh", EL2))
	}
	// /Ac/L3/Energy/Reverse  <- kWh  - bought
	if ContainString(".*/Export/L3$", msg.Topic()) {
		EL3 := bin2Float64(string(msg.Payload()))
		updateVariant(float64(EL3), "kWh", "/Ac/L3/Energy/Reverse")
		log.Debug(fmt.Sprintf("Energy Reverse L3: %.3f kWh", EL3))
	}
}
