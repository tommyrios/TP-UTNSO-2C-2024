package globals

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"net/http"
)

func IniciarProceso(w http.ResponseWriter, r *http.Request) {

	pcb := commons.PCB{
		Pid:           commons.PidCounter,
		Tid:           []int{},
		ContadorHilos: 1,
		Estado:        "NEW",
	}

	commons.MutexPidCounter.Lock()
	commons.PidCounter++
	commons.MutexPidCounter.Unlock()

	tcbMain := IniciarHilo(pcb.Pid, 0, 0)

	pcb.Tid = append(pcb.Tid, tcbMain.Tid)

	commons.ColaNew.Mutex.Lock()
	AgregarProceso(pcb, commons.ColaNew)
	commons.ColaNew.Mutex.Unlock()

	// TODO: Verificar que hay espacio en memoria para poner el proceso en READY
}

func IniciarHilo(pid int, prioridad int, tid int) commons.TCB {
	return commons.TCB{
		Pid:       pid,
		Tid:       tid,
		Prioridad: prioridad,
	}
}

func AgregarProceso(pcb commons.PCB, cola *commons.Colas) {
	cola.Procesos = append(cola.Procesos, pcb)
}

func SacarProceso(pcb commons.PCB, cola *commons.Colas) {
	for i, p := range cola.Procesos {
		if p.Pid == pcb.Pid {
			cola.Procesos = append(cola.Procesos[:i], cola.Procesos[i+1:]...)
			break
		}
	}
}
