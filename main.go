package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
)

const (
	AM2320_HUMIDITY = "am2320_humidity"
	AM2320_TEMP     = "am2320_temp"
	DS_TEMP         = "ds18x20_temp"
	LUX_CMD         = "lux"
)

var (
	portFlag     = flag.String("port", "/dev/ttyACM0", "Serial port to read")
	endpointFlag = flag.String("endpoint", "http://localhost:3000/readings", "Endpoint to PUT data")
)

func readCmd(port serial.Port, cmd string) (string, error) {
	n, err := port.Write([]byte(fmt.Sprintf("%s\r", cmd)))
	if err != nil {
		return "", err
	}

	time.Sleep(time.Second)

	buffer := make([]byte, 12)
	n, err = port.Read(buffer)
	if err != nil {
		return "", err
	}

	if n == 0 {
		return "err", nil
	}

	return string(buffer), nil
}

func readCmdTimeout(port serial.Port, cmd string, seconds uint) (string, error) {
	c1 := make(chan string, 1)
	go func() {
		res, err := readCmd(port, cmd)
		if err != nil {
			c1 <- ""
		} else {
			c1 <- res
		}
	}()

	select {
	case res := <-c1:
		return res, nil
	case <-time.After(10 * time.Second):
		return "", errors.New("Timeout")
	}
}

const CMD_TIMEOUT_SEC = 5

func main() {
	flag.Parse()

	mode := &serial.Mode{
		Parity:   serial.EvenParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(*portFlag, mode)
	if err != nil {
		log.Fatal(err)
	}

	amTemp, err := readCmdTimeout(port, AM2320_TEMP, CMD_TIMEOUT_SEC)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("AM2320 Temp: %s\n", amTemp)

	amHum, err := readCmdTimeout(port, AM2320_HUMIDITY, CMD_TIMEOUT_SEC)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("AM2320 Humidity: %s\n", amHum)

	dsTemp, err := readCmdTimeout(port, DS_TEMP, CMD_TIMEOUT_SEC)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DS Temp: %s\n", dsTemp)

	lux, err := readCmdTimeout(port, LUX_CMD, CMD_TIMEOUT_SEC)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("LUX: %s\n", lux)
}
