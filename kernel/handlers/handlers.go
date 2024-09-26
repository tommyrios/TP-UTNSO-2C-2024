package handlers

import (
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/request"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func HandleProcessCreate(w http.ResponseWriter, r *http.Request) {
	var proceso request.RequestProceso
	err := commons.DecodificarJSON(r.Body, &proceso)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	globals.CrearProceso(proceso.Pseudocodigo, proceso.TamanioMemoria, proceso.Prioridad)
}

func HandleThreadCreate(w http.ResponseWriter, r *http.Request) {
	var hilo request.RequestThreadCreate
	err := commons.DecodificarJSON(r.Body, &hilo)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	globals.CrearHilo(hilo.Pid, hilo.Prioridad, hilo.Pseudocodigo)
}

func HandleProcessExit(w http.ResponseWriter, r *http.Request) {

}

func HandleThreadExit(w http.ResponseWriter, r *http.Request)  {}
func HandleMutexCreate(w http.ResponseWriter, r *http.Request) {}
func HandleMutexLock(w http.ResponseWriter, r *http.Request)   {}
func HandleMutexUnlock(w http.ResponseWriter, r *http.Request) {}
