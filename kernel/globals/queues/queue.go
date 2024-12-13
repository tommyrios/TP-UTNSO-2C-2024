package queues

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func AgregarProcesoACola(pcb *commons.PCB, cola *[]*commons.PCB) {
	*cola = append(*cola, pcb)
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

func SacarHiloDeCola(tid int, pid int, cola *[]*commons.TCB) {
	for i, tcb := range *cola {
		if tcb.Tid == tid && tcb.Pid == pid {
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

func BuscarTCBenPCB(pid int, tid int) *commons.TCB {
	pcb := BuscarPCBEnColas(pid)

	for _, tcb := range pcb.Tid {
		if tcb.Tid == tid {
			return tcb
		}
	}

	return nil
}

func ConsultaBloqueado(pid int, tid int) bool {
	for _, tcb := range globals.Estructura.ColaBloqueados {
		if tcb.Pid == pid && tcb.Tid == tid {
			return true
		}
	}

	return false
}

func ConsultaExit(pid int, tid int) bool {
	for _, tcb := range globals.Estructura.ColaExit {
		if tcb.Pid == pid && tcb.Tid == tid {
			return true
		}
	}

	return false
}
