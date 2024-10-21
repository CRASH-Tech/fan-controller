package main

type Config struct {
	SerialPort   string     `yaml:"serialPort"`
	PortMap      []PortItem `yaml:"portMap"`
	CriticalTemp int        `yaml:"criticalTemp"`
}

type PortItem struct {
	ID   int   `yaml:"id"`
	Fans []int `yaml:"fans"`
}

type DeviceStats struct {
	Temp    int
	UtilGPU int
	UtilMem int
}
