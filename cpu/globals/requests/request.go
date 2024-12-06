package requests

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type RequestDispatch struct {
	Pid       int    `json:"pid"`
	Tid       int    `json:"tid"`
	Quantum   int    `json:"quantum"`
	Scheduler string `json:"scheduler"`
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
	Pid            int    `json:"pid,omitempty"`
	Tid            int    `json:"tid,omitempty"`
	PseudoCodigo   string `json:"pseudocodigo,omitempty"`
	TamanioMemoria int    `json:"tamanio_memoria,omitempty"`
	Prioridad      int    `json:"prioridad,omitempty"`
	Tiempo         int    `json:"tiempo,omitempty"`
	TidParametro   int    `json:"tid_parametro,omitempty"`
	TidAEliminar   int    `json:"tidAEliminar,omitempty"`
	NombreMutex    string `json:"nombre_mutex,omitempty"`
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

type RequestMutex struct {
	Nombre string `json:"nombre_mutex"`
	Pid    int    `json:"pid"`
	Tid    int    `json:"tid"`
}
