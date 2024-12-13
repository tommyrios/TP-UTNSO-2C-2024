package requests

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

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

type RequestActualizarContexto struct {
	Pid       int                `json:"pid"`
	Tid       int                `json:"tid"`
	Registros *commons.Registros `json:"registros"`
}

type RequestObtenerInstruccion struct {
	Pid int    `json:"pid"`
	Tid int    `json:"tid"`
	PC  uint32 `json:"pc"`
}

type ResponseObtenerInstruccion struct {
	Instruccion string `json:"instruccion"`
}

type RequestProcesoMemoria struct {
	Pid            int `json:"pid"`
	TamanioMemoria int `json:"tamanio_memoria"`
}

type RequestFinalizarProceso struct {
	Pid int `json:"pid"`
}

type RequestCrearHilo struct {
	Pid          int    `json:"pid"`
	Tid          int    `json:"tid"`
	Pseudocodigo string `json:"pseudocodigo"`
}

type RequestFinalizarHilo struct {
	Pid int `json:"pid"`
	Tid int `json:"tid"`
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

type RequestDumpMemory struct {
	Pid int `json:"pid"`
	Tid int `json:"tid"`
}

type DumpMemoryFS struct {
	Pid       uint32 `json:"pid"`
	Tid       uint32 `json:"tid"`
	Tamanio   int    `json:"tamanio"`
	Contenido []byte `json:"contenido"`
}
