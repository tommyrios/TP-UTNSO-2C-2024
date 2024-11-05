package schedulers

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
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
		switch globals.KConfig.SchedulerAlgorithm {
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
		if len(globals.Estructura.ColaReady) == 0 {
			continue
		}

		select {
		case <-globals.CpuLibre:
			globals.Estructura.HiloExecute = globals.Estructura.ColaReady[0]
			globals.Estructura.ColaReady = globals.Estructura.ColaReady[1:]

			go func() {
				time.Sleep(time.Duration(globals.KConfig.Quantum)) // Preguntar cómo sabríamos si terminó el proceso

				globals.CpuLibre <- true
			}()
		}
	}
}

func ManejarColaReadyPriority() {
	if len(globals.Estructura.ColaReady) == 0 {
		return
	}

	// Ordenar la cola de ready por prioridad
	sort.SliceStable(globals.Estructura.ColaReady, func(i, j int) bool {
		return globals.Estructura.ColaReady[i].Prioridad < globals.Estructura.ColaReady[j].Prioridad
	})

	globals.Estructura.HiloExecute = globals.Estructura.ColaReady[0]

	globals.Estructura.ColaReady = globals.Estructura.ColaReady[1:]

}

func ManejarColaReadyCMN() {
	for {
		if len(globals.Estructura.ColaReady) == 0 {
			continue
		}

		// Ordenar por prioridad y mantener orden FIFO para la misma prioridad
		sort.SliceStable(globals.Estructura.ColaReady, func(i, j int) bool {
			return globals.Estructura.ColaReady[i].Prioridad < globals.Estructura.ColaReady[j].Prioridad
		})

		// Crear un mapa para simular las colas por niveles de prioridad
		priorityMap := make(map[int][]*commons.TCB)
		for _, tcb := range globals.Estructura.ColaReady {
			priorityMap[tcb.Prioridad] = append(priorityMap[tcb.Prioridad], tcb)
		}

		// Iterar por las prioridades, de mayor a menor
		for priority := range priorityMap {
			queue := priorityMap[priority]

			for len(queue) > 0 {
				select {
				case <-globals.CpuLibre:
					globals.Estructura.HiloExecute = queue[0]
					queue = queue[1:]

					go func() {
						quantumAgotado := false
						for !quantumAgotado {
							time.Sleep(time.Duration(globals.KConfig.Quantum))

							if tieneMasPrioridad() {
								// Si llega un hilo de mayor prioridad, desaloja el hilo actual
								priorityMap[priority] = append(priorityMap[priority], globals.Estructura.HiloExecute)
								// Notificar que la CPU está libre
								globals.CpuLibre <- true
								return // Termina la ejecución para dar paso al hilo de mayor prioridad
							}

							quantumAgotado = checkQuantumAgotado()
						}

						// Si el hilo terminó su quantum, se reubica al final de la cola
						if !tieneMasPrioridad() {
							queue = append(queue, globals.Estructura.HiloExecute)
						}

						// Notificar que la CPU está libre
						globals.CpuLibre <- true
					}()
				}
			}
		}
	}
}

// Helper para verificar si hay un hilo de mayor prioridad
func tieneMasPrioridad() bool {
	for _, tcb := range globals.Estructura.ColaReady {
		if tcb.Prioridad < globals.Estructura.HiloExecute.Prioridad {
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
