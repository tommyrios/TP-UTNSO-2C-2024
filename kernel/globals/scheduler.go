package globals

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"sort"
	"time"
)

/* func PrepareProcess(pcb queues.PCB) {
	if pcb.Quantum > 0 && pcb.Quantum < globals.Config.Quantum && globals.IsVirtualRoundRobin() {
		ChangeState(&pcb, queues.PrioritizedReadyProcesses, "READY")

		log.Printf("Cola Ready+: [%s]",
			logs.IntArrayToString(queues.PrioritizedReadyProcesses.GetPids(), ", "))
	} else {
		pcb.Quantum = globals.Config.Quantum
		ChangeState(&pcb, queues.ReadyProcesses, "READY")

		log.Printf("Cola Ready: [%s]",
			logs.IntArrayToString(queues.ReadyProcesses.GetPids(), ", "))
	}

	<-globals.Ready
}

func SetProcessToReady() {
	for {
		globals.New <- 0
		globals.Multiprogramming <- 0

		globals.Plan()

		PrepareProcess(queues.NewProcesses.PopProcess())
	}
}
*/

func ManejarColaReady() {
	for {
		switch KConfig.SchedulerAlgorithm {
		case "FIFO":
			ManejarColaReadyFIFO()
		case "CMN":
			ManejarColaReadyCMN()
		case "PRIORITY":
			ManejarColaReadyPriority()
		}
	}
}

func ManejarColaReadyFIFO() {
	for {
		if len(Estructura.colaReady) == 0 {
			continue
		}

		select {
		case <-CpuLibre:
			Estructura.hiloExecute = Estructura.colaReady[0]
			Estructura.colaReady = Estructura.colaReady[1:]

			go func() {
				time.Sleep(time.Duration(KConfig.Quantum)) // Preguntar cómo sabríamos si terminó el proceso

				CpuLibre <- true
			}()
		}
	}
}

func ManejarColaReadyPriority() {
	if len(Estructura.colaReady) == 0 {
		return
	}

	// Ordenar la cola de ready por prioridad
	sort.SliceStable(Estructura.colaReady, func(i, j int) bool {
		return Estructura.colaReady[i].Prioridad < Estructura.colaReady[j].Prioridad
	})

	Estructura.hiloExecute = Estructura.colaReady[0]

	Estructura.colaReady = Estructura.colaReady[1:]

}

func ManejarColaReadyCMN() {
	for {
		if len(Estructura.colaReady) == 0 {
			continue
		}

		// Ordenar por prioridad y mantener orden FIFO para la misma prioridad
		sort.SliceStable(Estructura.colaReady, func(i, j int) bool {
			return Estructura.colaReady[i].Prioridad < Estructura.colaReady[j].Prioridad
		})

		// Crear un mapa para simular las colas por niveles de prioridad
		priorityMap := make(map[int][]*commons.TCB)
		for _, tcb := range Estructura.colaReady {
			priorityMap[tcb.Prioridad] = append(priorityMap[tcb.Prioridad], tcb)
		}

		// Iterar por las prioridades, de mayor a menor
		for priority := range priorityMap {
			queue := priorityMap[priority]

			for len(queue) > 0 {
				select {
				case <-CpuLibre:
					Estructura.hiloExecute = queue[0]
					queue = queue[1:]

					go func() {
						quantumAgotado := false
						for !quantumAgotado {
							time.Sleep(time.Duration(KConfig.Quantum))

							if tieneMasPrioridad() {
								// Si llega un hilo de mayor prioridad, desaloja el hilo actual
								priorityMap[priority] = append(priorityMap[priority], Estructura.hiloExecute)
								// Notificar que la CPU está libre
								CpuLibre <- true
								return // Termina la ejecución para dar paso al hilo de mayor prioridad
							}

							quantumAgotado = checkQuantumAgotado()
						}

						// Si el hilo terminó su quantum, se reubica al final de la cola
						if !tieneMasPrioridad() {
							queue = append(queue, Estructura.hiloExecute)
						}

						// Notificar que la CPU está libre
						CpuLibre <- true
					}()
				}
			}
		}
	}
}

// Helper para verificar si hay un hilo de mayor prioridad
func tieneMasPrioridad() bool {
	for _, tcb := range Estructura.colaReady {
		if tcb.Prioridad < Estructura.hiloExecute.Prioridad {
			return true
		}
	}
	return false
}

// Helper para verificar si el quantum ha expirado
func checkQuantumAgotado() bool {
	// Lógica para determinar si el quantum ha expirado
	return true // Modifica esta parte para que funcione según el sistema
}
