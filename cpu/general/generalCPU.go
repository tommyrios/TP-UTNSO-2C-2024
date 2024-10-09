package general

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func ObtenerInstruction() (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(commons.GetPedidoInstruccion{Pid: *globals.Pid, PC: globals.Registros.PC})
	if err != nil {
		return nil, err
	}
	return cliente.Post(globals.CConfig.IpMemory, globals.CConfig.PortMemory, "process", requestBody), nil
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
	NotifyKernel(respuesta)
}

func NotifyKernel(respuesta *commons.DespachoProceso) (*http.Response, error) {
	// Codificar el contexto de ejecución en JSON
	requestBody, err := commons.CodificarJSON(respuesta)
	if err != nil {
		return nil, err
	}

	return cliente.Post(globals.CConfig.IpKernel, globals.CConfig.PortKernel, "syscall", requestBody), nil
}
