package threads

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/request"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
	"net/http"
)

func CrearHilo(pid int, prioridad int, pseudocodigo string) {
	pcb := queues.BuscarPCBEnColas(pid)

	tcb := commons.TCB{
		Pid:          pcb.Pid,
		Tid:          pcb.ContadorHilos,
		Prioridad:    prioridad,
		Pseudocodigo: pseudocodigo,
	}

	pcb.ContadorHilos++

	pcb.Tid = append(pcb.Tid, tcb) // Chequear después si hay que agregar un mutex

	queues.AgregarHiloACola(&tcb, &globals.Estructura.ColaReady)

	//Avisar a memoria creacion de hilo!!! no hace falta la respuesta

	//<-globals.HilosReady

	log.Printf("## (%d:%d) Se crea el Hilo - Estado: READY", pcb.Pid, tcb.Tid)
}

func FinalizarHilo(pid int, tid int) {
	req := request.RequestFinalizarHilo{
		Pid: pid,
		Tid: tid,
	}

	solicitudCodificada, err := commons.CodificarJSON(req)

	if err != nil {
		log.Println("Error al codificar la solicitud de finalización de proceso")
		return
	}

	response := cliente.Post(globals.KConfig.IpMemory, globals.KConfig.PortMemory, "finalizar_hilo", solicitudCodificada)

	if response.StatusCode == http.StatusOK {
		pcb := queues.BuscarPCBEnColas(pid)
		tcb := BuscarHiloEnPCB(pid, tid)

		defer queues.SacarHiloDeCola(tid, queues.BuscarColaDeHilo(tcb))

		tcb.Estado = "EXIT"

		for i, thread := range pcb.Tid {
			if thread.Tid == tid {
				pcb.Tid = append(pcb.Tid[:i], pcb.Tid[i+1:]...)
				break
			}
		}

		for _, tcb := range tcb.TcbADesbloquear {
			DesbloquearHilo(tcb)
		}
	}

	log.Printf("## (%d:%d) Finaliza el hilo", pid, tid)

	commons.CpuLibre <- true
}

func BuscarHiloEnPCB(pid int, tid int) *commons.TCB {
	pcb := queues.BuscarPCBEnColas(pid)

	for _, tcb := range pcb.Tid {
		if tcb.Tid == tid {
			return &tcb
		}
	}

	return nil
}

func DesbloquearHilo(tcb *commons.TCB) {
	tcb.Estado = "READY"

	queues.SacarHiloDeCola(tcb.Tid, &globals.Estructura.ColaBloqueados)

	queues.AgregarHiloACola(tcb, &globals.Estructura.ColaReady)
}

func BloquearHilo(tcb *commons.TCB) {
	tcb.Estado = "BLOCKED"
	globals.Estructura.HiloExecute = nil
	queues.AgregarHiloACola(tcb, &globals.Estructura.ColaBloqueados)
	// VER que onda la CPU, cómo le avisamos o qué hace
}
