package requests

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type RequestDispatch struct {
	Pid int `json:"pid"`
	Tid int `json:"tid"`
}

type RequestContexto struct {
	Pid int `json:"pid"`
	Tid int `json:"tid"`
}

type ResponseContexto struct {
	Pid       int                `json:"pid"`
	Tid       int                `json:"tid"`
	Registros *commons.Registros `json:"registros"`
	Base      int                `json:"base"`
	Limite    int                `json:"limite"`
}

type RequestInstruccion struct {
	Pid int `json:"pid"`
	Tid int `json:"tid"`
	PC  int `json:"pc"`
}

type RequestActualizarRegistros struct {
	Pid       int                `json:"pid"`
	Tid       int                `json:"tid"`
	Registros *commons.Registros `json:"registros"`
}

type RequestSyscall struct {
	Pid            int    `json:"pid"`
	Tid            int    `json:"tid"`
	PseudoCodigo   string `json:"pseudocodigo"`
	TamanioMemoria int    `json:"tamanio_memoria"`
	Prioridad      int    `json:"prioridad"`
	Tiempo         int    `json:"tiempo"`
	TidParametro   int    `json:"tid_parametro"`
	TidAEliminar   int    `json:"tidAEliminar"`
	NombreMutex    string `json:"nombre_mutex"`
}

type RequestReadMemory struct {
	Direccion int `json:"direccion"`
	Pid       int `json:"pid"`
	Tid       int `json:"tid"`
}

type RequestWriteMemory struct {
	Direccion int    `json:"direccion"`
	Pid       int    `json:"pid"`
	Tid       int    `json:"tid"`
	Datos     []byte `json:"datos"`
}

type RequestDevolverPcb struct {
	Pid   int    `json:"pid"`
	Tid   int    `json:"tid"`
	Razon string `json:"razon"`
}

type RequestInterrupcion struct {
	Pid   int    `json:"pid"`
	Tid   int    `json:"tid"`
	Razon string `json:"razon"`
}
