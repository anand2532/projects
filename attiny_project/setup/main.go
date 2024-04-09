package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

const (
	i2cAddress = 0x04
	devPort    = "/dev/i2c-1"
)

var i2cMutex sync.Mutex

func main() {

	if _, err := host.Init(); err != nil {
		log.Fatalf("Failed to initialize periph: %v", err)
	}

	log.Println("starting")
	b, err := i2creg.Open(devPort)
	if err != nil {
		log.Fatal("Failed to open I2C bus: %v", err)
	}
	defer b.Close()

	d := &i2c.Dev{Addr: i2cAddress, Bus: b}

	go periodicallyGetTamperState(d)
	go periodicallyResetTamperState(d)

	for {
		time.Sleep(1000000000 * time.Second)
	}
}

func periodicallyGetTamperState(d *i2c.Dev) {
	i := 0
	for {
		i += 1
		time.Sleep(1550 * time.Millisecond)

		fmt.Println(i, "Get: Start")

		tamperState, err := sendAndReceive(d, byte(41))
		if err != nil {
			fmt.Println(i, "Get: Error:", err)
			continue
		}

		fmt.Println(i, "Get: Received:", tamperState)
	}
}

func periodicallyResetTamperState(d *i2c.Dev) {
	i := 0
	for {
		i += 1
		time.Sleep(5000 * time.Millisecond)

		fmt.Println(i, "Reset: Start")

		_, err := sendAndReceive(d, byte(40))
		if err != nil {
			fmt.Println(i, "Reset: Error:", err)
			continue
		}

		fmt.Println(i, "Reset: Done")
	}
}

func sendAndReceive(d *i2c.Dev, data byte) (byte, error) {

	// i2cMutex.Lock()
	// defer i2cMutex.Unlock()

	write := []byte{data}
	read := make([]byte, 1)

	if err := d.Tx(write, read); err != nil {
		return 0, err
	}

	// time.Sleep(time.Millisecond * 550)
	// if err := d.Tx(nil, read); err != nil {
	// 	return 0, err
	// }

	return read[0], nil
}

func testWrite(device *i2c.Dev, value byte) error {
	// i2cMutex.Lock()
	// defer i2cMutex.Unlock()

	writeBuffer := []byte{value}

	_, err := device.Write(writeBuffer)

	return err
}
