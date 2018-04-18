package service

import (
	"github.com/muka/go-bluetooth/bluez"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
	"github.com/fatih/structs"
)


type LEAdvertisement1Config struct {
	app 			*Application
	objectPath		dbus.ObjectPath
	conn 			*dbus.Conn
	device_path		dbus.ObjectPath
}

type LEAdvertisement1Properties struct {
	Type				string
	ServiceUUIDs		[]string
	ManufacturerData	map[string][]byte
	SolicitUUIDs		[]string
	ServiceData			map[string][]byte
	Includes			[]string
	LocalName			string
	Appearance			uint16
	Duration			uint16
	Timeout				uint16
}

func (p *LEAdvertisement1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(p), nil
}

func NewLEAdvertisement1(config *LEAdvertisement1Config, props *LEAdvertisement1Properties) (*LEAdvertisement1, error) {
	propInterface, err := NewProperties(config.conn)
	if err != nil {
		return nil, err
	}

	a := &LEAdvertisement1{
		config:              config,
		properties:          props,
		PropertiesInterface: propInterface,
	}

	err = propInterface.AddProperties(a.Interface(), props)
	if err != nil {
		return nil, err
	}

	return a, nil
}

type LEAdvertisement1 struct {
	config 				*LEAdvertisement1Config
	properties			*LEAdvertisement1Properties
	PropertiesInterface *Properties
}

//Properties return the properties of the service
func (ad *LEAdvertisement1) Properties() map[string]bluez.Properties {
	p := make(map[string]bluez.Properties)
	p[ad.Interface()] = ad.properties
	return p
}


func (ad *LEAdvertisement1) Interface() string {
	return bluez.LEAdvertisement1Interface
}

func (ad *LEAdvertisement1) Application() *Application {
	return ad.config.app
}

func (ad *LEAdvertisement1) Path() dbus.ObjectPath {
	return ad.config.objectPath
}

func (ad *LEAdvertisement1) Release() *dbus.Error {
	// TODO: Release anything we need?
	return nil
}

//Expose the service to dbus
func (ad *LEAdvertisement1) Expose() error {

	conn := ad.config.conn

	err := conn.Export(ad, ad.Path(), ad.Interface())
	if err != nil {
		return err
	}

	for iface, props := range ad.Properties() {
		ad.PropertiesInterface.AddProperties(iface, props)
	}

	ad.PropertiesInterface.Expose(ad.Path())

	node := &introspect.Node{
		Interfaces: []introspect.Interface{
			//Introspect
			introspect.IntrospectData,
			//Properties
			prop.IntrospectData,
			//GattService1
			{
				Name:       ad.Interface(),
				Methods:    introspect.Methods(ad),
				Properties: ad.PropertiesInterface.Introspection(ad.Interface()),
			},
		},
	}

	err = conn.Export(
		introspect.NewIntrospectable(node),
		ad.Path(),
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		return err
	}

	return nil
}
