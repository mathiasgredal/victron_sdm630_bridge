package dbus

import (
	"strings"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	log "github.com/sirupsen/logrus"
)

var introspection = &introspect.Node{
	Interfaces: []introspect.Interface{
		introspect.IntrospectData,
		{
			Name: "com.victronenergy.BusItem",
			Signals: []introspect.Signal{
				{
					Name: "PropertiesChanged",
					Args: []introspect.Arg{
						{
							Name: "properties",
							Type: "a{sv}",
						},
					},
				},
			},
			Methods: []introspect.Method{
				{
					Name: "SetValue",
					Args: []introspect.Arg{
						{
							Name:      "value",
							Type:      "v",
							Direction: "in",
						},
						{
							Type:      "i",
							Direction: "out",
						},
					},
				},
				{
					Name: "GetText",
					Args: []introspect.Arg{
						{
							Type:      "s",
							Direction: "out",
						},
					},
				},
				{
					Name: "GetValue",
					Args: []introspect.Arg{
						{
							Type:      "v",
							Direction: "out",
						},
					},
				},
			},
		},
	},
}

type ObjectPath struct {
	path  dbus.ObjectPath
	state *State
}

func (f ObjectPath) GetValue() (dbus.Variant, *dbus.Error) {
	log.Debug("GetValue() called for ", f.path)
	log.Debug("...returning ", f.state.Items[f.path].Value)
	return f.state.Items[f.path].Value, nil
}

func (f ObjectPath) GetText() (string, *dbus.Error) {
	log.Debug("GetText() called for ", f)
	log.Debug("...returning ", f.state.Items[f.path].Text)
	return strings.Trim(f.state.Items[f.path].Text.String(), "\""), nil
}
