package handlers

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals/mutexes"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/processes"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/threads"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/request"
	"github.com/sisoputnfrba/tp-golang/memoria/globals/schemes"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"log"
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

	w.WriteHeader(processes.CrearProceso(req.Pseudocodigo, req.TamanioMemoria, req.Prioridad))
	w.Write([]byte("Proceso creado"))
}

func HandleThreadCreate(w http.ResponseWriter, r *http.Request) {
	var req request.RequestThreadCreate
	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	threads.CrearHilo(req.Pid, req.Prioridad, req.Pseudocodigo)
}

func HandleProcessExit(w http.ResponseWriter, r *http.Request) {
	var req request.RequestProcessExit
	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}
	if req.Tid == 0 {
		processes.FinalizarProceso(req.Pid)
	} else {
		http.Error(w, "La finalizacion de un proceso solo puede ser solicitada por el TID 0", http.StatusBadRequest)
	}
}

// THREAD_EXIT Finaliza el hilo que la invoca (el tid que se manda es del propio hilo)

func HandleThreadExit(w http.ResponseWriter, r *http.Request) {
	var req request.RequestThreadExit
	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	threads.FinalizarHilo(req.Pid, req.Tid)
}

//THREAD_JOIN, esta syscall recibe como parámetro un TID, mueve el hilo que la invocó al estado BLOCK hasta que el TID pasado por parámetro finalice.
//En caso de que el TID pasado por parámetro no exista o ya haya finalizado, esta syscall no hace nada y el hilo que la invocó continuará su ejecución.

func HandleThreadJoin(w http.ResponseWriter, r *http.Request) {
	var req request.RequestThreadJoin
	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	threads.BloquearHilo(globals.Estructura.HiloExecute)
	// REPLANIFICAR !!
	tcbParametro := threads.BuscarHiloEnPCB(req.PidParametro, req.TidParametro)

	tcbParametro.TcbADesbloquear = append(tcbParametro.TcbADesbloquear, globals.Estructura.HiloExecute)
}

// THREAD_CANCEL Finaliza el hilo cuyo tid se pasa por parámetro (desde otro hilo)

func HandleThreadCancel(w http.ResponseWriter, r *http.Request) {
	var req request.RequestThreadCancel
	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	if threads.BuscarHiloEnPCB(req.Pid, req.TidAEliminar) != nil {
		threads.FinalizarHilo(req.Pid, req.TidAEliminar)
	}
}

func HandleMutexCreate(w http.ResponseWriter, r *http.Request) {
	var req request.RequestMutex

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	mutexes.CrearMutex(req.Nombre, req.Pid)
}

func HandleMutexLock(w http.ResponseWriter, r *http.Request) {
	var req request.RequestMutex

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	mutexes.BloquearMutex(req.Nombre, req.Pid, req.Tid)
}

func HandleMutexUnlock(w http.ResponseWriter, r *http.Request) {
	var req request.RequestMutex

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	mutexes.DesbloquearMutex(req.Nombre, req.Pid, req.Tid)
}

func HandleDumpMemory(w http.ResponseWriter, r *http.Request) {
	var req request.RequestDumpMemory
	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	tcb := threads.BuscarHiloEnPCB(req.Pid, req.Tid)
	if tcb == nil {
		http.Error(w, "No se encontró el hilo", http.StatusNotFound)
		return
	}
	threads.BloquearHilo(tcb)

	//mutex!!
	response, err := SolicitarDumpMemory(req.Pid, req.Tid)

	if err != nil || response.StatusCode != http.StatusOK {
		http.Error(w, "Error al solicitar el dump de memoria", http.StatusInternalServerError)
		processes.FinalizarProceso(req.Pid)
		return
	}

	threads.DesbloquearHilo(tcb)
}

func HandleIO(w http.ResponseWriter, r *http.Request) {
	var req request.RequestIO

	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	tcb := globals.Estructura.HiloExecute
	threads.BloquearHilo(tcb)
	// REPLANIFICAR !!
	time.Sleep(time.Duration(req.Tiempo) * time.Second)

	// Desbloquear Hilo y mandarlo a la cola de Ready devuelta
	threads.DesbloquearHilo(tcb)
	log.Printf("## (%d:%d) finalizó IO y pasa a READY", tcb.Pid, tcb.Tid)
}

func SolicitarDumpMemory(pid int, tid int) (*http.Response, error) {
	request := request.RequestDumpMemory{
		Pid: pid,
		Tid: tid,
	}
	requestCodificado, _ := commons.CodificarJSON(request)
	cliente.Post(globals.KConfig.IpMemory, globals.KConfig.PortMemory, "dump", requestCodificado)
	return nil, nil
}

func Dispatch(pcb *commons.PCB, tid int) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(request.RequestDispatcher{PCB: pcb, Tid: tid})
	if err != nil {
		return nil, err
	}

	return cliente.Post(globals.KConfig.IpCpu, globals.KConfig.PortCpu, "dispatch", requestBody), err
}

func Interrupt(interruption string, pid int) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(request.RequestInterrupcion{Razon: interruption, Pid: pid})
	if err != nil {
		return nil, err
	}

	return cliente.Post(globals.KConfig.IpCpu, globals.KConfig.PortCpu, "interrupt", requestBody), err
}

func HandleCompactacion(w http.ResponseWriter, r *http.Request) {
	//pausarPlanificacion() // Pausar planificación de corto plazo

	// Responder a Memoria para permitir la compactación
	w.WriteHeader(http.StatusOK)

	// Notificar a Memoria que puede proceder
	schemes.CompactacionCond.L.Lock()
	schemes.CompactacionCond.Signal() // Aviso a Memoria que puede comenzar la compactación
	schemes.CompactacionCond.L.Unlock()
}

func HandleCompactacionFinalizada(w http.ResponseWriter, r *http.Request) {
	//reanudarPlanificacion() // Reanudar planificación
}
