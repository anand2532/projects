package mmwaveread

import (
	"bytes"
	"log"
	"time"

	"github.com/tarm/serial"
)

// const(
// 	var MagicWord []byte = []byte{0x02, 0x01, 0x04, 0x03, 0x06, 0x05, 0x08, 0x07}
// 	Timeout := 10 * time.Second
// 	msPerFrame := 100
// )

const (
	Timeout    = 10 * time.Second
	msPerFrame = 100
)

var MagicWord []byte = []byte{0x02, 0x01, 0x04, 0x03, 0x06, 0x05, 0x08, 0x07}

type Driver struct {
	port       *serial.Port
	MagicWord  []byte
	Timeout    time.Duration
	msPerFrame int
	dataPort   *bytes.Buffer
	start      time.Time
	bufferTemp *bytes.Buffer
}

func NewSensor(MagicWord []byte, Timeout time.Duration, msPerFrame int) *Driver {
	return &Driver{
		MagicWord:  MagicWord,
		Timeout:    Timeout,
		msPerFrame: msPerFrame,
		dataPort:   &bytes.Buffer{},
		start:      time.Now(),
		bufferTemp: &bytes.Buffer{},
	}
}

func NewDriver(port *serial.Port) *Driver {
	return &Driver{
		port: port,
	}
}

func (d *Driver) Start() {
	for {
		go d.readData()
	}
}

func (d *Driver) readData() {
	buffer := make([]byte, 0)

	for {
		buf := make([]byte, 1024)
		n, err := d.port.Read(buf)
		if err != nil {
			continue
		}

		buffer := append(buffer, buf[:n]...)
		log.Println(buffer)
		// return buf[:n]
		// go d.parseData(buf[:n])
	}
}

// func (d *Driver) parseData(buffer []byte) {
// 	start := bytes.Index(buffer, d.MagicWord) + 40
// 	log.Println(start)

// 	var tlvType, tlvLength uint32

// 	err := binary.Read(bytes.NewReader(buffer[start:start+8]), binary.LittleEndian, &tlvType)
// 	if err != nil {
// 		log.Println("Parsing tlvType failed")
// 		// return nil
// 	}

// 	err = binary.Read(bytes.NewReader(buffer[start+4:start+8]), binary.LittleEndian, &tlvLength)
// 	if err != nil {
// 		log.Println("Parsing tlvType failed")
// 		// return nil
// 	}

// 	numPoints := int(tlvLength / 16)
// 	data := make([][]float64, numPoints)

// 	if tlvType == 1 {
// 		for j := 0; j < numPoints; j++ {

// 			var x, y, z, vel float32

// 			reader := bytes.NewReader(buffer[start+8+j*16 : start+12+j*16])
// 			err = binary.Read(reader, binary.LittleEndian, &x)
// 			if err != nil {
// 				log.Printf("Parsing x failed: %v", err)
// 			}

// 			reader = bytes.NewReader(buffer[start+12+j*16 : start+16+j*16])
// 			err = binary.Read(reader, binary.LittleEndian, &y)
// 			if err != nil {
// 				log.Printf("Parsing y failed: %v", err)
// 			}

// 			reader = bytes.NewReader(buffer[start+16+j*16 : start+20+j*16])
// 			err = binary.Read(reader, binary.LittleEndian, &z)
// 			if err != nil {
// 				log.Printf("Parsing z failed: %v", err)
// 			}

// 			reader = bytes.NewReader(buffer[start+20+j*16 : start+24+j*16])
// 			err = binary.Read(reader, binary.LittleEndian, &vel)
// 			if err != nil {
// 				log.Printf("Parsing vel failed: %v", err)
// 			}
// 			data[j] = []float64{float64(y), -float64(x), float64(z), float64(vel)}
// 			log.Println(data)
// 		}
// 	} else {
// 		log.Println("wrong")
// 		return nil
// 	}

// 	err = binary.Read(bytes.NewReader(buffer[start+8+numPoints*16:start+12+numPoints*16]), binary.LittleEndian, &tlvType)
// 	if err != nil {
// 		log.Println("Parsing tlv2 failed")
// 		return nil
// 	}

// 	if tlvType == 7 {
// 		for j := 0; j < numPoints; j++ {
// 			var snr int16
// 			err = binary.Read(bytes.NewReader(buffer[start+12+numPoints*16+j*4:start+14+numPoints*16+j*4]), binary.LittleEndian, &snr)
// 			if err != nil {
// 				log.Println("Parsing snr failed")
// 				// return nil
// 			}
// 			data[j][3] = float64(snr) / 10.0
// 		}
// 	} else {
// 		log.Println("wrong message")
// 		// return nil
// 	}
// 	// return data
// }

// func (d *Driver) Read() {
// 	msgChan := make(chan [][]float64)

// 	go func() {
// 		defer close(msgChan)
// 		for time.Now().Before(d.start.Add(d.Timeout)) {
// 			d.bufferTemp.Write(d.dataPort.Bytes())
// 			if d.bufferTemp.Len() != 0 {
// 				idxStart := bytes.Index(d.bufferTemp.Bytes(), d.MagicWord)
// 				if idxStart != -1 {
// 					idxEnd := bytes.Index(d.bufferTemp.Bytes()[idxStart+1:], d.MagicWord)
// 					if idxEnd != -1 {
// 						msg := d.parseData(d.bufferTemp.Bytes()[idxStart : idxStart+idxEnd])
// 						log.Println(msg)
// 						// msgChan <- msg
// 						d.resetTimer()
// 						d.clearBuffer()
// 					}
// 				}
// 			} else {
// 				time.Sleep(time.Duration(d.msPerFrame) * time.Millisecond)
// 			}
// 			// if d.shutDown {
// 			// 	break
// 			// }
// 		}
// 		d.close()
// 	}()
// 	// return msgChan
// }

// func (d *Driver) resetTimer() {
// 	d.start = time.Now()
// }

// func (d *Driver) clearBuffer() {
// 	d.bufferTemp.Reset()
// }

// func (d *Driver) close() {
// 	// Close any resources if needed
// }
