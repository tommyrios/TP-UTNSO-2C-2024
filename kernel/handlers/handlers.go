package handlers

import (
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/request"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func HandleProcessCreate(w http.ResponseWriter, r *http.Request) {
	var req request.RequestProcessCreate
	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	globals.CrearProceso(req.Pseudocodigo, req.TamanioMemoria, req.Prioridad)
}

func HandleThreadCreate(w http.ResponseWriter, r *http.Request) {
	var req request.RequestThreadCreate
	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	globals.CrearHilo(req.Pid, req.Prioridad, req.Pseudocodigo)
}

func HandleProcessExit(w http.ResponseWriter, r *http.Request) {
	var req request.RequestProcessExit
	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}
	if req.Tid == 0 {
		globals.FinalizarProceso(req.Pid)
	} else {
		http.Error(w, "La finalizacion de un proceso solo puede ser solicitada por el TID 0", http.StatusBadRequest)
	}

	//Falta avisar a memoria la finalizacion del proceso

}

// THREAD_EXIT Finaliza el hilo que la invoca (el tid que se manda es del propio hilo)

func HandleThreadExit(w http.ResponseWriter, r *http.Request) {
	var req request.RequestThreadExit
	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	globals.FinalizarHilo(req.Pid, req.Tid)
}

func HandleThreadJoin(w http.ResponseWriter, r *http.Request) {
	var req request.RequestThreadJoin
	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Bloquear hilo (tid: request.tid) para darle lugar a que ejecute el hilo (tid: request.tidParametro) y luego desbloquearlo

	//rutina con if esperando a que finalice el otro y cuando finalice el otro, desbloquear hilo

}

// THREAD_CANCEL Finaliza el hilo cuyo tid se pasa por parámetro (desde otro hilo)

func HandleThreadCancel(w http.ResponseWriter, r *http.Request) {
	var req request.RequestThreadCancel
	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	if globals.BuscarHiloEnPCB(req.Pid, req.TidAEliminar) != nil {
		globals.FinalizarHilo(req.Pid, req.TidAEliminar)
	}
}

func HandleMutexCreate(w http.ResponseWriter, r *http.Request) {
	var req request.RequestMutex

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	globals.CrearMutex(req.Nombre, req.Pid)
}

func HandleMutexLock(w http.ResponseWriter, r *http.Request) {
	var req request.RequestMutex

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	globals.BloquearMutex(req.Nombre, req.Pid, req.Tid)
}

func HandleMutexUnlock(w http.ResponseWriter, r *http.Request) {
	var req request.RequestMutex

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	globals.DesbloquearMutex(req.Nombre, req.Pid, req.Tid)
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
		globals.FinalizarProceso(req.Pid)
		return
	}

	globals.DesbloquearHilo(tcb)
}

func HandleIO(w http.ResponseWriter, r *http.Request) {
	var req request.RequestIO

	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Bloquear Hilo
	time.Sleep(time.Duration(req.Tiempo) * time.Second)

	// Desbloquear Hilo y mandarlo a la cola de Ready devuelta
}
