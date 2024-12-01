package requests

import (
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
)

type RequestContexto struct {
	Pid int `json:"pid"`
	Tid int `json:"tid"`
}

type ResponseContexto struct {
	Pid       int                  `json:"pid"`
	Tid       int                  `json:"tid"`
	Registros globals.ContextoHilo `json:"registros"`
	Base      int                  `json:"base"`
	Limite    int                  `json:"limite"`
}

type RequestActualizarContexto struct {
	Pid       int                  `json:"pid"`
	Tid       int                  `json:"tid"`
	Registros globals.ContextoHilo `json:"registros"`
}

type RequestObtenerInstruccion struct {
	Pid int    `json:"pid"`
	Tid int    `json:"tid"`
	PC  uint32 `json:"pc"`
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
	Byte byte `json:"byte"`
	Pid  int  `json:"pid"`
	Tid  int  `json:"tid"`
}

type RequestWriteMemory struct {
	Byte  byte   `json:"byte"`
	Pid   int    `json:"pid"`
	Tid   int    `json:"tid"`
	Datos []byte `json:"datos"`
}

type RequestDumpMemory struct {
	Pid int `json:"pid"`
	Tid int `json:"tid"`
}

type DumpMemoryFS struct {
	Pid       int    `json:"pid"`
	Tid       int    `json:"tid"`
	Tamanio   int    `json:"tamanio"`
	Contenido []byte `json:"contenido"`
}
