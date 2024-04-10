package main

import (
	"log"

	"github.com/mavneo/Obstacle-Avoidance/tree/mmwave_driver/mmwave_driver/go_driver/configuration"
	mmwaveread "github.com/mavneo/Obstacle-Avoidance/tree/mmwave_driver/mmwave_driver/go_driver/mmwaveRead"
	"github.com/tarm/serial"
)

type Driver struct {
	CommandPort  *serial.Port
	DataPort     *serial.Port
	ConfigDriver *configuration.Driver
	ReadDriver   *mmwaveread.Driver
}

func NewDriverConfig(commandPort, dataPort *serial.Port) *Driver {
	return &Driver{
		CommandPort: commandPort,
		DataPort:    dataPort,
	}
}

func (d *Driver) InitializeDrivers() error {
	d.ConfigDriver = configuration.NewDriver(d.CommandPort)
	d.ConfigDriver.LoadCfgs()

	d.ReadDriver = mmwaveread.NewDriver(d.DataPort)

	return nil
}

func (d *Driver) RunDrivers() {
	d.ReadDriver.Start()
}

func main() {
	commandPort, dataPort, err := initializePorts()
	if err != nil {
		log.Fatalf("Failed to initialize ports: %v", err)
	}
	defer func() {
		closePort(commandPort)
		closePort(dataPort)
	}()

	driverConfig := NewDriverConfig(commandPort, dataPort)

	err = driverConfig.InitializeDrivers()
	if err != nil {
		log.Fatalf("Failed to initialize drivers: %v", err)
	}

	driverConfig.RunDrivers()
}

func initializePorts() (*serial.Port, *serial.Port, error) {
	commandPort, err := openPort("/dev/ttyUSB0", 115200)
	if err != nil {
		return nil, nil, err
	}

	dataPort, err := openPort("/dev/ttyUSB1", 921600)
	if err != nil {
		closePort(commandPort)
		return nil, nil, err
	}

	return commandPort, dataPort, nil
}

func openPort(name string, baud int) (*serial.Port, error) {
	port, err := serial.OpenPort(&serial.Config{Name: name, Baud: baud})
	if err != nil {
		return nil, err
	}
	return port, nil
}

func closePort(port *serial.Port) {
	if port != nil {
		port.Close()
	}
}
