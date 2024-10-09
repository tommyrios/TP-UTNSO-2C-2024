package instrucciones

import (
	"log"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/utils"
)

type InstruccionStruct struct {
	Partes            []string
	CodigoInstruccion string
	Operandos         []string
}

var Instruccion *InstruccionStruct

func Set() {
	registro := Instruccion.Operandos[0]
	valor := utils.ConvertirStringAEntero(Instruccion.Operandos[1])
	aplicarCambio(registro, valor)
}

func Sum() {
	registroDestino := utils.ValorRegistros(Instruccion.Operandos[0])
	registroOrigen := utils.ValorRegistros(Instruccion.Operandos[1])

	aplicarCambio(Instruccion.Operandos[0], registroDestino+registroOrigen)

}

func Sub() {
	registroDestino := utils.ValorRegistros(Instruccion.Operandos[0])
	registroOrigen := utils.ValorRegistros(Instruccion.Operandos[1])

	aplicarCambio(Instruccion.Operandos[0], registroDestino-registroOrigen)

}

func Jnz() bool {
	registro := utils.ValorRegistros(Instruccion.Operandos[0])

	if registro != 0 {
		aplicarCambio("PC", utils.ValorRegistros(Instruccion.Operandos[1]))
	}

	return registro != 0
}

func Log() {
	registro := globals.Instruccion.Operandos[0]
	valor := utils.ValorRegistros(registro)

	log.Printf("TID: %d - LOG - Registro: %s - Valor: %d", *globals.Pid, registro, valor)
}

func aplicarCambio(registro string, valor uint32) {
	registroAux := globals.Regis[registro]
	switch v := registroAux.(type) {
	case *uint32:
		*v = valor
	case *uint8:
		*v = uint8(valor)
	}
}
