package handlers

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/mutexes"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/processes"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/threads"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/request"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log/slog"
	"net/http"
	"time"
)

func HandleProcessCreate(w http.ResponseWriter, r *http.Request) {
	var req request.RequestProcessCreate

	err := commons.DecodificarJSON(r.Body, &req)

	slog.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: PROCESS_CREATE", req.Pid, req.Tid))

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	statusCode := processes.CrearProceso(req.Pseudocodigo, req.TamanioMemoria, req.Prioridad)

	//log.Println("Respuesta de la creación de proceso: ", statusCode)

	w.WriteHeader(statusCode)
}

func HandleProcessExit(w http.ResponseWriter, r *http.Request) {
	var req request.RequestProcessExit
	err := commons.DecodificarJSON(r.Body, &req)

	slog.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: PROCESS_EXIT", req.Pid, req.Tid))

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

	slog.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: THREAD_CREATE", req.Pid, req.Tid))

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

	slog.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: THREAD_JOIN", req.Pid, req.Tid))

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

	slog.Info(fmt.Sprintf("## (%d:%d) - Bloqueado por: THREAD_JOIN", tcbExecute.Pid, tcbExecute.Tid))

	w.WriteHeader(http.StatusOK)
}

func HandleThreadCancel(w http.ResponseWriter, r *http.Request) {
	var req request.RequestThreadCancel
	err := commons.DecodificarJSON(r.Body, &req)

	slog.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: THREAD_CANCEL", req.Pid, req.Tid))

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

	slog.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: THREAD_EXIT", req.Pid, req.Tid))

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	threads.FinalizarHilo(req.Pid, req.Tid)

	w.WriteHeader(http.StatusOK)
}

func HandleMutexCreate(w http.ResponseWriter, r *http.Request) {
	var req request.RequestMutex

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	slog.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: MUTEX_CREATE", req.Pid, req.Tid))

	mutexes.CrearMutex(req.Nombre, req.Pid)

	w.WriteHeader(http.StatusOK)
}

func HandleMutexLock(w http.ResponseWriter, r *http.Request) {
	var req request.RequestMutex

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	slog.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: MUTEX_LOCK", req.Pid, req.Tid))

	mutexes.BloquearMutex(req.Nombre, req.Pid, req.Tid)
}

func HandleMutexUnlock(w http.ResponseWriter, r *http.Request) {
	var req request.RequestMutex

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	slog.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: MUTEX_UNLOCK", req.Pid, req.Tid))

	mutexes.DesbloquearMutex(req.Nombre, req.Pid, req.Tid)
}

func HandleDumpMemory(w http.ResponseWriter, r *http.Request) {
	var req request.RequestDumpMemory

	slog.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: DUMP_MEMORY", req.Pid, req.Tid))

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
	statusCode, _, mensaje := SolicitarDumpMemory(req.Pid, req.Tid)

	slog.Debug(fmt.Sprintf("Respuesta al solicitar el dump de memoria: %s", mensaje))

	if statusCode != http.StatusOK {
		processes.FinalizarProceso(req.Pid)
		return
	}

	threads.DesbloquearHilo(tcb)
}

func HandleCompactacion(w http.ResponseWriter, r *http.Request) {
	PausarPlanificacion()

	w.WriteHeader(http.StatusOK)

	slog.Debug(fmt.Sprintf("Compactación aceptada"))
}

func HandleCompactacionFinalizada(w http.ResponseWriter, r *http.Request) {
	ReanudarPlanificacion()

	w.WriteHeader(http.StatusOK)
}

