package processes

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/threads"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/request"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
	"net/http"
	"strconv"
)

func ProcesoInicial(argumentos []string) {
	pseudocodigo := argumentos[1]
	tamanio, _ := strconv.Atoi(argumentos[2])
	prioridadHiloMain := 0

	CrearProceso(pseudocodigo, tamanio, prioridadHiloMain)
}

func CrearProceso(pseudocodigo string, tamanioMemoria int, prioridad int) int {
	pcb := CrearPCB(pseudocodigo, tamanioMemoria, prioridad)

	log.Printf("## (%d:0) Se crea el proceso - Estado: NEW", pcb.Pid)

	if len(globals.Estructura.ColaNew) == 0 {
		log.Println("Cola NEW está vacía, solicitando memoria.")

		// Solicitar espacio en memoria
		respuestaMemoria, err := SolicitarProcesoMemoria(pcb.Pid, pseudocodigo, tamanioMemoria)

		if err != nil {
			log.Println("Error al solicitar espacio en memoria.")
		}
		// Si la memoria aceptó el proceso, crearlo y pasarlo a READY
		if respuestaMemoria.StatusCode == http.StatusOK {

			threads.CrearHilo(pcb.Pid, prioridad, pseudocodigo)

			log.Printf("Hilo main movido a READY")
		} else {
			if respuestaMemoria.StatusCode == http.StatusConflict {
				log.Println("Memoria no tiene suficiente espacio. Proceso en espera.")
			}
			queues.AgregarProcesoACola(pcb, globals.Estructura.ColaNew)
		}
	} else {
		log.Println("Cola NEW no está vacía, proceso se encola en NEW.")

		queues.AgregarProcesoACola(pcb, globals.Estructura.ColaNew)
	}

	return http.StatusOK
}

func CrearPCB(pseudocodigo string, tamanio int, prioridad int) *commons.PCB {
	pcb := commons.PCB{
		Pid:               globals.Estructura.ContadorPid,
		Estado:            "NEW",
		Tid:               []commons.TCB{},
		ContadorHilos:     0,
		Tamanio:           tamanio,
		PseudoCodigoHilo0: pseudocodigo,
		PrioridadTID0:     prioridad,
		Mutex:             []commons.Mutex{},
	}

	globals.Estructura.ContadorPid++

	globals.Estructura.Procesos[pcb.Pid] = &pcb

	return &pcb
}

func FinalizarProceso(pid int) {
	req := request.RequestFinalizarProceso{
		Pid: pid,
	}

	solicitudCodificada, err := commons.CodificarJSON(req)

	if err != nil {
		log.Println("Error al codificar la solicitud de finalización de proceso")
		return
	}

	response := cliente.Post(globals.KConfig.IpMemory, globals.KConfig.PortMemory, "finalizar_proceso", solicitudCodificada)

	if response.StatusCode == http.StatusOK {
		pcb := queues.BuscarPCBEnColas(pid)

		defer delete(globals.Estructura.Procesos, pid)

		for _, tcb := range pcb.Tid {
			threads.FinalizarHilo(pid, tcb.Tid)
		}

		pcb.Estado = "EXIT"

		log.Printf("## Finaliza el proceso %d", pid)

		if len(globals.Estructura.ColaNew) != 0 {
			procesoNuevo := globals.Estructura.ColaNew[0]
			queues.SacarProcesoDeCola(procesoNuevo.Pid, &globals.Estructura.ColaNew)
			threads.CrearHilo(procesoNuevo.Pid, procesoNuevo.PrioridadTID0, procesoNuevo.PseudoCodigoHilo0)
			queues.AgregarHiloACola(threads.BuscarHiloEnPCB(procesoNuevo.Pid, 0), &globals.Estructura.ColaReady)
		}

	} else {
		log.Printf("## Error al finalizar el proceso %d", pid)
	}
}

func SolicitarProcesoMemoria(pid int, pseudocodigo string, tamanio int) (*http.Response, error) {
	req := request.RequestProcessCreate{
		Pid:            pid,
		Pseudocodigo:   pseudocodigo,
		TamanioMemoria: tamanio,
	}

	solicitudCodificada, err := commons.CodificarJSON(req)

	if err != nil {
		return nil, err
	}

	return cliente.Post(globals.KConfig.IpMemory, globals.KConfig.PortMemory, "process", solicitudCodificada), nil
}
