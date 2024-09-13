package globals

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"net/http"
)

func IniciarProceso(w http.ResponseWriter, r *http.Request) {
	pcb := commons.PCB{
		Pid:    commons.PidCounter,
		Tid:    []int{},
		Estado: "NEW",
	}

	// TODO: Agregar mutex para incrementar el contador de Pid porque se puede interrumpir
	commons.PidCounter++

	tcbMain := commons.TCB{
		Tid:       0,
		Prioridad: 0,
	}

	pcb.Tid = append(pcb.Tid, tcbMain.Tid)

	AgregarProceso(pcb, commons.ColaNew)
	// TODO: Verificar que hay espacio en memoria para poner el proceso en READY
}

func IniciarHilo(pid int, prioridad int) {

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
