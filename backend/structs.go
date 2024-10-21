package main

type Config struct {
	SerialPort   string   `yaml:"serialPort"`
	Devices      []Device `yaml:"devices"`
	CriticalTemp int      `yaml:"criticalTemp"`
}

type Device struct {
	ID    int         `yaml:"id"`
	Fans  []int       `yaml:"fans"`
	Curve map[int]int `yaml:"curve"`
}

type DeviceStats struct {
	Temp    int
	UtilGPU int
	UtilMem int
}
