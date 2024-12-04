package request

// STRUCTS SYSCALLS
type RequestProcessCreate struct {
	Pid            int    `json:"pid"`
	Tid            int    `json:"tid"`
	Pseudocodigo   string `json:"pseudocodigo"`
	TamanioMemoria int    `json:"tamanio_memoria"`
	Prioridad      int    `json:"prioridad"`
}

type RequestProcessCreateMemoria struct {
	Pid            int    `json:"pid"`
	Pseudocodigo   string `json:"pseudocodigo"`
	TamanioMemoria int    `json:"tamanio_memoria"`
}

type RequestProcessExit struct {
	Pid int `json:"pid"`
	Tid int `json:"tid"`
}

type RequestThreadCreate struct {
	Pid          int    `json:"pid"`
	Tid          int    `json:"tid"`
	Pseudocodigo string `json:"pseudocodigo"`
	Prioridad    int    `json:"prioridad"`
}

type RequestThreadJoin struct {
	Pid          int `json:"pid"`
	Tid          int `json:"tid"`
	TidParametro int `json:"tid_parametro"`
}

type RequestThreadCancel struct {
	Pid          int `json:"pid"`
	Tid          int `json:"tid"`
	TidAEliminar int `json:"tidAEliminar"`
}

type RequestThreadExit struct {
	Pid int `json:"pid"`
	Tid int `json:"tid"`
}

type RequestMutex struct {
	Nombre string `json:"nombre_mutex"`
	Pid    int    `json:"pid"`
	Tid    int    `json:"tid"`
}

type RequestDumpMemory struct {
	Pid int `json:"pid"`
	Tid int `json:"tid"`
}

type RequestIO struct {
	Pid    int `json:"pid"`
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
	Pid int `json:"pid"`
	Tid int `json:"tid"`
}

type RequestCrearHilo struct {
	Pid          int    `json:"pid"`
	Tid          int    `json:"tid"`
	Pseudocodigo string `json:"pseudocodigo"`
}

type RequestInterrupcion struct {
	Razon string `json:"razon"`
	Pid   int    `json:"pid"`
}

type RequestDevolucionPCB struct {
	Pid   int    `json:"pid"`
	Tid   int    `json:"tid"`
	Razon string `json:"razon"`
}
