package mmwaveread

// import (
// 	"bytes"
// 	"encoding/binary"
// 	"fmt"
// 	"time"
// )

func jack() {

}

// type Sensor struct {
// 	magicWord  []byte
// 	timeout    time.Duration
// 	msPerFrame int
// 	dataPort   *bytes.Buffer
// 	start      time.Time
// 	bufferTemp *bytes.Buffer
// 	shutDown   bool
// }

// func NewSensor(magicWord []byte, timeout time.Duration, msPerFrame int) *Sensor {
// 	return &Sensor{
// 		magicWord:  magicWord,
// 		timeout:    timeout,
// 		msPerFrame: msPerFrame,
// 		dataPort:   &bytes.Buffer{},
// 		start:      time.Now(),
// 		bufferTemp: &bytes.Buffer{},
// 	}
// }

// func (s *Sensor) parseData(buffer []byte) [][]float64 {
// 	start := bytes.Index(buffer, s.magicWord) + 40
// 	var tlvType, tlvLength uint32
// 	err := binary.Read(bytes.NewReader(buffer[start:start+8]), binary.LittleEndian, &tlvType)
// 	if err != nil {
// 		fmt.Println("Parsing tlvType failed")
// 		return nil
// 	}
// 	err = binary.Read(bytes.NewReader(buffer[start+4:start+8]), binary.LittleEndian, &tlvLength)
// 	if err != nil {
// 		fmt.Println("Parsing tlvLength failed")
// 		return nil
// 	}
// 	numPoints := int(tlvLength / 16)
// 	data := make([][]float64, numPoints)
// 	if tlvType == 1 {
// 		for j := 0; j < numPoints; j++ {
// 			var x, y, z, vel float32
// 			err = binary.Read(bytes.NewReader(buffer[start+8+j*16:start+12+j*16]), binary.LittleEndian, &x)
// 			if err != nil {
// 				fmt.Println("Parsing x failed")
// 				return nil
// 			}
// 			err = binary.Read(bytes.NewReader(buffer[start+12+j*16:start+16+j*16]), binary.LittleEndian, &y)
// 			if err != nil {
// 				fmt.Println("Parsing y failed")
// 				return nil
// 			}
// 			err = binary.Read(bytes.NewReader(buffer[start+16+j*16:start+20+j*16]), binary.LittleEndian, &z)
// 			if err != nil {
// 				fmt.Println("Parsing z failed")
// 				return nil
// 			}
// 			err = binary.Read(bytes.NewReader(buffer[start+20+j*16:start+24+j*16]), binary.LittleEndian, &vel)
// 			if err != nil {
// 				fmt.Println("Parsing vel failed")
// 				return nil
// 			}
// 			data[j] = []float64{float64(y), -float64(x), float64(z), float64(vel)}
// 		}
// 	} else {
// 		fmt.Println("OUTPUT_MSG_DETECTED_POINTS WRONG")
// 		return nil
// 	}

// 	// For SDK 3.x, intensity is replaced by snr in sideInfo and is parsed in the READ_SIDE_INFO code
// 	err = binary.Read(bytes.NewReader(buffer[start+8+numPoints*16:start+12+numPoints*16]), binary.LittleEndian, &tlvType)
// 	if err != nil {
// 		fmt.Println("Parsing tlvType 2 failed")
// 		return nil
// 	}
// 	if tlvType == 7 {
// 		for j := 0; j < numPoints; j++ {
// 			var snr int16
// 			err = binary.Read(bytes.NewReader(buffer[start+12+numPoints*16+j*4:start+14+numPoints*16+j*4]), binary.LittleEndian, &snr)
// 			if err != nil {
// 				fmt.Println("Parsing snr failed")
// 				return nil
// 			}
// 			data[j][3] = float64(snr) / 10.0
// 		}
// 	} else {
// 		fmt.Println("OUTPUT_MSG_DETECTED_POINTS_SIDE_INFO WRONG")
// 		return nil
// 	}
// 	return data
// }

// func (s *Sensor) Read() <-chan [][]float64 {
// 	msgChan := make(chan [][]float64)
// 	go func() {
// 		defer close(msgChan)
// 		for time.Now().Before(s.start.Add(s.timeout)) {
// 			s.bufferTemp.Write(s.dataPort.Bytes())
// 			if s.bufferTemp.Len() != 0 {
// 				idxStart := bytes.Index(s.bufferTemp.Bytes(), s.magicWord)
// 				if idxStart != -1 {
// 					idxEnd := bytes.Index(s.bufferTemp.Bytes()[idxStart+1:], s.magicWord)
// 					if idxEnd != -1 {
// 						msg := s.parseData(s.bufferTemp.Bytes()[idxStart : idxStart+idxEnd])
// 						msgChan <- msg
// 						s.resetTimer()
// 						s.clearBuffer()
// 					}
// 				}
// 			} else {
// 				time.Sleep(time.Duration(s.msPerFrame) * time.Millisecond)
// 			}
// 			if s.shutDown {
// 				break
// 			}
// 		}
// 		s.close()
// 	}()
// 	return msgChan
// }

// func (s *Sensor) resetTimer() {
// 	s.start = time.Now()
// }

// func (s *Sensor) clearBuffer() {
// 	s.bufferTemp.Reset()
// }

// func (s *Sensor) close() {
// 	// Close any resources if needed
// }

// func charli() {
// 	magicWord := []byte{0x01, 0x02} // Change to your magic word
// 	timeout := 10 * time.Second     // Change to your desired timeout
// 	msPerFrame := 100               // Change to your desired ms per frame

// 	sensor := NewSensor(magicWord, timeout, msPerFrame)
// 	msgChannel := sensor.Read()

// 	for msg := range msgChannel {
// 		fmt.Println(msg)
// 	}
// }
