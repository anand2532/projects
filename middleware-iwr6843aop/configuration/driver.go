package configuration

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tarm/serial"
)

type Driver struct {
	port *serial.Port
}

func NewDriver(port *serial.Port) *Driver {
	return &Driver{
		port: port,
	}
}

func (d *Driver) LoadCfgs() {
	tiConfig, err := d.readConfigFile()
	if err != nil {
		log.Println(err)
		return
	}

	err = d.Configure(tiConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	// d.processConfigurations(tiConfig)
}

func (d *Driver) readConfigFile() ([]string, error) {
	file, err := os.Open("/home/legolas/git-file/Obstacle-Avoidance/mmwave_driver/go_driver/file.cfg")
	// file, err := os.Open("/home/cplid/file.cfg")
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var tiConfig []string

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] != '%' {
			tiConfig = append(tiConfig, strings.TrimRight(line, "\r\n"))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return tiConfig, nil
}

func (d *Driver) Configure(tiConfig []string) error {
	for _, i := range tiConfig {
		_, err := d.port.Write([]byte(i + "\n"))
		if err != nil {
			return fmt.Errorf("error writing to command port: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	go d.readConfigurationResp()
	log.Println("Configuration Done!")
	return nil
}

func (d *Driver) readConfigurationResp() {
	resp := ""

	for {
		buf := make([]byte, 1024)
		n, _ := d.port.Read(buf)
		resp += string(buf[:n])
		log.Println(resp)
	}
}

// func (d *Driver) processConfigurations(tiConfig []string) {
// 	foundRate := d.findFrameCfg(tiConfig)
// 	if !foundRate {
// 		fmt.Println("cfg parameters wrong")
// 		d.close()
// 		return
// 	}
// }

// func (d *Driver) findFrameCfg(tiConfig []string) bool {
// 	for _, i := range tiConfig {
// 		if d.isFrameCfg(i) {
// 			msPerFrame, err := d.extractMillisecondsPerFrame(i)
// 			if err != nil {
// 				fmt.Println("Error parsing milliseconds per frame:", err)
// 				return false
// 			}
// 			d.msPerFrame = msPerFrame
// 			fmt.Println("Found frameCfg, milliseconds per frame is ", d.msPerFrame)
// 			d.foundRate = true
// 			return true
// 		}
// 	}
// 	return false
// }

// func (d *Driver) isFrameCfg(line string) bool {
// 	splitWords := strings.Split(line, " ")
// 	return !d.foundRate && strings.Contains(splitWords[0], "frameCfg")
// }

// func (d *Driver) extractMillisecondsPerFrame(line string) (float64, error) {
// 	splitWords := strings.Split(line, " ")
// 	return strconv.ParseFloat(splitWords[5], 64)
// }
