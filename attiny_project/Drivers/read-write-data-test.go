package main

import (
	"fmt"
	"log"
	"time"

	"periph.io/x/host/v3"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

const (
	i2cAddress = 0x04
	devPort = "/dev/i2c-1"
)



func main() {
	\
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

	sendAndReceive := func(data byte) (byte, error) {
		
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

	dataToSend := byte(40)
	receivedData, err := sendAndReceive(dataToSend)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Sent: \n", dataToSend)
	fmt.Printf("Received: \n", receivedData)
}


