package globals

import (
	"log"
	"net/http"
	"sync"

	"github.com/sisoputnfrba/tp-golang/handlers/request"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func CrearProceso(pseudocodigo string, tamanioMemoria int, prioridad int) {
	pcb := CrearPCB(pseudocodigo, tamanioMemoria, prioridad)

	if len(commons.ColaNew.Procesos) == 0 {
		log.Println("Cola NEW está vacía, solicitando memoria.")

		// Solicitar espacio en memoria
		respuestaMemoria, err := request.SolicitarProcesoMemoria(pseudocodigo, tamanioMemoria)

		if err != nil {
			log.Println("Error al solicitar espacio en memoria.")
			return
		}

		// Si la memoria aceptó el proceso, crearlo y pasarlo a READY
		if respuestaMemoria.StatusCode == http.StatusOK {
			AgregarProcesoACola(pcb, commons.ColaReady)

			go CrearHilo(pcb.Pid, prioridad, pseudocodigo) // Crear el hilo TID 0

			log.Println("Proceso creado y movido a READY")

		} else {
			log.Println("Memoria no tiene suficiente espacio. Proceso en espera.")

			AgregarProcesoACola(pcb, commons.ColaNew)
		}

	} else {
		log.Println("Cola NEW no está vacía, proceso se encola en NEW.")

		AgregarProcesoACola(pcb, commons.ColaNew)
	}
	log.Printf("## (%d:0) Se crea el proceso - Estado: NEW", pcb.Pid)
}

func CrearPCB(pseudocodigo string, tamanio int, prioridad int) commons.PCB {
	pcb := commons.PCB{
		Pid:           commons.PidCounter,
		Estado:        "NEW",
		Tid:           []commons.TCB{},
		ContadorHilos: 0,
		Tamanio:       tamanio,
		PseudoCodigo:  pseudocodigo,
		PrioridadTID0: prioridad,
		Mutex:         []sync.Mutex{},
	}

	commons.MutexPidCounter.Lock()
	commons.PidCounter++
	commons.MutexPidCounter.Unlock()

	return pcb
}

func CrearHilo(pid int, prioridad int, instrucciones string) commons.TCB {
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

	pcb.Tid = append(pcb.Tid, tcb) // Chequear después

	AgregarHiloACola(tcb, commons.ColaReady)

	log.Printf("## (%d:%d) Se crea el Hilo - Estado: READY", pcb.Pid, tcb.Tid)

	return tcb
}

func AgregarProcesoACola(pcb commons.PCB, cola *commons.Colas) {
	cola.Mutex.Lock()
	cola.Procesos = append(cola.Procesos, pcb)
	cola.Mutex.Unlock()
}

func AgregarHiloACola(tcb commons.TCB, cola *commons.Colas) {
	cola.Mutex.Lock()
	cola.Hilos = append(cola.Hilos, tcb)
	cola.Mutex.Unlock()
}

func BuscarPCBEnColas(pid int) *commons.PCB {
	colas := []*commons.Colas{commons.ColaNew, commons.ColaReady, commons.ColaBlocked}

	for _, cola := range colas {
		cola.Mutex.Lock()
		for _, pcb := range cola.Procesos {
			if pcb.Pid == pid {
				cola.Mutex.Unlock()
				return &pcb
			}
		}
		cola.Mutex.Unlock()
	}

	return nil
}
