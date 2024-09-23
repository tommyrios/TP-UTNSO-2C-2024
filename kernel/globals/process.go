package globals

import (
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
	"net/http"
	"sync"
)

func CrearProceso(pseudocodigo string, tamanioMemoria int, prioridad int) {
	pcb := CrearPCB(pseudocodigo, tamanioMemoria)

	if len(commons.ColaNew.Procesos) == 0 {
		log.Println("Cola NEW está vacía, solicitando memoria.")

		// Solicitar espacio en memoria
		respuestaMemoria, err := SolicitarProcesoMemoria(pseudocodigo, tamanioMemoria)

		if err != nil {
			log.Println("Error al solicitar espacio en memoria.")
			return
		}

		// Si la memoria aceptó el proceso, crearlo y pasarlo a READY
		if respuestaMemoria.StatusCode == http.StatusOK {
			commons.ColaReady.Mutex.Lock()
			AgregarProcesoACola(pcb, commons.ColaReady)
			commons.ColaReady.Mutex.Unlock()

			go IniciarHilo(pcb, prioridad, pcb.ContadorHilos) // Crear el hilo TID 0

			log.Println("Proceso creado y movido a READY")

		} else {
			log.Println("Memoria no tiene suficiente espacio. Proceso en espera.")
		}
	}
}

func SolicitarProcesoMemoria(pseudocodigo string, tamanio int) (*http.Response, error) {
	request := RequestProceso{
		Pseudocodigo:   pseudocodigo,
		TamanioMemoria: tamanio,
	}

	solicitudCodificada, err := commons.CodificarJSON(request)

	if err != nil {
		return nil, err
	}

	return cliente.Post(KConfig.IpMemory, KConfig.PortMemory, "process", solicitudCodificada), nil
}

func CrearPCB(pseudocodigo string, tamanio int) commons.PCB {
	pcb := commons.PCB{
		Pid:           commons.PidCounter,
		Estado:        "NEW",
		Tid:           []commons.TCB{},
		ContadorHilos: 0,
		Tamanio:       tamanio,
		PseudoCodigo:  pseudocodigo,
		Mutex:         []sync.Mutex{},
	}

	commons.MutexPidCounter.Lock()
	commons.PidCounter++
	commons.MutexPidCounter.Unlock()

	commons.ColaNew.Mutex.Lock()
	AgregarProcesoACola(pcb, commons.ColaNew)
	commons.ColaNew.Mutex.Unlock()

	return pcb
}

func IniciarHilo(pcb commons.PCB, prioridad int, tid int) commons.TCB {
	tcb := commons.TCB{
		Pid:       pcb.Pid,
		Tid:       tid,
		Prioridad: prioridad,
	}

	pcb.Tid = append(pcb.Tid, tcb) // Chequear después

	commons.ColaReady.Mutex.Lock()
	AgregarHiloACola(tcb, commons.ColaReady)
	commons.ColaReady.Mutex.Unlock()

	return tcb
}

func AgregarProcesoACola(pcb commons.PCB, cola *commons.Colas) {
	cola.Procesos = append(cola.Procesos, pcb)
}

func AgregarHiloACola(tcb commons.TCB, cola *commons.Colas) {
	cola.Hilos = append(cola.Hilos, tcb)
}
