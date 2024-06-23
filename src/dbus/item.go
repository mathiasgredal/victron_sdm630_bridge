package dbus

import "github.com/godbus/dbus"

/* dbus item */
type Item struct {
	Text  dbus.Variant
	Value dbus.Variant
}

/* Item Text constructor */
func MakeTextItem(text string) Item {
	return Item{
		Text:  dbus.MakeVariant(text),
		Value: dbus.MakeVariant(text),
	}
}

/* Item Value constructor */
func MakeValueItem(text string, value interface{}) Item {
	return Item{
		Text:  dbus.MakeVariant(text),
		Value: dbus.MakeVariant(value),
	}
}

/* Item Signature constructor */
func MakeSignatureItem(text string, value interface{}, signature interface{}) Item {
	return Item{
		Text:  dbus.MakeVariant(text),
		Value: dbus.MakeVariantWithSignature(value, dbus.SignatureOf(signature)),
	}
}
