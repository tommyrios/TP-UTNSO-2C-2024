package instrucciones

import (
	"log"

	"github.com/sisoputnfrba/tp-golang/cpu/general"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/utils"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

var Instruccion *globals.InstruccionStruct

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

func HandleSyscall(respuesta *commons.DespachoProceso) {
	respuesta.Pcb.Registros = *globals.Registros
	respuesta.Pcb.ProgramCounter = int(globals.Registros.PC)

	switch globals.Instruccion.CodigoInstruccion {
	case "DUMP_MEMORY":
		respuesta.Reason = "DUMP_MEMORY"
	case "IO":
		respuesta.Reason = "IO"
	case "PROCESS_CREATE":
		respuesta.Reason = "PROCESS_CREATE"
	case "THREAD_CREATE":
		respuesta.Reason = "THREAD_CREATE"
	case "THREAD_JOIN":
		respuesta.Reason = "THREAD_JOIN"
	case "THREAD_CANCEL":
		respuesta.Reason = "THREAD_CANCEL"
	case "MUTEX_CREATE":
		respuesta.Reason = "MUTEX_CREATE"
	case "MUTEX_LOCK":
		respuesta.Reason = "MUTEX_LOCK"
	case "MUTEX_UNLOCK":
		respuesta.Reason = "MUTEX_UNLOCK"
	case "THREAD_EXIT":
		respuesta.Reason = "THREAD_EXIT"
	case "PROCESS_EXIT":
		respuesta.Reason = "PROCESS_EXIT"
	}
	general.NotifyKernel(respuesta)
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
