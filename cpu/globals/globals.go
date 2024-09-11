package globals

import (
	"sync"
)

type Config struct {
	IpMemory   string `json:"ip_memory"`
	PortMemory int    `json:"port_memory"`
	IpKernel   string `json:"ip_kernel"`
	PortKernel int    `json:"port_kernel"`
	Port       int    `json:"port"`
	LogLevel   string `json:"log_level"`
}

type PCB struct {
	Pid   int `json:"pid"`
	Tid   int `json:"tid"`
	Mutex sync.Mutex
}

type Process struct {
	Pid    int    `json:"pid"`
	Estado string `json:"estado"`
	PCB    PCB    `json:"pcb"`
}

var CConfig *Config

var ColaNEW []Process
