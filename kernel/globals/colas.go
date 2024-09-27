package globals

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"sync"
)

func AgregarProcesoACola(pcb *commons.PCB, cola []*commons.PCB, mutex *sync.Mutex) {
	mutex.Lock()
	cola = append(cola, pcb)
	mutex.Unlock()
}

func SacarProcesoDeCola(pid int, cola *[]*commons.PCB, mutex *sync.Mutex) {
	mutex.Lock()
	for i, pcb := range *cola {
		if pcb.Pid == pid {
			*cola = append((*cola)[:i], (*cola)[i+1:]...)
			mutex.Unlock()
			return
		}
	}
	mutex.Unlock()
}

func AgregarHiloACola(tcb *commons.TCB, cola *[]*commons.TCB, mutex *sync.Mutex) {
	mutex.Lock()
	*cola = append(*cola, tcb)
	mutex.Unlock()
}

func SacarHiloDeCola(tid int, cola *[]*commons.TCB, mutex *sync.Mutex) {
	mutex.Lock()
	for i, tcb := range *cola {
		if tcb.Tid == tid {
			*cola = append((*cola)[:i], (*cola)[i+1:]...)
			mutex.Unlock()
			return
		}
	}
	mutex.Unlock()
}

func BuscarPCBEnColas(pid int) *commons.PCB {
	if pcb := Estructura.procesos[pid]; pcb != nil {
		return pcb
	}

	return nil
}
