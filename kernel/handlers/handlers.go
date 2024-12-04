package handlers

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/mutexes"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/processes"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/threads"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/request"
	"github.com/sisoputnfrba/tp-golang/memoria/globals/schemes"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
	"net/http"
	"sync"
	"time"
)

func HandleProcessCreate(w http.ResponseWriter, r *http.Request) {
	var req request.RequestProcessCreate
	err := commons.DecodificarJSON(r.Body, &req)

	log.Printf("## (%d:%d) - Solicitó syscall: PROCESS_CREATE", req.Pid, req.Tid)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	statusCode := processes.CrearProceso(req.Pseudocodigo, req.TamanioMemoria, req.Prioridad)

	w.WriteHeader(statusCode)
}

func HandleProcessExit(w http.ResponseWriter, r *http.Request) {
	var req request.RequestProcessExit
	err := commons.DecodificarJSON(r.Body, &req)

	log.Printf("## (%d:%d) - Solicitó syscall: PROCESS_EXIT", req.Pid, req.Tid)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}
	if req.Tid == 0 {
		processes.FinalizarProceso(req.Pid)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "La finalizacion de un proceso solo puede ser solicitada por el TID 0", http.StatusBadRequest)
	}
}

func HandleThreadCreate(w http.ResponseWriter, r *http.Request) {
	var req request.RequestThreadCreate
	err := commons.DecodificarJSON(r.Body, &req)

	log.Printf("## (%d:%d) - Solicitó syscall: THREAD_CREATE", req.Pid, req.Tid)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	threads.CrearHilo(req.Pid, req.Prioridad, req.Pseudocodigo)
	w.WriteHeader(http.StatusOK)
}

func HandleThreadJoin(w http.ResponseWriter, r *http.Request) {
	var req request.RequestThreadJoin
	err := commons.DecodificarJSON(r.Body, &req)

	log.Printf("## (%d:%d) - Solicitó syscall: THREAD_JOIN", req.Pid, req.Tid)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	tcbParametro := queues.BuscarTCBenPCB(req.Pid, req.TidParametro)
	tcbExecute := queues.BuscarTCBenPCB(req.Pid, req.Tid)

	if tcbParametro == nil || tcbParametro.Estado == "EXIT" {
		return
	}

	threads.BloquearHilo(tcbExecute)

	tcbParametro.TcbADesbloquear = append(tcbParametro.TcbADesbloquear, tcbExecute)

	w.WriteHeader(http.StatusOK)
}

func HandleThreadCancel(w http.ResponseWriter, r *http.Request) {
	var req request.RequestThreadCancel
	err := commons.DecodificarJSON(r.Body, &req)

	log.Printf("## (%d:%d) - Solicitó syscall: THREAD_CANCEL", req.Pid, req.Tid)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	if threads.BuscarHiloEnPCB(req.Pid, req.TidAEliminar) != nil {
		threads.FinalizarHilo(req.Pid, req.TidAEliminar)
		w.WriteHeader(http.StatusOK)
	}
}

func HandleThreadExit(w http.ResponseWriter, r *http.Request) {
	var req request.RequestThreadExit
	err := commons.DecodificarJSON(r.Body, &req)

	log.Printf("## (%d:%d) - Solicitó syscall: THREAD_EXIT", req.Pid, req.Tid)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	threads.FinalizarHilo(req.Pid, req.Tid)

	w.WriteHeader(http.StatusOK)
}

func HandleMutexCreate(w http.ResponseWriter, r *http.Request) {
	var req request.RequestMutex

	log.Printf("## (%d:%d) - Solicitó syscall: MUTEX_CREATE", req.Pid, req.Tid)

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	mutexes.CrearMutex(req.Nombre, req.Pid)

	w.WriteHeader(http.StatusOK)
}

func HandleMutexLock(w http.ResponseWriter, r *http.Request) {
	var req request.RequestMutex

	log.Printf("## (%d:%d) - Solicitó syscall: MUTEX_LOCK", req.Pid, req.Tid)

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	mutexes.BloquearMutex(req.Nombre, req.Pid, req.Tid)
}

