// package main
package main

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"go.bug.st/serial"
)

const (
	FAN_COUT = 4
)

var (
	config      Config
	port        serial.Port
	deviceStats map[int]DeviceStats
	fanStats    map[int]int
)

func init() {
	deviceStats = make(map[int]DeviceStats)
	fanStats = make(map[int]int)

	configPath := flag.String("c", "config.yaml", "Path to the YAML configuration file")
	flag.Parse()

	var err error
	config, err = readConfig(*configPath)
	if err != nil {
		panic(fmt.Errorf("cannot read config file: %s", err))
	}

	err = nvInit()
	if err != nil {
		panic(fmt.Errorf("cannot init nvml: %s", err))
	}

	mode := &serial.Mode{
		BaudRate: 115200,
	}

	port, err = serial.Open(config.SerialPort, mode)
	if err != nil {
		panic(fmt.Errorf("cannot open serial port: %s", err))
	}
}

func main() {
	setDeviceFans(0, 40)

	time.Sleep(1 * time.Second)
	for {
		getFanStats()
		getDeviceStats()

		checkAllDevices()
		watchdog()
		time.Sleep(1 * time.Second)
	}
}

func getDeviceStats() {
	deviceCount, err := getDeviceCount()
	if err != nil {
		setEmergencyMode()
		panic(fmt.Errorf("cannot get devices count: %s", err))
	}

	for deviceIndex := range deviceCount {
		temp, err := getDeviceGPUTemp(deviceIndex)
		if err != nil {
			setEmergencyMode()
			panic(fmt.Errorf("cannot get device temperature: %s", err))
		}

		utilGPU, utilMem, err := getDeviceUtil(deviceIndex)
		if err != nil {
			setEmergencyMode()
			panic(fmt.Errorf("cannot get device utilization: %s", err))
		}

		deviceStats[deviceIndex] = DeviceStats{Temp: temp, UtilGPU: utilGPU, UtilMem: utilMem}

		var fans string
		for _, fanId := range config.PortMap[deviceIndex].Fans {
			fans = fmt.Sprintf("%s FAN_%d: %d", fans, fanId, fanStats[fanId])

		}

		fmt.Printf("DEVICE: %d TEMP: %d GPU: %d%% MEM: %d%%%s\n", deviceIndex, temp, utilGPU, utilMem, fans)
	}
}

func getFanStats() {
	for fanId := range FAN_COUT {
		fanId = fanId + 1

		rpm, err := getFanSpeed(fanId)
		if err != nil {
			setEmergencyMode()
			panic(fmt.Errorf("cannot get fan speed: %s", err))
		}
		fanStats[fanId] = rpm
	}
}

func setDeviceFans(id, percent int) error {
	for _, portMap := range config.PortMap {
		if portMap.ID == id {
			for _, fanId := range portMap.Fans {
				setFanSpeed(fanId, percent)
			}
		}
	}

	return errors.New("cannot find device id")
}

// func printDeviceStats() {
// 	for deviceIndex, stats := range deviceStats {
// 		if stats.Temp >= config.CriticalTemp {
// 			fmt.Printf("DEVICE: %d TEMP: %d GPU: %d%% MEM: %d%%\n", deviceIndex, stats.Temp, stats.UtilGPU, stats.UtilMem)
// 		}
// 	}
// }

func checkAllDevices() {

	for deviceIndex, stats := range deviceStats {
		if stats.Temp >= config.CriticalTemp {
			fmt.Printf("Device %d is over critical temperature(%d>=%d)! Turn on Emergency mode!\n", deviceIndex, stats.Temp, config.CriticalTemp)
			setEmergencyMode()
		}
	}
}
