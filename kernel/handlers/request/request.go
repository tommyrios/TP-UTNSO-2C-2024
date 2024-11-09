package request

import "github.com/sisoputnfrba/tp-golang/utils/commons"

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
	PidParametro int `json:"pidparametro"`
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

type RequestFinalizarProceso struct {
	Pid int `json:"pid"`
}

type RequestFinalizarHilo struct {
	Pid int `json:"pid"`
	Tid int `json:"tid"`
}

type RequestDispatcher struct {
	PCB *commons.PCB `json:"pcb"`
	Tid int          `json:"tid"`
}
