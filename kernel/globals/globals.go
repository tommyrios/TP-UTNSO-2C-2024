package globals

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"net/http"
)

type Config struct {
	Port               int    `json:"port"`
	IpMemory           string `json:"ip_memory"`
	PortMemory         int    `json:"port_memory"`
	IpCpu              string `json:"ip_cpu"`
	PortCpu            int    `json:"port_cpu"`
	SchedulerAlgorithm string `json:"scheduler_algorithm"`
	Quantum            int    `json:"quantum"`
	LogLevel           string `json:"log_level"`
}

type RequestProceso struct {
	Pseudocodigo   string `json:"pseudocodigo"`
	TamanioMemoria int    `json:"tamanio_memoria"`
}

type Syscall struct {
	Nombre         string `json:"nombre"`
	Tiempo         int    `json:"tiempo"`
	Pseudocodigo   string `json:"archivo_pseudocodigo"`
	TamanioMemoria int    `json:"tamanio_memoria"`
	PrioridadTID0  int    `json:"prioridad"`
	TID            int    `json:"tid"`
	Recurso        string `json:"recurso"`
}

var KConfig *Config

func Syscalls(w http.ResponseWriter, r *http.Request) {

	var syscall Syscall

	if r.Body == nil {
		http.Error(w, "Cuerpo de solicitud vacío", http.StatusBadRequest)
		return
	}

	err := commons.DecodificarJSON(r.Body, &syscall)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	switch syscall.Nombre {
	case "PROCESS_CREATE":
		CrearProceso(syscall.Pseudocodigo, syscall.TamanioMemoria, syscall.PrioridadTID0)
	case "PROCESS_EXIT":

	case "THREAD_CREATE":

	case "THREAD_EXIT":

	case "THREAD_JOIN":

	case "THREAD_CANCEL":

	case "DUMP_MEMORY":

	case "IO":

	case "MUTEX_CREATE":

	case "MUTEX_LOCK":

	case "MUTEX_UNLOCK":

	default:
		http.Error(w, "Syscall no reconocida", http.StatusBadRequest)
	}
}
