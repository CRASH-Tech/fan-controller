package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func watchdog() error {
	_, err := port.Write([]byte("WATCHDOG\n"))
	if err != nil {
		return err
	}

	return nil
}

func setEmergencyMode() {
	for fanId := range FAN_COUT {
		fanId = fanId + 1

		setFanSpeed(fanId, 100)
	}
}

func setFanSpeed(id, percent int) error {
	fmt.Printf("SET FAN %d TO %d%%\n", id, percent)

	_, err := port.Write([]byte(fmt.Sprintf("SET %d %d\n", id, percent)))
	if err != nil {
		time.Sleep(500 * time.Millisecond)
		return err
	}

	time.Sleep(500 * time.Millisecond)
	return nil
}

func getFanSpeed(id int) (int, error) {
	//fmt.Printf("GET FAN %d SPEED\n", id)

	// _, err := port.Write([]byte(fmt.Sprintf("\n", id)))
	// if err != nil {
	// 	return -1, err
	// }

	_, err := port.Write([]byte(fmt.Sprintf("GET %d\n", id)))
	if err != nil {
		time.Sleep(500 * time.Millisecond)
		return -1, err
	}

	buff := make([]byte, 100)
	for {
		err := port.ResetInputBuffer()
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			return -1, err
		}

		err = port.ResetOutputBuffer()
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			return -1, err
		}

		port.SetReadTimeout(5 * time.Second)
		n, err := port.Read(buff)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			return -1, err
		}
		if n == 0 {
			break
		}

		data := strings.TrimSpace(string(buff[:n]))
		rpm, err := strconv.Atoi(data)

		time.Sleep(500 * time.Millisecond)
		return rpm, err
	}

	time.Sleep(500 * time.Millisecond)
	return -1, err
}
