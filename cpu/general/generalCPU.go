package general

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func ObtenerInstruction() (commons.GetRespuestaInstruccion, error) {
	requestBody, err := commons.CodificarJSON(commons.GetPedidoInstruccion{Pid: *globals.Pid, Tid: *globals.Tid, PC: globals.Registros.PC})

	if err != nil {
		return commons.GetRespuestaInstruccion{}, err
	}

	response, instruccion := cliente.Post2(globals.CConfig.IpMemory, globals.CConfig.PortMemory, "obtener_instruccion", requestBody)

	defer response.Body.Close()

	instruccionDecodificada := commons.GetRespuestaInstruccion{Instruccion: instruccion}

	//err = commons.DecodificarJSON(response.Body, &instruccionDecodificada)

	if err != nil {
		return commons.GetRespuestaInstruccion{}, err
	}

	return instruccionDecodificada, nil
}

func NotifyKernel(respuesta *commons.DespachoProceso, ruta string) (*http.Response, error) {
	// Codificar el contexto de ejecución en JSON
	requestBody, err := commons.CodificarJSON(respuesta)
	if err != nil {
		return nil, err
	}

	return cliente.Post(globals.CConfig.IpKernel, globals.CConfig.PortKernel, ruta, requestBody), nil
}

func NotifyMemory(tcb commons.TCB) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(tcb)
	if err != nil {
		return nil, err
	}

	return cliente.Post(globals.CConfig.IpMemory, globals.CConfig.PortMemory, "contexto actualizado", requestBody), nil
}
