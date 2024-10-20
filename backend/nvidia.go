package main

import (
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

func nvInit() error {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		return fmt.Errorf("unable to initialize NVML: %v", nvml.ErrorString(ret))
	}

	// defer func() {
	// 	ret := nvml.Shutdown()
	// 	if ret != nvml.SUCCESS {
	// 		panic(fmt.Sprintf("unable to shutdown NVML: %v", nvml.ErrorString(ret)))
	// 	}
	// }()

	return nil
}

func getDeviceCount() (int, error) {
	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		return -1, fmt.Errorf("unable to get device count: %v", nvml.ErrorString(ret))
	}

	return count, nil
}

func getDeviceGPUTemp(id int) (int, error) {
	device, ret := nvml.DeviceGetHandleByIndex(id)
	if ret != nvml.SUCCESS {
		return -1, fmt.Errorf("unable to get device at index %d: %v", id, nvml.ErrorString(ret))
	}

	temp, ret := device.GetTemperature(nvml.TEMPERATURE_GPU)
	if ret != nvml.SUCCESS {
		return -1, fmt.Errorf("unable to get device temperature %d: %v", id, nvml.ErrorString(ret))
	}

	return int(temp), nil
}

func getDeviceUtil(id int) (int, int, error) {
	device, ret := nvml.DeviceGetHandleByIndex(id)
	if ret != nvml.SUCCESS {
		return -1, -1, fmt.Errorf("unable to get device at index %d: %v", id, nvml.ErrorString(ret))
	}

	util, ret := device.GetUtilizationRates()
	if ret != nvml.SUCCESS {
		return -1, -1, fmt.Errorf("unable to get device utilization %d: %v", id, nvml.ErrorString(ret))
	}

	return int(util.Gpu), int(util.Memory), nil
}
