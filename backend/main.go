// package main
package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.bug.st/serial"
)

var (
	port serial.Port
)

func init() {
	err := nvInit()
	if err != nil {
		panic(err)
	}
	mode := &serial.Mode{
		BaudRate: 115200,
	}

	port, err = serial.Open("/dev/ttyUSB0", mode)
	if err != nil {
		panic(err)
	}

	// port = aport
}

func main() {
	deviceCount, err := getDeviceCount()
	if err != nil {
		panic(err)
	}
	fmt.Printf("DEVICE COUNT: %d\n", deviceCount)

	// for {
	// 	for deviceIndex := range deviceCount {
	// 		//fmt.Printf("DEVICE INDEX: %d\n", deviceIndex)

	// 		temp, err := getDeviceGPUTemp(deviceIndex)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		utilGPU, utilMem, err := getDeviceUtil(deviceIndex)
	// 		if err != nil {
	// 			panic(err)
	// 		}

	// 		fmt.Printf("DEVICE: %d TEMP: %d GPU: %d%% MEM: %d%%\n", deviceIndex, temp, utilGPU, utilMem)
	// 	}
	// 	time.Sleep(1 * time.Second)
	// }
	// speeds := make(map[int]int)

	// for i := 1; i <= 4; i++ {
	// 	rpm, err := getFanSpeed(i)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	speeds[i] = rpm
	// 	time.Sleep(2 * time.Second)
	// }

	// for fan, rpm := range speeds {
	// 	fmt.Printf("FAN %d: %d\n", fan, rpm)
	// }

	// var err error
	// err = setFanSpeed(1, 30)
	// if err != nil {
	// 	panic(err)
	// }
	// err = setFanSpeed(2, 30)
	// if err != nil {
	// 	panic(err)
	// }

	// for {
	// 	err := watchdog()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	time.Sleep(10 * time.Second)
	// }

	// rpm, err := getFanSpeed(1)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(rpm)

}

func watchdog() error {
	_, err := port.Write([]byte("WATCHDOG\n"))
	if err != nil {
		return err
	}

	return nil
}

func setFanSpeed(id, percent int) error {
	//fmt.Printf("SET FAN %d TO %d\n", id, percent)

	_, err := port.Write([]byte(fmt.Sprintf("SET %d %d\n", id, percent)))
	if err != nil {
		return err
	}

	return nil
}

func getFanSpeed(id int) (int, error) {
	//fmt.Printf("GET FAN %d SPEED\n", id)
	_, err := port.Write([]byte(fmt.Sprintf("GET %d\n", id)))
	if err != nil {
		return -1, err
	}

	buff := make([]byte, 100)
	for {
		err := port.ResetInputBuffer()
		if err != nil {
			return -1, err
		}

		err = port.ResetOutputBuffer()
		if err != nil {
			return -1, err
		}

		port.SetReadTimeout(5 * time.Second)
		n, err := port.Read(buff)
		if err != nil {
			return -1, err
		}
		if n == 0 {
			break
		}

		//fmt.Printf("%v", string(buff[:n]))
		data := strings.TrimSpace(string(buff[:n]))
		rpm, err := strconv.Atoi(data)
		return rpm, err
	}

	return -1, err
}
