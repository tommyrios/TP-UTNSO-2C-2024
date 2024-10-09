package globals

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

var planificador Planificador

func SeleccionarPlanificador(tipo string) {
	switch tipo {
	case "FIFO":
		planificador = &PlanificadorFIFO{}
	case "Prioridades":
		planificador = &PlanificadorPrioridades{}
	case "ColasMultinivel":
		planificador = &PlanificadorColasMultinivel{
			quantum: KConfig.Quantum,
		}
	}
}

type PlanificadorFIFO struct{}

func (p *PlanificadorFIFO) Planificar() *commons.TCB {
	if len(Estructura.colaReady) == 0 {
		return nil
	}
	hilo := Estructura.colaReady[0]
	Estructura.colaReady = Estructura.colaReady[1:]
	return hilo
}

func (p *PlanificadorFIFO) AgregarHilo(tcb *commons.TCB) {
	Estructura.colaReady = append(Estructura.colaReady, tcb)
}

func (p *PlanificadorFIFO) DesalojarHilo(tcb *commons.TCB) {
	// FIFO does not require preemption
}

type PlanificadorPrioridades struct{}

func (p *PlanificadorPrioridades) Planificar() *commons.TCB {
	if len(Estructura.colaReady) == 0 {
		return nil
	}
	indice := 0
	for i, tcb := range Estructura.colaReady {
		if tcb.Prioridad < Estructura.colaReady[indice].Prioridad {
			indice = i
		}
	}
	hilo := Estructura.colaReady[indice]
	Estructura.colaReady = append(Estructura.colaReady[:indice], Estructura.colaReady[indice+1:]...)
	return hilo
}

func (p *PlanificadorPrioridades) AgregarHilo(tcb *commons.TCB) {
	Estructura.colaReady = append(Estructura.colaReady, tcb)
}

func (p *PlanificadorPrioridades) DesalojarHilo(tcb *commons.TCB) {
	p.AgregarHilo(tcb)
}

/*type PlanificadorColasMultinivel struct {
	quantum int
}

func (p *PlanificadorColasMultinivel) Planificar() *commons.TCB {
	for _, cola := range Estructura.colasMultinivel {
		if len(cola) > 0 {
			hilo := cola[0]
			cola = cola[1:]
			return hilo
		}
	}
	return nil
}

func (p *PlanificadorColasMultinivel) AgregarHilo(tcb *commons.TCB) {
	Estructura.colasMultinivel[tcb.Prioridad] = append(Estructura.colasMultinivel[tcb.Prioridad], tcb)
}

func (p *PlanificadorColasMultinivel) DesalojarHilo(tcb *commons.TCB) {
	p.AgregarHilo(tcb)
}
*/

type PlanificadorColasMultinivel struct {
	quantum int
}

func (p *PlanificadorColasMultinivel) Planificar() *commons.TCB {
	if len(Estructura.colaReady) == 0 {
		return nil
	}
	hilo := Estructura.colaReady[0]
	Estructura.colaReady = Estructura.colaReady[1:]
	return hilo
}

func (p *PlanificadorColasMultinivel) AgregarHilo(tcb *commons.TCB) {
	inserted := false
	for i, existingTCB := range Estructura.colaReady {
		if tcb.Prioridad < existingTCB.Prioridad {
			Estructura.colaReady = append(Estructura.colaReady[:i], append([]*commons.TCB{tcb}, Estructura.colaReady[i:]...)...)
			inserted = true
			break
		}
	}
	if !inserted {
		Estructura.colaReady = append(Estructura.colaReady, tcb)
	}
}

func (p *PlanificadorColasMultinivel) DesalojarHilo(tcb *commons.TCB) {
	p.AgregarHilo(tcb)
}
