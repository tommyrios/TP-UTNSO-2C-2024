package interrupciones

import (
	"log"
	"net/http"
	"sync"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type Interrupcion struct {
	Mutex  sync.Mutex
	Status bool
	Reason string
	Tid    int
}

var InterrupcionActual *Interrupcion

func RecibirInterrupcion(w http.ResponseWriter, r *http.Request) {
	var req globals.InterrupcionRecibida

	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Verificar si la interrupción es para el proceso actual
	if req.Tid == *globals.Tid {
		ActivarInterrupcion(true, req.Reason, req.Tid)
		log.Printf("Interrupción recibida para el TID: %d con motivo: %s", req.Tid, req.Reason)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Interrupción recibida"))
	if err != nil {
		log.Printf("Error al escribir la respuesta: %v", err)
	}
}

func ActivarInterrupcion(status bool, reason string, tid int) {
	InterrupcionActual.Mutex.Lock()
	defer InterrupcionActual.Mutex.Unlock()

	// Configurar los detalles de la interrupción
	InterrupcionActual.Status = status
	InterrupcionActual.Reason = reason
	InterrupcionActual.Tid = tid
}

func ObtenerYResetearInterrupcion() (bool, string, int) {
	InterrupcionActual.Mutex.Lock()
	defer InterrupcionActual.Mutex.Unlock()

	status := InterrupcionActual.Status
	reason := InterrupcionActual.Reason
	tid := InterrupcionActual.Tid

	// Resetear el estado
	InterrupcionActual.Status = false
	InterrupcionActual.Reason = ""
	InterrupcionActual.Tid = 0

	return status, reason, tid
}
