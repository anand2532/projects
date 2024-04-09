package main

import (
	"fmt"
	"log"
	"time"
	"os"

	"periph.io/x/host/v3"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

const (
	i2cAddress = 0x04
	devPort = "/dev/i2c-1"
)

var (
	isPaired = false
	Code uint8 = 0
	pairingCode = &Code

)

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
	}else if len(os.Args) > 1 && os.Args[1] == "pair" {
        // Send the value 60 and receive some data
        dataToSend := byte(60)
        receivedData, err := sendAndReceive(d, dataToSend)
        if err != nil {
            log.Printf("Error sending and receiving data: %s", err)
        } else {
            *pairingCode = receivedData
			log.Println("Pairing Code", *pairingCode)
			isPaired = true
        }
        os.Exit(0)
    }

	if isPaired == true {

		dataToSend := byte(61)
		receivedData, err := sendAndReceive(d, dataToSend)
		if err != nil {
			log.Fatal(err)
		}
        
		if receivedData != *pairingCode {
			fmt.Println("Received: Tampered")
		}else {
			log.Println("Pairing Done")
			for {
				dataToSend := byte(10)
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
		
	}else {
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
}

func sendAndReceive(d *i2c.Dev, data byte) (byte, error) {

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





