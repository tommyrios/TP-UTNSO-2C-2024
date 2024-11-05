package request

import "github.com/sisoputnfrba/tp-golang/utils/commons"

type RequestContexto struct {
	Pid int `json:"pid"`
	Tid int `json:"tid"`
}

type ResponseContexto struct {
	Pid       int               `json:"pid"`
	Tid       int               `json:"tid"`
	Registros commons.Registros `json:"registros"`
	Base      int               `json:"base"`
	Limite    int               `json:"limite"`
}

type RequestActualizarContexto struct {
	Pid       int               `json:"pid"`
	Tid       int               `json:"tid"`
	Registros commons.Registros `json:"registros"`
}

type RequestObtenerInstruccion struct {
	PC  uint32 `json:"pc"`
	Pid int    `json:"pid"`
	Tid int    `json:"tid"`
}

type RequestMemory struct {
	Byte  byte   `json:"byte"`
	Datos []byte `json:"datos"`
	Pid   int    `json:"pid"`
}

type RequestProcesoMemoria struct {
	Pid            int `json:"pid"`
	TamanioMemoria int `json:"tamanio_memoria"`
}
