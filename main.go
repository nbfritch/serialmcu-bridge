package main

import (
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
	portFlag = flag.String("port", "/dev/ttyACM0", "Serial port to read")
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

	amTemp, err := readCmd(port, AM2320_TEMP)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("AM2320 Temp: %s\n", amTemp)

	amHum, err := readCmd(port, AM2320_HUMIDITY)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("AM2320 Humidity: %s\n", amHum)

	dsTemp, err := readCmd(port, DS_TEMP)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DS Temp: %s\n", dsTemp)

	lux, err := readCmd(port, LUX_CMD)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("LUX: %s\n", lux)
}
