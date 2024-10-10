package handlers

import (
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/request"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func HandleProcessCreate(w http.ResponseWriter, r *http.Request) {
	var proceso request.RequestProcessCreate
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
	var request request.RequestProcessExit
	err := commons.DecodificarJSON(r.Body, &request)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}
	if request.Tid == 0 {
		globals.FinalizarProceso(request.Pid)
	} else {
		http.Error(w, "La finalizacion de un proceso solo puede ser solicitada por el TID 0", http.StatusBadRequest)
	}

}

// THREAD_EXIT Finaliza el hilo que la invoca (el tid que se manda es del propio hilo)

func HandleThreadExit(w http.ResponseWriter, r *http.Request) {
	var hilo request.RequestThreadExit
	err := commons.DecodificarJSON(r.Body, &hilo)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	globals.FinalizarHilo(hilo.Pid, hilo.Tid)
}

func HandleThreadJoin(w http.ResponseWriter, r *http.Request) {
	var hilo request.RequestThreadJoin
	err := commons.DecodificarJSON(r.Body, &hilo)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Bloquear hilo (tid: request.tid) para darle lugar a que ejecute el hilo (tid: request.tidParametro) y luego desbloquearlo

}

// THREAD_CANCEL Finaliza el hilo cuyo tid se pasa por parámetro (desde otro hilo)

func HandleThreadCancel(w http.ResponseWriter, r *http.Request) {
	var request request.RequestThreadCancel
	err := commons.DecodificarJSON(r.Body, &request)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	globals.FinalizarHilo(request.Pid, request.TidAEliminar)
}

func HandleMutexCreate(w http.ResponseWriter, r *http.Request) {
	var mutex request.RequestMutex

	err := commons.DecodificarJSON(r.Body, &mutex)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	globals.CrearMutex(mutex.Nombre, mutex.Pid)
}

func HandleMutexLock(w http.ResponseWriter, r *http.Request) {
	var mutex request.RequestMutex

	err := commons.DecodificarJSON(r.Body, &mutex)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	globals.BloquearMutex(mutex.Nombre, mutex.Pid, mutex.Tid)
}

func HandleMutexUnlock(w http.ResponseWriter, r *http.Request) {
	var mutex request.RequestMutex

	err := commons.DecodificarJSON(r.Body, &mutex)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	globals.DesbloquearMutex(mutex.Nombre, mutex.Pid, mutex.Tid)
}

func HandleDumpMemory(w http.ResponseWriter, r *http.Request) {
	var req request.RequestDumpMemory
	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	tcb := globals.BuscarHiloEnPCB(req.Pid, req.Tid)
	if tcb == nil {
		http.Error(w, "No se encontró el hilo", http.StatusNotFound)
		return
	}
	globals.BloquearHilo(tcb)

	//mutex!!
	response, err := request.SolicitarDumpMemory(req.Pid, req.Tid)

	if err != nil || response.StatusCode != http.StatusOK {
		http.Error(w, "Error al solicitar el dump de memoria", http.StatusInternalServerError)
		return
	}

	globals.DesbloquearHilo(tcb)
}

func HandleIO(w http.ResponseWriter, r *http.Request) {
	var io request.RequestIO

	err := commons.DecodificarJSON(r.Body, &io)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Bloquear Hilo
	time.Sleep(time.Duration(io.Tiempo) * time.Second)

	// Desbloquear Hilo y mandarlo a la cola de Ready devuelta
}
