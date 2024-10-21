// package main
package main

import (
	"flag"
	"fmt"
	"sort"
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
	for {
		getFanStats()
		getDeviceStats()

		setDevicesCurve()

		checkAllDevices()
		watchdog()
		time.Sleep(1 * time.Second)
	}
}

func setDevicesCurve() {
	for _, deviceConfig := range config.Devices {
		keys := make([]int, 0, len(deviceConfig.Curve))
		for k, _ := range deviceConfig.Curve {
			keys = append(keys, k)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(keys)))

		for _, k := range keys {
			if deviceStats[deviceConfig.ID].Temp >= k {
				setDeviceFans(deviceConfig.ID, deviceConfig.Curve[k])
				break
			}
		}
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

		deviceConfig, err := getDeviceConfig(deviceIndex)
		if err != nil {
			setEmergencyMode()
			panic(err)
		}

		var fans string
		for _, fanId := range deviceConfig.Fans {
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
	device, err := getDeviceConfig(id)
	if err != nil {
		setEmergencyMode()
		panic(err)
	}
	for _, fanId := range device.Fans {
		err := setFanSpeed(fanId, percent)
		if err != nil {
			setEmergencyMode()
			panic(err)
		}
	}

	return nil
}

func getDeviceConfig(id int) (Device, error) {
	for _, device := range config.Devices {
		if device.ID == id {
			return device, nil
		}
	}

	return Device{}, fmt.Errorf("cannot find device %d in config", id)
}

func checkAllDevices() {
	for deviceIndex, stats := range deviceStats {
		if stats.Temp >= config.CriticalTemp {
			fmt.Printf("Device %d is over critical temperature(%d>=%d)! Turn on Emergency mode!\n", deviceIndex, stats.Temp, config.CriticalTemp)
			setEmergencyMode()
		}
	}
}