func HandleDesalojoCpu(w http.ResponseWriter, r *http.Request) {
	var req request.RequestDevolucionPCB
	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	globals.Estructura.HiloExecute = nil

	if req.Razon == "SEGMENTATION FAULT" {
		processes.FinalizarProceso(req.Pid)

		slog.Debug(fmt.Sprintf("## (PID:TID) - (%d:%d) - Hilo recibido de CPU - Razon: %s", req.Pid, req.Tid, req.Razon))
	} else {
		if req.Razon == "SYSCALL" || req.Razon == "INTERRUPCION" || (req.Razon == "MEMORY_DUMP" && !queues.ConsultaExit(req.Pid, req.Tid) || req.Razon == "END_OF_QUANTUM") {
			globals.Estructura.MtxReady.Lock()
			if !queues.ConsultaBloqueado(req.Pid, req.Tid) {
				tcb := threads.BuscarHiloEnPCB(req.Pid, req.Tid)

				tcb.Estado = "READY"

				queues.AgregarHiloACola(threads.BuscarHiloEnPCB(req.Pid, req.Tid), &globals.Estructura.ColaReady)
				if req.Razon != "END_OF_QUANTUM" {
					slog.Debug(fmt.Sprintf("## (PID:TID) - (%d:%d) - Hilo recibido de CPU - Razon: %s", req.Pid, req.Tid, req.Razon))
				} else {
					slog.Info(fmt.Sprintf("## (%d:%d) - Desalojado por fin de Quantum", req.Pid, req.Tid))
				}
			}
			globals.Estructura.MtxReady.Unlock()
		}
	}

	w.WriteHeader(http.StatusOK)

	if len(globals.Estructura.ColaReady) != 0 {
		globals.Planificar <- true
		globals.CpuLibre <- true
	}

}

func HandleIO(w http.ResponseWriter, r *http.Request) {
	var req request.RequestIO

	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	slog.Info(fmt.Sprintf("## (%d:%d) - Solicitó syscall: IO por %d ms", req.Pid, req.Tid, req.Tiempo))

	tcb := threads.BuscarHiloEnPCB(req.Pid, req.Tid)

	IO := globals.IO{Tcb: tcb, Tiempo: req.Tiempo}

	if len(globals.Estructura.ColaIO) == 0 {
		threads.BloquearHilo(tcb)
		slog.Info(fmt.Sprintf("## (%d:%d) - Bloqueado por: IO", req.Pid, req.Tid))
		globals.Estructura.ColaIO = append(globals.Estructura.ColaIO, &IO)
		go ejecutarIO()
	} else {
		threads.BloquearHilo(tcb)
		slog.Info(fmt.Sprintf("## (%d:%d) - Bloqueado por: IO", req.Pid, req.Tid))
		globals.Estructura.ColaIO = append(globals.Estructura.ColaIO, &IO)
	}
}

func ejecutarIO() {
	for len(globals.Estructura.ColaIO) != 0 {
		IO := globals.Estructura.ColaIO[0]

		time.Sleep(time.Duration(IO.Tiempo) * time.Millisecond)

		if len(globals.Estructura.ColaReady) == 0 {
			threads.DesbloquearHilo(IO.Tcb)

			globals.Planificar <- true
			globals.CpuLibre <- true
		} else {
			threads.DesbloquearHilo(IO.Tcb)
		}

		slog.Info(fmt.Sprintf("## (%d:%d) - Finalizó IO y pasa a READY", IO.Tcb.Pid, IO.Tcb.Tid))

		globals.Estructura.ColaIO = globals.Estructura.ColaIO[1:]
	}
}

func SolicitarDumpMemory(pid int, tid int) (int, error, string) {
	req := request.RequestDumpMemory{
		Pid: pid,
		Tid: tid,
	}

	requestCodificado, _ := commons.CodificarJSON(req)

	response, mensaje := cliente.Post2(globals.KConfig.IpMemory, globals.KConfig.PortMemory, "memory_dump", requestCodificado)

	defer response.Body.Close()

	return response.StatusCode, nil, string(mensaje)
}

func Dispatch(pid int, tid int) (*http.Response, error) {
	requestBody, err := commons.CodificarJSON(request.RequestDispatcher{Pid: pid, Tid: tid, Quantum: globals.KConfig.Quantum, Scheduler: globals.KConfig.SchedulerAlgorithm})
	if err != nil {
		return nil, err
	}

	return cliente.Post(globals.KConfig.IpCpu, globals.KConfig.PortCpu, "dispatch", requestBody), err
}

func PausarPlanificacion() {
	globals.Estructura.MtxReady.Lock()
}

func ReanudarPlanificacion() {
	globals.Estructura.MtxReady.Unlock()
}
