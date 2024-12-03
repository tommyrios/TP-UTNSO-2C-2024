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
	Pid       int               `json:"pid"`
	Tid       int               `json:"tid"`
	Registros commons.Registros `json:"registros"`
	Base      int               `json:"base"`
	Limite    int               `json:"limite"`
}

type RequestInstruccion struct {
	Pid int `json:"pid"`
	Tid int `json:"tid"`
	PC  int `json:"pc"`
}
