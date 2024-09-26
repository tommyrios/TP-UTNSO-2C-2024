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

	if len(ColaNew.Procesos) == 0 {
		log.Println("Cola NEW está vacía, solicitando memoria.")

		// Solicitar espacio en memoria
		respuestaMemoria, err := request.SolicitarProcesoMemoria(pseudocodigo, tamanioMemoria)

		if err != nil {
			log.Println("Error al solicitar espacio en memoria.")
			return
		}

		// Si la memoria aceptó el proceso, crearlo y pasarlo a READY
		if respuestaMemoria.StatusCode == http.StatusOK {
			AgregarProcesoACola(pcb, ColaReady)

			go CrearHilo(pcb.Pid, prioridad, pseudocodigo) // Crear el hilo TID 0

			log.Println("Proceso creado y movido a READY")

		} else {
			log.Println("Memoria no tiene suficiente espacio. Proceso en espera.")

			AgregarProcesoACola(pcb, ColaNew)
		}

	} else {
		log.Println("Cola NEW no está vacía, proceso se encola en NEW.")

		AgregarProcesoACola(pcb, ColaNew)
	}
	log.Printf("## (%d:0) Se crea el proceso - Estado: NEW", pcb.Pid)
}

func CrearPCB(pseudocodigo string, tamanio int, prioridad int) commons.PCB {
	pcb := commons.PCB{
		Pid:           PidCounter,
		Estado:        "NEW",
		Tid:           []commons.TCB{},
		ContadorHilos: 0,
		Tamanio:       tamanio,
		PseudoCodigo:  pseudocodigo,
		PrioridadTID0: prioridad,
		Mutex:         []sync.Mutex{},
	}

	MutexPidCounter.Lock()
	PidCounter++
	MutexPidCounter.Unlock()

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

	AgregarHiloACola(tcb, ColaReady)

	log.Printf("## (%d:%d) Se crea el Hilo - Estado: READY", pcb.Pid, tcb.Tid)

	return tcb
}

func FinalizarProceso(pid int) {
	pcb := BuscarPCBEnColas(pid)

	for _, tcb := range pcb.Tid {
		go FinalizarHilo(pid, tcb.Tid)
	}

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
	//Estrategia de hacer una lista de procesos activos en una estructura kernel VER!!!
	//delete(kernel.processes, tid)
	// Pasarlo de cola

	log.Printf("## (%d:%d) Finaliza el hilo", pid, tid)
}

/*type Kernel struct {
	processes    map[int]*PCB // Mapa de procesos activos
	readyQueue   []*TCB       // Cola de hilos listos para ejecución
	blockedQueue []*TCB       // Cola de hilos bloqueados
	nextPID      int          // PID autoincremental
	nextTID      int          // TID autoincremental
}

k := Kernel{
	processes:    make(map[int]*PCB),
	readyQueue:   []*TCB{},
	blockedQueue: []*TCB{},
	nextPID:      1,
	nextTID:      1,
}



func ExitProcess(pid int, k *Kernel) error {
	pcb, ok := k.processes[pid]
	if !ok {
		return fmt.Errorf("Proceso con PID %d no encontrado", pid)
	}

	delete(k.processes, pid)
	log.Printf("Proceso con PID %d finalizado", pid)
	return nil
}

func (k *Kernel) CreateProcess(priority int) *PCB {
	pid := k.nextPID
	k.nextPID++

	pcb := &PCB{
		PID: pid,
		Priority: priority,
		TIDs: []int{},
		Mutexes: []int{},
	}

	k.processes[pid] = pcb

	// Crear el primer hilo (TID 0) para el proceso
	tid := k.nextTID
	k.nextTID++
	tcb := &TCB{
		TID: tid,
		PID: pid,
		Priority: priority,
		State: "NEW",
	}
	pcb.TIDs = append(pcb.TIDs, tid)

	// Añadir el hilo a la cola de READY
	k.readyQueue = append(k.readyQueue, tcb)

	type PCB struct {
		PID      int       // Identificador del proceso
		TIDs     []int     // Hilos asociados al proceso
		Priority int       // Prioridad del proceso
		Mutexes  []int     // Lista de mutexes
	}
*/
