package globals

import (
	"sync"

	"github.com/sisoputnfrba/tp-golang/utils/commons"
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

type InstruccionStruct struct {
	Partes            []string
	CodigoInstruccion string
	Operandos         []string
}

var CConfig *Config

var ColaNEW []Process

var Registros *commons.Registros

var Tid *int

var Pid *int

var Instruccion *InstruccionStruct

var Regis map[string]interface{}

func CargarRegistros() {
	Regis = map[string]interface{}{
		"PC": &Registros.PC,
		"AX": &Registros.AX,
		"BX": &Registros.BX,
		"CX": &Registros.CX,
		"DX": &Registros.DX,
		"EX": &Registros.EX,
		"FX": &Registros.FX,
		"GX": &Registros.GX,
		"HX": &Registros.HX,
	}
}