func HandleMutexUnlock(w http.ResponseWriter, r *http.Request) {
	var req request.RequestMutex

	log.Printf("## (%d:%d) - Solicitó syscall: MUTEX_UNLOCK", req.Pid, req.Tid)

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	mutexes.DesbloquearMutex(req.Nombre, req.Pid, req.Tid)
}

func HandleDumpMemory(w http.ResponseWriter, r *http.Request) {
	var req request.RequestDumpMemory

	log.Printf("## (%d:%d) - Solicitó syscall: DUMP_MEMORY", req.Pid, req.Tid)

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
	statusCode, _ := SolicitarDumpMemory(req.Pid, req.Tid)

	if statusCode != http.StatusOK {
		http.Error(w, "Error al solicitar el dump de memoria", http.StatusInternalServerError)
		processes.FinalizarProceso(req.Pid)
		return
	}

	threads.DesbloquearHilo(tcb)
}

var mtxIO sync.Mutex

func HandleIO(w http.ResponseWriter, r *http.Request) {
	var req request.RequestIO

	log.Printf("## (%d:%d) - Solicitó syscall: IO", req.Pid, req.Tid)

	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	mtxIO.Lock()
	globals.Estructura.ColaIO = append(globals.Estructura.ColaIO, &req)
	mtxIO.Unlock()

	globals.IO <- 1

	log.Printf("## (%d:%d) finalizó IO y pasa a READY", req.Pid, req.Tid)
}

func SolicitarDumpMemory(pid int, tid int) (int, error) {
	request := request.RequestDumpMemory{
		Pid: pid,
		Tid: tid,
	}
	requestCodificado, _ := commons.CodificarJSON(request)
	response := cliente.Post(globals.KConfig.IpMemory, globals.KConfig.PortMemory, "dump", requestCodificado)
	return response.StatusCode, nil
}

func ManejadorIO() {
	for {
		<-globals.IO
		mtxIO.Lock()
		if len(globals.Estructura.ColaIO) > 0 {
			req := globals.Estructura.ColaIO[0]
			globals.Estructura.ColaIO = globals.Estructura.ColaIO[1:]
			mtxIO.Unlock()
			threads.BloquearHilo(threads.BuscarHiloEnPCB(req.Pid, req.Tid))
			time.Sleep(time.Duration(req.Tiempo))
			threads.DesbloquearHilo(threads.BuscarHiloEnPCB(req.Pid, req.Tid))
		} else {
			mtxIO.Unlock()
		}
	}
}

func Dispatch(pid int, tid int) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(request.RequestDispatcher{Pid: pid, Tid: tid})
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
	PausarPlanificacion() // Pausar planificación de corto plazo

	// Responder a Memoria para permitir la compactación
	w.WriteHeader(http.StatusOK)

	// Notificar a Memoria que puede proceder
	schemes.CompactacionCond.L.Lock()
	schemes.CompactacionCond.Signal() // Aviso a Memoria que puede comenzar la compactación
	schemes.CompactacionCond.L.Unlock()

	log.Println("Compactación aceptada")
}

func HandleCompactacionFinalizada(w http.ResponseWriter, r *http.Request) {
	ReanudarPlanificacion() // Reanudar planificación
	w.WriteHeader(http.StatusOK)
}

func HandleDesalojoCpu(w http.ResponseWriter, r *http.Request) {
	var req request.RequestDevolucionPCB
	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	if req.Razon == "SEGMENTATION FAULT" {
		processes.FinalizarProceso(req.Pid)
	} else {
		if req.Razon == "SYSCALL" || req.Razon == "INTERRUPCION" {
			globals.Estructura.HiloExecute = nil
			queues.AgregarHiloACola(threads.BuscarHiloEnPCB(req.Pid, req.Tid), &globals.Estructura.ColaReady)
			<-globals.Planificar
		}
	}

	<-commons.CpuLibre
}

func PausarPlanificacion() {
	globals.MutexPlanificacion.Lock()
}

func ReanudarPlanificacion() {
	globals.MutexPlanificacion.Unlock()
}
