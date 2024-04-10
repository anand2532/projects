// Main Package to read

// package main

// import (
// 	"log"

// 	"github.com/mavneo/Obstacle-Avoidance/tree/mmwave_driver/mmwave_driver/go_driver/configuration"
// 	mmwaveread "github.com/mavneo/Obstacle-Avoidance/tree/mmwave_driver/mmwave_driver/go_driver/mmwaveRead"
// 	"github.com/tarm/serial"
// )

// func main() {
// 	commandPort, err := serial.OpenPort(&serial.Config{Name: "/dev/ttyUSB0", Baud: 115200})
// 	if err != nil {
// 		log.Fatalf("Failed to open write port: %v", err)
// 	}
// 	defer commandPort.Close()
// 	dataPort, err := serial.OpenPort(&serial.Config{Name: "/dev/ttyUSB1", Baud: 921600})
// 	if err != nil {
// 		log.Fatalf("Failed to open write port: %v", err)
// 	}
// 	defer dataPort.Close()
// 	driver := configuration.NewDriver(commandPort)
// 	driver.LoadCfgs()

// 	readdriver := mmwaveread.NewDriver(dataPort)
// 	readdriver.Start()

// }
