package mutexes

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/threads"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log/slog"
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

func BloquearMutex(nombre string, pid int, tid int) {
	tcb := threads.BuscarHiloEnPCB(pid, tid)

	pcb := queues.BuscarPCBEnColas(pid)

	for i := range pcb.Mutex {
		if pcb.Mutex[i].Nombre == nombre {
			if pcb.Mutex[i].Valor == 1 {
				pcb.Mutex[i].Valor--
				tcb.Mutex = pcb.Mutex[i]
				return
			} else {
				pcb.Mutex[i].HilosBloqueados = append(pcb.Mutex[i].HilosBloqueados, tcb)
				slog.Info(fmt.Sprintf("## (<%d>:<%d>) - Bloqueado por: MUTEX", pid, tid))
				threads.BloquearHilo(tcb)
				return
			}
		}
	}

	slog.Debug(fmt.Sprintf("No existe el mutex solicitado con el nombre: %s\n", nombre))
}

func DesbloquearMutex(nombre string, pid int, tid int) {
	tcb := threads.BuscarHiloEnPCB(pid, tid)
	pcb := queues.BuscarPCBEnColas(pid)

	for i := range pcb.Mutex {
		if pcb.Mutex[i].Nombre == nombre {
			if tcb.Mutex.Nombre == nombre {
				if len(pcb.Mutex[i].HilosBloqueados) > 0 {
					tcbADesbloquear := pcb.Mutex[i].HilosBloqueados[0]
					pcb.Mutex[i].HilosBloqueados = pcb.Mutex[i].HilosBloqueados[1:]
					tcbADesbloquear.Mutex = pcb.Mutex[i]
					threads.DesbloquearHilo(tcbADesbloquear)
				} else {
					pcb.Mutex[i].Valor++
				}
				tcb.Mutex = commons.Mutex{}
				return
			} else {
				slog.Debug(fmt.Sprintf("El hilo %d no tiene asignado el mutex %s\n", tid, nombre))
				return
			}
		}
	}

	slog.Debug(fmt.Sprintf("No existe el mutex solicitado con el nombre: %s\n", nombre))
}
