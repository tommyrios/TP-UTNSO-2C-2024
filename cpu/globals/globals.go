package globals

import (
	"github.com/sisoputnfrba/tp-golang/cpu/globals/requests"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
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

var CConfig *Config

var Syscall = make(chan bool)

type InstruccionStruct struct {
	CodOperacion string   `json:"cod_operacion"`
	Operandos    []string `json:"operandos"`
}

type InterrupcionStruct struct {
	Pid       int  `json:"pid"`
	Tid       int  `json:"tid"`
	Presencia bool `json:"presencia"`
}

func ValorRegistros(registroBuscado string, registros *commons.Registros) uint32 {
	if registroBuscado == "PC" {
		return registros.PC
	} else if registroBuscado == "AX" {
		return registros.AX
	} else if registroBuscado == "BX" {
		return registros.BX
	} else if registroBuscado == "CX" {
		return registros.CX
	} else if registroBuscado == "DX" {
		return registros.DX
	} else if registroBuscado == "EX" {
		return registros.EX
	} else if registroBuscado == "FX" {
		return registros.FX
	} else if registroBuscado == "GX" {
		return registros.GX
	} else if registroBuscado == "HX" {
		return registros.HX
	}

	return 999
}

func CambiarValorRegistros(registro string, valor uint32, registros *commons.Registros) {
	if registro == "PC" {
		registros.PC = valor
	} else if registro == "AX" {
		registros.AX = valor
	} else if registro == "BX" {
		registros.BX = valor
	} else if registro == "CX" {
		registros.CX = valor
	} else if registro == "DX" {
		registros.DX = valor
	} else if registro == "EX" {
		registros.EX = valor
	} else if registro == "FX" {
		registros.FX = valor
	} else if registro == "GX" {
		registros.GX = valor
	} else if registro == "HX" {
		registros.HX = valor
	}
}

func Mmu(desplazamiento int, base int, limite int) (int, int) {
	if desplazamiento < 0 || desplazamiento >= limite || desplazamiento+base >= limite {
		return -1, 1
	}

	return desplazamiento + base, 0
}

func DevolverPCB(pid int, tid int, razon string) {
	reqDispatch := requests.RequestDevolverPcb{Pid: pid, Tid: tid, Razon: razon}

	reqCodificada, _ := commons.CodificarJSON(reqDispatch)

	cliente.Post(CConfig.IpKernel, CConfig.PortKernel, "pcb", reqCodificada)

}
