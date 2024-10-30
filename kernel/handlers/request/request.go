package request

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

// STRUCTS SYSCALLS
type RequestProcessCreate struct {
	Pid            int    `json:"pid"`
	Pseudocodigo   string `json:"pseudocodigo"`
	TamanioMemoria int    `json:"tamanio_memoria"`
	Prioridad      int    `json:"prioridad"`
}

type RequestThreadCreate struct {
	Pid          int    `json:"pid"`
	Pseudocodigo string `json:"pseudocodigo"`
	Prioridad    int    `json:"prioridad"`
}

type RequestInterrupcion struct {
	Razon string `json:"razon"`
	Pid   int    `json:"pid"`
}

type RequestProcessExit struct {
	Pid int `json:"pid"`
	Tid int `json:"tid"`
}

type RequestThreadExit struct {
	Pid int `json:"pid"`
	Tid int `json:"tid"`
}

type RequestThreadJoin struct {
	Tid          int `json:"tid"`
	TidParametro int `json:"tidparametro"`
}

type RequestThreadCancel struct {
	TidAEliminar int `json:"tid"`
	Pid          int `json:"pid"`
}

type RequestMutex struct {
	Nombre string `json:"nombre"`
	Pid    int    `json:"pid"`
	Tid    int    `json:"tid"`
}

type RequestDumpMemory struct {
	Pid int `json:"pid"`
	Tid int `json:"tid"`
}

type RequestIO struct {
	Tid    int `json:"tid"`
	Tiempo int `json:"tiempo"`
}

func SolicitarProcesoMemoria(pid int, pseudocodigo string, tamanio int) (*http.Response, error) {
	request := RequestProcessCreate{
		Pid:            pid,
		Pseudocodigo:   pseudocodigo,
		TamanioMemoria: tamanio,
	}

	solicitudCodificada, err := commons.CodificarJSON(request)

	if err != nil {
		return nil, err
	}

	return cliente.Post(globals.KConfig.IpMemory, globals.KConfig.PortMemory, "process", solicitudCodificada), nil
}

func SolicitarDumpMemory(pid int, tid int) (*http.Response, error) {
	request := RequestDumpMemory{
		Pid: pid,
		Tid: tid,
	}
	requestCodificado, _ := commons.CodificarJSON(request)
	cliente.Post(globals.KConfig.IpMemory, globals.KConfig.PortMemory, "dump", requestCodificado)
	return nil, nil
}

func Dispatch(pcb commons.PCB) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(pcb)
	if err != nil {
		return nil, err
	}

	return cliente.Post(globals.KConfig.IpCpu, globals.KConfig.PortCpu, "dispatch", requestBody), err
}

func Interrupt(interruption string, pid int) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(RequestInterrupcion{Razon: interruption, Pid: pid})
	if err != nil {
		return nil, err
	}

	return cliente.Post(globals.KConfig.IpCpu, globals.KConfig.PortCpu, "interrupt", requestBody), err
}
