package globals

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
)

func CrearHilo(pid int, prioridad int, pseudocodigo string) {
	pcb := BuscarPCBEnColas(pid)

	tcb := commons.TCB{
		Pid:          pcb.Pid,
		Tid:          pcb.ContadorHilos,
		Prioridad:    prioridad,
		Pseudocodigo: pseudocodigo,
	}

	pcb.ContadorHilos++

	pcb.Tid = append(pcb.Tid, tcb) // Chequear después si hay que agregar un mutex

	AgregarHiloACola(&tcb, &Estructura.colaReady)

	//Avisar a memoria creacion de hilo!!! no hace falta la respuesta

	log.Printf("## (%d:%d) Se crea el Hilo - Estado: READY", pcb.Pid, tcb.Tid)
}

func FinalizarHilo(pid int, tid int) {
	pcb := BuscarPCBEnColas(pid)
	tcb := BuscarHiloEnPCB(pid, tid)

	defer SacarHiloDeCola(tid, BuscarColaDeHilo(tcb))

	tcb.Estado = "EXIT"

	for i, thread := range pcb.Tid {
		if thread.Tid == tid {
			pcb.Tid = append(pcb.Tid[:i], pcb.Tid[i+1:]...)
			break
		}
	}

	// Y avisarle a memoria que no está más este hilo!!!!
	// Liberar todos los hilos bloqueados por culpa del hilo a finalizar

	log.Printf("## (%d:%d) Finaliza el hilo", pid, tid)
}

func BuscarHiloEnPCB(pid int, tid int) *commons.TCB {
	pcb := BuscarPCBEnColas(pid)

	for _, tcb := range pcb.Tid {
		if tcb.Tid == tid {
			return &tcb
		}
	}

	return nil
}

func DesbloquearHilo(tcb *commons.TCB) {
	tcb.Estado = "READY"

	SacarHiloDeCola(tcb.Tid, &Estructura.colaBloqueados)

	AgregarHiloACola(tcb, &Estructura.colaReady)
}

func BloquearHilo(tcb *commons.TCB) {
	tcb.Estado = "BLOCKED"
	Estructura.hiloExecute = nil
	AgregarHiloACola(tcb, &Estructura.colaBloqueados)
	// VER que onda la CPU, cómo le avisamos o qué hace
}
