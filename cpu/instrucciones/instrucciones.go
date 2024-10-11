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
		general.NotifyKernel(respuesta, "--")
	case "IO":
		respuesta.Reason = "IO"
		general.NotifyKernel(respuesta, "--")
	case "PROCESS_CREATE":
		respuesta.Reason = "PROCESS_CREATE"
		general.NotifyKernel(respuesta, "/syscall/process_create")
	case "THREAD_CREATE":
		respuesta.Reason = "THREAD_CREATE"
		general.NotifyKernel(respuesta, "/syscall/thread_create")
	case "THREAD_JOIN":
		respuesta.Reason = "THREAD_JOIN"
		general.NotifyKernel(respuesta, "/syscall/thread_join")
	case "THREAD_CANCEL":
		respuesta.Reason = "THREAD_CANCEL"
		general.NotifyKernel(respuesta, "/syscall/thread_cancel")
	case "MUTEX_CREATE":
		respuesta.Reason = "MUTEX_CREATE"
		general.NotifyKernel(respuesta, "/syscall/mutex_create")
	case "MUTEX_LOCK":
		respuesta.Reason = "MUTEX_LOCK"
		general.NotifyKernel(respuesta, "/syscall/mutex_lock")
	case "MUTEX_UNLOCK":
		respuesta.Reason = "MUTEX_UNLOCK"
		general.NotifyKernel(respuesta, "/syscall/mutex_unlock")
	case "THREAD_EXIT":
		respuesta.Reason = "THREAD_EXIT"
		general.NotifyKernel(respuesta, "/syscall/thread_exit")
	case "PROCESS_EXIT":
		respuesta.Reason = "PROCESS_EXIT"
		general.NotifyKernel(respuesta, "/syscall/process_exit")
	}

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
