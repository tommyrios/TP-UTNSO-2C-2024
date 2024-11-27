package queues

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func AgregarProcesoACola(pcb *commons.PCB, cola []*commons.PCB) {
	cola = append(cola, pcb)
}

func SacarProcesoDeCola(pid int, cola *[]*commons.PCB) {
	for i, pcb := range *cola {
		if pcb.Pid == pid {
			*cola = append((*cola)[:i], (*cola)[i+1:]...)
			return
		}
	}
}

func AgregarHiloACola(tcb *commons.TCB, cola *[]*commons.TCB) {
	*cola = append(*cola, tcb)
}

func SacarHiloDeCola(tid int, cola *[]*commons.TCB) {
	for i, tcb := range *cola {
		if tcb.Tid == tid {
			*cola = append((*cola)[:i], (*cola)[i+1:]...)
			return
		}
	}
}

func BuscarPCBEnColas(pid int) *commons.PCB {
	if pcb := globals.Estructura.Procesos[pid]; pcb != nil {
		return pcb
	}

	return nil
}

func BuscarColaDeHilo(tcbBuscado *commons.TCB) *[]*commons.TCB {
	switch tcbBuscado.Estado {
	case "READY":
		return &globals.Estructura.ColaReady
	case "BLOCKED":
		return &globals.Estructura.ColaBloqueados
	case "EXIT":
		return &globals.Estructura.ColaExit
	}
	return nil
}
