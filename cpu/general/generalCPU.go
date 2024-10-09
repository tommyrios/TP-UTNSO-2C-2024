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

func NotifyKernel(respuesta *commons.DespachoProceso) (*http.Response, error) {
	// Codificar el contexto de ejecución en JSON
	requestBody, err := commons.CodificarJSON(respuesta)
	if err != nil {
		return nil, err
	}

	return cliente.Post(globals.CConfig.IpKernel, globals.CConfig.PortKernel, "syscall", requestBody), nil
}
