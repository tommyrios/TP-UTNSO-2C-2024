package mmu

import (
	"fmt"
	"log"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type MemoriaRequest struct {
	DireccionFisica uint32 `json:"direccion_fisica"`
	Valor           uint32 `json:"valor"`
}

type MemoriaResponse struct {
	Valor uint32 `json:"valor"`
	Error string `json:"error"`
}

func TraducirDireccion(base, desplazamiento, limite uint32) (uint32, error) {
	if desplazamiento >= limite {
		return 0, fmt.Errorf("SEGMENTATION FAULT")
	}
	return desplazamiento + base, nil
}

func LeerMemoria(direccionLogica uint32) (uint32, error) {
	direccionFisica, err := TraducirDireccion(globals.Registros.Base, direccionLogica, globals.Registros.Limite)
	if err != nil {
		return 0, err
	}

	request := MemoriaRequest{DireccionFisica: direccionFisica}
	datosEnvioMemoria, err := commons.CodificarJSON(request)
	if err != nil {
		return 0, err
	}
	respuesta := cliente.Post(globals.CConfig.IpMemory, globals.CConfig.PortMemory, "/memory/read", datosEnvioMemoria)
	defer respuesta.Body.Close()

	var respuestaMemoria MemoriaResponse
	err = commons.DecodificarJSON(respuesta.Body, &respuestaMemoria)
	if err != nil {
		return 0, fmt.Errorf("error al decodificar respuesta de lectura: %s", err)
	}

	log.Printf("## TID: %d - Acción: LEER - Dirección Física: %d", *globals.Tid, direccionFisica)
	return respuestaMemoria.Valor, nil
}

func EscribirMemoria(direccionLogica uint32, valor uint32) error {
	direccionFisica, err := TraducirDireccion(globals.Registros.Base, direccionLogica, globals.Registros.Limite)
	if err != nil {
		return err
	}

	request := MemoriaRequest{DireccionFisica: direccionFisica, Valor: valor}
	datosEnvioMemoria, err := commons.CodificarJSON(request)
	if err != nil {
		return err
	}
	respuesta := cliente.Post(globals.CConfig.IpMemory, globals.CConfig.PortMemory, "/memory/write", datosEnvioMemoria)
	defer respuesta.Body.Close()

	var respuestaMemoria MemoriaResponse
	err = commons.DecodificarJSON(respuesta.Body, &respuestaMemoria)
	if err != nil {
		return fmt.Errorf("error al decodificar respuesta de escritura: %s", err)
	}

	log.Printf("## TID: %d - Acción: ESCRIBIR - Dirección Física: %d", *globals.Tid, direccionFisica)
	return nil
}
