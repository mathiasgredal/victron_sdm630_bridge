package dbus

import (
	log "github.com/sirupsen/logrus"

	"github.com/godbus/dbus"
)

/* dbus state */
type State struct {
	Conn        *dbus.Conn
	Items       map[dbus.ObjectPath]Item
	DCPower     float64 // DC power value(W)
	LastDCPower float64 // DC power value from last mqtt update(W)
}

/* Initialize dbus State */
func New() *State {
	conn, err := dbus.SystemBus()
	if err != nil {
		panic(err)
	}

	state := &State{
		Conn: conn,
		Items: map[dbus.ObjectPath]Item{
			Connected:          MakeValueItem("1", 1),
			CustomName:         MakeTextItem("Grid meter"),
			DeviceInstance:     MakeValueItem("30", 30),
			DeviceType:         MakeValueItem("71", 71),
			ErrorCode:          MakeSignatureItem("0", 0, 123),
			FirmwareVersion:    MakeValueItem("2", 2),
			MgmtConnection:     MakeTextItem("/dev/ttyUSB0"),
			MgmtProcessName:    MakeTextItem("/opt/color-control/dbus-cgwacs/dbus-cgwacs"),
			MgmtProcessVersion: MakeTextItem("1.8.0"),
			Position:           MakeSignatureItem("0", 0, 123),
			ProductId:          MakeValueItem("45058", 45058),
			ProductName:        MakeTextItem("Grid meter"),
			Serial:             MakeTextItem("BP98305081235"),
			AcL1Power:          MakeValueItem("0 W", 0.0),
			AcL2Power:          MakeValueItem("0 W", 0.0),
			AcL3Power:          MakeValueItem("0 W", 0.0),
			// AcL1Voltage:        MakeValueItem("230 V", 230),
			// AcL2Voltage:        MakeValueItem("230 V", 230),
			// AcL3Voltage:        MakeValueItem("230 V", 230),
			// AcL1Current:        MakeValueItem("0 A", 0.0),
			// AcL2Current:        MakeValueItem("0 A", 0.0),
			// AcL3Current:        MakeValueItem("0 A", 0.0),
			// AcL1EnergyForward:  MakeValueItem("0 kWh", 0.0),
			// AcL2EnergyForward:  MakeValueItem("0 kWh", 0.0),
			// AcL3EnergyForward:  MakeValueItem("0 kWh", 0.0),
			// AcL1EnergyReverse:  MakeValueItem("0 kWh", 0.0),
			// AcL2EnergyReverse:  MakeValueItem("0 kWh", 0.0),
			// AcL3EnergyReverse:  MakeValueItem("0 kWh", 0.0),
		},
	}

	reply, err := conn.RequestName("com.victronenergy.grid.cgwacs_ttyUSB0_di30_mb1", dbus.NameFlagDoNotQueue)

	if err != nil {
		log.Panic("Something went horribly wrong in the dbus connection")
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		log.Panic("name cgwacs_ttyUSB0_di30_mb1 already taken on dbus.")
	}

	for path := range state.Items {
		log.Debug("Registering dbus path: ", path)
		conn.Export(ObjectPath{path, state}, path, "com.victronenergy.BusItem")
		conn.Export(introspection, path, "org.freedesktop.DBus.Introspectable")
	}

	return state
}

/* Send dbus notification about a state change */
func (s *State) Update(path dbus.ObjectPath, item Item) {
	s.Items[path] = item
	emit := make(map[string]dbus.Variant)
	emit["Text"] = item.Text
	emit["Value"] = item.Value
	s.Conn.Emit(dbus.ObjectPath(path), "com.victronenergy.BusItem.PropertiesChanged", emit)
}
