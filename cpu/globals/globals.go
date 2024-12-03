package globals

import (
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

type InstruccionStruct struct {
	CodOperacion string   `json:"cod_operacion"`
	Operandos    []string `json:"operandos"`
}

type InterrupcionRecibida struct {
	Tid    int    `json:"tid"`
	Reason string `json:"reason"`
}

var CConfig *Config

var Registros *commons.Registros

var Regis map[string]interface{}

func CargarRegistros() {
	Regis = map[string]interface{}{
		"PC":     &Registros.PC,
		"AX":     &Registros.AX,
		"BX":     &Registros.BX,
		"CX":     &Registros.CX,
		"DX":     &Registros.DX,
		"EX":     &Registros.EX,
		"FX":     &Registros.FX,
		"GX":     &Registros.GX,
		"HX":     &Registros.HX,
		"Base":   &Registros.Base,
		"Limite": &Registros.Limite,
	}
}
