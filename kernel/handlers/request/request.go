package request

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

// STRUCTS SYSCALLS
type RequestProceso struct {
	Pseudocodigo   string `json:"pseudocodigo"`
	TamanioMemoria int    `json:"tamanio_memoria"`
	Prioridad      int    `json:"prioridad"`
}

type RequestThreadCreate struct {
	Pseudocodigo string `json:"pseudocodigo"`
	Prioridad    int    `json:"prioridad"`
}

func SolicitarProcesoMemoria(pseudocodigo string, tamanio int) (*http.Response, error) {
	request := RequestProceso{
		Pseudocodigo:   pseudocodigo,
		TamanioMemoria: tamanio,
		// VER SI HAY QUE PASAR LA PRIORIDAD O NO HACE FALTA
	}

	solicitudCodificada, err := commons.CodificarJSON(request)

	if err != nil {
		return nil, err
	}

	return cliente.Post(globals.KConfig.IpMemory, globals.KConfig.PortMemory, "process", solicitudCodificada), nil
}
