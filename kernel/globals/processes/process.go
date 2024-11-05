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

	log.Println(pseudocodigo, tamanio, prioridadHiloMain)

}

func CrearProceso(pseudocodigo string, tamanioMemoria int, prioridad int) {
	pcb := CrearPCB(pseudocodigo, tamanioMemoria, prioridad)

	if len(globals.Estructura.ColaNew) == 0 {
		log.Println("Cola NEW está vacía, solicitando memoria.")

		// Solicitar espacio en memoria
		respuestaMemoria, err := SolicitarProcesoMemoria(pcb.Pid, pseudocodigo, tamanioMemoria)

		if err != nil {
			log.Println("Error al solicitar espacio en memoria.")
			return
		}

		// Si la memoria aceptó el proceso, crearlo y pasarlo a READY
		if respuestaMemoria.StatusCode == http.StatusOK {

			threads.CrearHilo(pcb.Pid, prioridad, pseudocodigo) // Crear el hilo TID 0

			log.Println("Proceso creado y movido a READY")

		} else {
			if respuestaMemoria.StatusCode == http.StatusConflict {
				log.Println("Memoria no tiene suficiente espacio. Proceso en espera.")
			} else {
				aceptarCompactación()
			}
			queues.AgregarProcesoACola(pcb, globals.Estructura.ColaNew)
		}

	} else {
		log.Println("Cola NEW no está vacía, proceso se encola en NEW.")

		queues.AgregarProcesoACola(pcb, globals.Estructura.ColaNew)
	}
	log.Printf("## (%d:0) Se crea el proceso - Estado: NEW", pcb.Pid)
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

	return &pcb
}

func FinalizarProceso(pid int) {
	pcb := queues.BuscarPCBEnColas(pid)

	defer delete(globals.Estructura.Procesos, pid)

	for _, tcb := range pcb.Tid {
		threads.FinalizarHilo(pid, tcb.Tid)
	}

	pcb.Estado = "EXIT"

	// esperar el ok de memoria. generar aviso de que se finaliza proceso. Liberar en memoria y avisar al planificador x si alguno en new

	log.Printf("## Finaliza el proceso %d", pid)
}

func SolicitarProcesoMemoria(pid int, pseudocodigo string, tamanio int) (*http.Response, error) {
	request := request.RequestProcessCreate{
		Pid:            pid,
		Pseudocodigo:   pseudocodigo,
		TamanioMemoria: tamanio,
	}

	solicitudCodificada, err := commons.CodificarJSON(request)

	if err != nil {
		return nil, err
	}

	return cliente.Post(globals.KConfig.IpMemory, globals.KConfig.PortMemory, "process", solicitudCodificada), nil
}
