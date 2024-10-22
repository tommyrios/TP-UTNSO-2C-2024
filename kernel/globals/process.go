package globals

import (
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/request"
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

	if len(Estructura.colaNew) == 0 {
		log.Println("Cola NEW está vacía, solicitando memoria.")

		// Solicitar espacio en memoria
		respuestaMemoria, err := request.SolicitarProcesoMemoria(pseudocodigo, tamanioMemoria)

		if err != nil {
			log.Println("Error al solicitar espacio en memoria.")
			return
		}

		// Si la memoria aceptó el proceso, crearlo y pasarlo a READY
		if respuestaMemoria.StatusCode == http.StatusOK {

			CrearHilo(pcb.Pid, prioridad, pseudocodigo) // Crear el hilo TID 0

			log.Println("Proceso creado y movido a READY")

		} else {
			log.Println("Memoria no tiene suficiente espacio. Proceso en espera.")

			AgregarProcesoACola(pcb, Estructura.colaNew)
		}

	} else {
		log.Println("Cola NEW no está vacía, proceso se encola en NEW.")

		AgregarProcesoACola(pcb, Estructura.colaNew)
	}
	log.Printf("## (%d:0) Se crea el proceso - Estado: NEW", pcb.Pid)
}

func CrearPCB(pseudocodigo string, tamanio int, prioridad int) *commons.PCB {
	pcb := commons.PCB{
		Pid:               Estructura.contadorPid,
		Estado:            "NEW",
		Tid:               []commons.TCB{},
		ContadorHilos:     0,
		Tamanio:           tamanio,
		PseudoCodigoHilo0: pseudocodigo,
		PrioridadTID0:     prioridad,
		Mutex:             []commons.Mutex{},
	}

	Estructura.contadorPid++

	return &pcb
}

func FinalizarProceso(pid int) {
	pcb := BuscarPCBEnColas(pid)

	defer delete(Estructura.procesos, pid)

	for _, tcb := range pcb.Tid {
		FinalizarHilo(pid, tcb.Tid)
	}

	pcb.Estado = "EXIT"

	// esperar el ok de memoria. generar aviso de que se finaliza proceso. Liberar en memoria y avisar al planificador x si alguno en new

	log.Printf("## Finaliza el proceso %d", pid)
}
