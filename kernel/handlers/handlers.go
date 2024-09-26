package handlers

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/handlers/request"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func HandleProcessCreate(w http.ResponseWriter, r *http.Request) {
	var request request.RequestProceso
	err := commons.DecodificarJSON(r.Body, &request)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	globals.CrearProceso(request.Pseudocodigo, request.TamanioMemoria, request.Prioridad)
}

func HandleThreadCreate(w http.ResponseWriter, r *http.Request) {
	var request request.RequestThreadCreate
	err := commons.DecodificarJSON(r.Body, &request)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	//globals.CrearHilo()
}

func HandleProcessExit(w http.ResponseWriter, r *http.Request) {}
func HandleThreadExit(w http.ResponseWriter, r *http.Request)  {}
func HandleMutexCreate(w http.ResponseWriter, r *http.Request) {}
func HandleMutexLock(w http.ResponseWriter, r *http.Request)   {}
func HandleMutexUnlock(w http.ResponseWriter, r *http.Request) {}
