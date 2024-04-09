package main

import (
	"fmt"
	"log"
	"os"
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

	if len(os.Args) > 1 && os.Args[1] == "reset" {
		err := testWrite(d, 40)
		if err != nil {
			log.Printf("Error writing data: %s", err)
		}
		os.Exit(0)
	}

	for {
		dataToSend := byte(41)
		receivedData, err := sendAndReceive(d, dataToSend)
		if err != nil {
			log.Fatal(err)
		}

		if receivedData == 20 {
			fmt.Println("Received: Tampered")
		} else {
			fmt.Println("Received: Untampered")
		}

		time.Sleep(time.Second)
	}
}

func sendAndReceive(d *i2c.Dev, data byte) (byte, error) {

	i2cMutex.Lock()
	defer i2cMutex.Unlock()

	write := []byte{data}
	if err := d.Tx(write, nil); err != nil {
		return 0, err
	}

	time.Sleep(time.Millisecond * 10)

	read := make([]byte, 1)
	if err := d.Tx(nil, read); err != nil {
		return 0, err
	}

	return read[0], nil
}

func testWrite(device *i2c.Dev, value byte) error {
	writeBuffer := []byte{value}

	_, err := device.Write(writeBuffer)

	return err
}



