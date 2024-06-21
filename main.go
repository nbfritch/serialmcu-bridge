package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.bug.st/serial"
)

const (
	AM2320_HUMIDITY = "am2320_humidity"
	AM2320_TEMP     = "am2320_temp"
	DS_TEMP         = "ds18x20_temp"
	LUX_CMD         = "lux"
	TEMPERATURE     = 1
	HUMIDITY        = 2
	LUX             = 3
	INVALID         = "INVALID"
)

var (
	portFlag     = flag.String("port", "/dev/ttyACM0", "Serial port to read")
	endpointFlag = flag.String("endpoint", "http://localhost:3000/readings", "Endpoint to PUT data")
	sendFlag     = flag.Bool("send", false, "Whether to submit readings")
	nameFlag     = flag.String("name", INVALID, "Name of sensor")
	verbose      = flag.Bool("verbose", false, "Verbose logging")
)

func validateReading(reading float64, readingType uint) bool {
	if readingType == TEMPERATURE {
		return reading > 0 && reading < 120
	}

	if readingType == HUMIDITY {
		return reading > 0 && reading < 100
	}

	if readingType == LUX {
		return reading > 0 && reading < 65336
	}

	return false
}

type CreateReadingReq struct {
	SensorName   string  `json:"sensorName"`
	ReadingType  int     `json:"reading_type"`
	ReadingValue float64 `json:"reading_value"`
}

func submitReading(reading float64, name string) error {
	url := "https://httpbin.org/put"
	reqData := CreateReadingReq{
		SensorName:   name,
		ReadingType:  TEMPERATURE,
		ReadingValue: reading,
	}
	data, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return errors.New("Create failed")
	}

	return nil
}

func readCmd(port serial.Port, cmd string) (string, error) {
	if *verbose {
		fmt.Printf("readCmd: %s\n", cmd)
	}
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

func readCmdTimeout(port serial.Port, cmd string, seconds int) (string, error) {
	if *verbose {
		fmt.Printf("readCmdTimeout(%d): %s\n", seconds, cmd)
	}
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
	case <-time.After(time.Second * time.Duration(seconds)):
		if *verbose {
			fmt.Printf("Timeout: %s\n", cmd)
		}
		return "", errors.New("Timeout")
	}
}

const CMD_TIMEOUT_SEC = 5

func main() {
	flag.Parse()
	name := *nameFlag

	if *sendFlag && name == INVALID {
		fmt.Println("Name is required when sending")
		return
	}

	mode := &serial.Mode{
		Parity:   serial.EvenParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	if _, err := os.Stat(*portFlag); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("%s does not exist\n", *portFlag)
		return
	}

	port, err := serial.Open(*portFlag, mode)
	if err != nil {
		log.Fatal(err)
	}

	amTempStr, amTempErr := readCmdTimeout(port, AM2320_TEMP, CMD_TIMEOUT_SEC)
	amTemp := 0.0
	if amTempErr != nil {
		amTemp, amTempErr = strconv.ParseFloat(amTempStr, 32)
	}

	if amTemp > 0 {
		fmt.Printf("AM2320 Temp: %s\n", amTempStr)
	}

	amHumStr, amHumErr := readCmdTimeout(port, AM2320_HUMIDITY, CMD_TIMEOUT_SEC)
	amHum := 0.0
	if amHumErr != nil {
		amHum, amHumErr = strconv.ParseFloat(amHumStr, 32)
	}

	if amHum > 0 {
		fmt.Printf("AM2320 Humidity: %s\n", amHumStr)
	}

	dsTempStr, dsTempErr := readCmdTimeout(port, DS_TEMP, CMD_TIMEOUT_SEC)
	dsTemp := 0.0
	if dsTempErr != nil {
		log.Fatal(err)
		dsTemp, dsTempErr = strconv.ParseFloat(dsTempStr, 32)
	}

	if dsTemp > 0 {
		fmt.Printf("DS Temp: %s\n", dsTempStr)
	}

	luxStr, luxErr := readCmdTimeout(port, LUX_CMD, CMD_TIMEOUT_SEC)
	lux := 0.0
	if err != nil {
		lux, luxErr = strconv.ParseFloat(luxStr, 32)
	}

	if lux > 0 {
		fmt.Printf("LUX: %f\n", lux)
	}

	if *sendFlag {
		temp := 0.0
		tempCount := 0
		if amTempErr == nil && validateReading(amTemp, TEMPERATURE) {
			tempCount = tempCount + 1
			temp = temp + amTemp
		}

		if dsTempErr != nil && validateReading(dsTemp, TEMPERATURE) {
			tempCount = tempCount + 1
			temp = temp + dsTemp
		}

		if tempCount > 0 {
			temp = temp / float64(tempCount)
			submitReading(temp, name)
		}

		if amHumErr != nil && validateReading(amHum, HUMIDITY) {
			submitReading(amHum, name)
		}

		if luxErr != nil && validateReading(lux, LUX) {
			submitReading(lux, name)
		}
	}
}
