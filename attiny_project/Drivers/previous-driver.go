package tamper

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/aler9/gomavlib/pkg/msg"
	"github.com/nextuav/thunderlink-go/pkg/dialects/umt"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

const (
	DriverName                = "tamper"
	keyTamperActionEnabled    = "TAMP_ACT_ENBL"
	keyTamperActionDevicePort = "TAMP_ACT_DEV_PORT"
	i2cAddress                = 04
	cmdReadState              = 0
	cmdResetState             = 40
	timeoutResetState         = time.Second
	valTampered               = 20
	readPeriod                = time.Second
	mavlinkMessagePeriod      = time.Second
)

type Driver struct {
	msgWrite chan<- msg.Message
	msgRead  <-chan msg.Message
	status   *status
	device   *i2c.Dev
}

func NewDriver(msgWrite chan<- msg.Message, msgRead <-chan msg.Message) (*Driver, error) {
	if _, err := host.Init(); err != nil {
		return nil, fmt.Errorf("%s: cannot init drivers: %s", DriverName, err)
	}

	devPort := os.Getenv(keyTamperActionDevicePort)
	log.Printf("%s: devicePort = %s", DriverName, devPort)

	b, err := i2creg.Open(devPort)
	if err != nil {
		return nil, fmt.Errorf("%s: cannot open device %s: %s", DriverName, devPort, err)
	}

	device := &i2c.Dev{Addr: i2cAddress, Bus: b}

	driver := &Driver{
		msgWrite: msgWrite,
		msgRead:  msgRead,
		status:   &status{hardwareState: umt.UMT_TAMPER_STATE_UNKNOWN},
		device:   device,
	}

	// If disabled reset the device and return an error
	if enabled, err := strconv.ParseBool(os.Getenv(keyTamperActionEnabled)); err != nil || !enabled {
		driver.resetDevice(10)
		return nil, fmt.Errorf("%s: disabled", DriverName)
	}

	return driver, nil
}

func (d *Driver) Start() {
	go d.readTamperState()
	go d.relayHardwareTamperStatus()
	go d.observeResetCommand()

	for {
		time.Sleep(time.Duration(math.MaxInt64))
	}
}

func (d *Driver) relayHardwareTamperStatus() {
	for {
		time.Sleep(mavlinkMessagePeriod)

		msg := &umt.MessageUmtHardwareTamperStatus{
			TimeUsec: uint64(time.Now().UnixMicro()),
			State:    d.status.hardwareState,
		}
		d.msgWrite <- msg
	}
}

func (d *Driver) readTamperState() {
	log.Printf("%s: reading tamper status", DriverName)

	writeBuffer := []byte{cmdReadState}
	readBuffer := make([]byte, 1)

	for {
		time.Sleep(readPeriod)

		if err := d.device.Tx(writeBuffer, readBuffer); err != nil {
			d.status.updateHardwareState(umt.UMT_TAMPER_STATE_TAMPERED)
		} else {
			if readBuffer[0] == valTampered {
				d.status.updateHardwareState(umt.UMT_TAMPER_STATE_TAMPERED)
			} else {
				d.status.updateHardwareState(umt.UMT_TAMPER_STATE_UNTAMPERED)
			}
		}

		log.Printf("%s: value = %d state = %d", DriverName, readBuffer[0], d.status.hardwareState)
	}
}

func (d *Driver) observeResetCommand() {
	for message := range d.msgRead {
		switch msg := message.(type) {
		case *umt.MessageUmtReqResetHardwareTamper:
			log.Printf("%s: received reset command: requestId = %d", DriverName, msg.RequestId)
			d.resetDevice(msg.RequestId)
		}
	}
}

func (d *Driver) resetDevice(requestId uint32) {
	time.Sleep(5 * time.Millisecond)
	ack := &umt.MessageUmtAckResetHardwareTamper{
		RequestId:     requestId,
		TimeoutMillis: uint32(timeoutResetState.Milliseconds()),
	}
	d.msgWrite <- ack

	if _, err := d.device.Write([]byte{cmdResetState}); err != nil {
		log.Printf("%s: cannot reset device: %s", DriverName, err)

		time.Sleep(5 * time.Millisecond)
		resp := &umt.MessageUmtRespResetHardwareTamper{
			RequestId:     requestId,
			Result:        umt.UMT_RESULT_FAILURE,
			FailureReason: "Cannot write to device",
		}
		d.msgWrite <- resp
		return
	}

	log.Printf("%s: reset success", DriverName)

	time.Sleep(5 * time.Millisecond)
	resp := &umt.MessageUmtRespResetHardwareTamper{
		RequestId: requestId,
		Result:    umt.UMT_RESULT_SUCCESS,
	}
	d.msgWrite <- resp
}


