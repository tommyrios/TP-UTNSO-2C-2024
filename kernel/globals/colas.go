package globals

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"sync"
)

type Colas struct {
	Mutex    sync.Mutex
	Procesos []commons.PCB
	Hilos    []commons.TCB
}

var ColaNew = &Colas{
	Procesos: []commons.PCB{},
	Hilos:    []commons.TCB{},
}

var ColaReady = &Colas{
	Procesos: []commons.PCB{},
	Hilos:    []commons.TCB{},
}
var ColaBlocked = &Colas{
	Procesos: []commons.PCB{},
	Hilos:    []commons.TCB{},
}

func AgregarProcesoACola(pcb commons.PCB, cola *Colas) {
	cola.Mutex.Lock()
	cola.Procesos = append(cola.Procesos, pcb)
	cola.Mutex.Unlock()
}

func SacarProcesoDeCola(pid int, cola *Colas) {
	cola.Mutex.Lock()
	for i, pcb := range cola.Procesos {
		if pcb.Pid == pid {
			cola.Procesos = append(cola.Procesos[:i], cola.Procesos[i+1:]...)
			cola.Mutex.Unlock()
			return
		}
	}
	cola.Mutex.Unlock()
}

func AgregarHiloACola(tcb commons.TCB, cola *Colas) {
	cola.Mutex.Lock()
	cola.Hilos = append(cola.Hilos, tcb)
	cola.Mutex.Unlock()
}

func SacarHiloDeCola(tid int, cola *Colas) {
	cola.Mutex.Lock()
	for i, tcb := range cola.Hilos {
		if tcb.Tid == tid {
			cola.Hilos = append(cola.Hilos[:i], cola.Hilos[i+1:]...)
			cola.Mutex.Unlock()
			return
		}
	}
	cola.Mutex.Unlock()
}

func BuscarPCBEnColas(pid int) *commons.PCB {
	colas := []*Colas{ColaNew, ColaReady, ColaBlocked}

	for _, cola := range colas {
		cola.Mutex.Lock()
		for _, pcb := range cola.Procesos {
			if pcb.Pid == pid {
				cola.Mutex.Unlock()
				return &pcb
			}
		}
		cola.Mutex.Unlock()
	}

	return nil
}
