package general

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func GetInstruction() (*http.Response, error) {
	//Hay que crear el json de la instruccion con el PID y el PC de la PCB

	var requestBody bytes.Buffer

	err := commons.CodificarJSON(&requestBody, commons.PedidoInstruccion{Pid: *globals.Pid, PC: globals.Registros.PC})

	if err != nil {
		return nil, err
	}

	response := cliente.Post(globals.CConfig.IpMemory, globals.CConfig.PortMemory, "instruction", requestBody.Bytes())

	// Verificar si ocurrió algún error en la respuesta
	//OPCIONAL
	if response.StatusCode >= 400 {
		log.Printf("Error en la solicitud POST: %d %s", response.StatusCode, response.Status)
		return response, fmt.Errorf("error en la solicitud POST: %d", response.StatusCode)
	}

	return response, nil
}
