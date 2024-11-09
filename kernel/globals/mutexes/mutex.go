package mutexes

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/threads"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
)

func CrearMutex(nombre string, pid int) {
	mutex := commons.Mutex{
		Nombre:          nombre,
		Valor:           1,
		HilosBloqueados: make([]*commons.TCB, 0),
	}

	pcb := globals.Estructura.Procesos[pid]

	pcb.Mutex = append(pcb.Mutex, mutex)
}

// MUTEX_LOCK

func BloquearMutex(nombre string, pid int, tid int) {
	tcb := threads.BuscarHiloEnPCB(pid, tid)

	pcb := queues.BuscarPCBEnColas(pid)

	for _, mutex := range pcb.Mutex {
		if mutex.Nombre == nombre {
			if mutex.Valor == 1 {
				mutex.Valor = 0
				tcb.Mutex = mutex
			} else {
				mutex.HilosBloqueados = append(mutex.HilosBloqueados, tcb)
				threads.BloquearHilo(tcb)
			}
		} else {
			log.Printf("No existe el mutex solicitado con el nombre: %s\n", nombre)
		}
	}
}

// MUTEX_UNLOCK

func DesbloquearMutex(nombre string, pid int, tid int) {
	tcb := threads.BuscarHiloEnPCB(pid, tid)

	pcb := queues.BuscarPCBEnColas(pid)

	for _, mutex := range pcb.Mutex {
		if mutex.Nombre == nombre {
			if tcb.Mutex.Nombre == nombre {
				if len(mutex.HilosBloqueados) > 0 {
					tcbADesbloquear := mutex.HilosBloqueados[0]
					mutex.HilosBloqueados = mutex.HilosBloqueados[1:]
					tcbADesbloquear.Mutex = mutex
					threads.DesbloquearHilo(tcbADesbloquear)
				}
				mutex.Valor++
				tcb.Mutex = commons.Mutex{} // Le saco el mutex al hilo que ejecutó la syscall
			} else {
				log.Printf("El hilo %d no tiene asignado el mutex %s\n", tid, nombre)
			}
		} else {
			log.Printf("No existe el mutex solicitado con el nombre: %s\n", nombre)
		}
	}
}
