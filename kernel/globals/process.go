package globals

import (
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/request"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
	"net/http"
	"sync"
)

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

			go CrearHilo(pcb.Pid, prioridad, pseudocodigo) // Crear el hilo TID 0

			log.Println("Proceso creado y movido a READY")

		} else {
			log.Println("Memoria no tiene suficiente espacio. Proceso en espera.")

			AgregarProcesoACola(pcb, Estructura.colaNew, &Estructura.mutexColaNew)
		}

	} else {
		log.Println("Cola NEW no está vacía, proceso se encola en NEW.")

		AgregarProcesoACola(pcb, Estructura.colaNew, &Estructura.mutexColaNew)
	}
	log.Printf("## (%d:0) Se crea el proceso - Estado: NEW", pcb.Pid)
}

func CrearPCB(pseudocodigo string, tamanio int, prioridad int) *commons.PCB {
	pcb := commons.PCB{
		Pid:           Estructura.contadorPid,
		Estado:        "NEW",
		Tid:           []commons.TCB{},
		ContadorHilos: 0,
		Tamanio:       tamanio,
		PseudoCodigo:  pseudocodigo,
		PrioridadTID0: prioridad,
		Mutex:         []sync.Mutex{},
	}

	Estructura.mutexContador.Lock()
	Estructura.contadorPid++
	Estructura.mutexContador.Unlock()

	return &pcb
}

func CrearHilo(pid int, prioridad int, instrucciones string) {
	pcb := BuscarPCBEnColas(pid)

	tcb := commons.TCB{
		Pid:           pcb.Pid,
		Tid:           pcb.ContadorHilos,
		Prioridad:     prioridad,
		Instrucciones: instrucciones,
	}

	pcb.Mutex[0].Lock()
	pcb.ContadorHilos++
	pcb.Mutex[0].Unlock()

	pcb.Tid = append(pcb.Tid, tcb) // Chequear después si hay que agregar un mutex

	AgregarHiloACola(&tcb, &Estructura.colaReady, &Estructura.mutexReady)

	log.Printf("## (%d:%d) Se crea el Hilo - Estado: READY", pcb.Pid, tcb.Tid)
}

func FinalizarProceso(pid int) {
	pcb := BuscarPCBEnColas(pid)

	defer delete(Estructura.procesos, pid)

	for _, tcb := range pcb.Tid {
		go FinalizarHilo(pid, tcb.Tid)
	}

	pcb.Estado = "EXIT"

	// Liberar memoria

	log.Printf("## Finaliza el proceso %d", pid)
}

func FinalizarHilo(pid int, tid int) {
	pcb := BuscarPCBEnColas(pid)

	for i, tcb := range pcb.Tid {
		if tcb.Tid == tid {
			pcb.Tid = append(pcb.Tid[:i], pcb.Tid[i+1:]...)
			break
		}
	}

	// Pasarlo de cola

	log.Printf("## (%d:%d) Finaliza el hilo", pid, tid)
}
