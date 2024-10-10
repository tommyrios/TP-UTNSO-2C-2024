package globals

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"sort"
	"time"
)

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
	if len(Estructura.colaReady) == 0 {
		return
	}
	for {
		Estructura.hiloExecute = Estructura.colaReady[0]
		Estructura.colaReady = Estructura.colaReady[1:]
	}

}

func ManejarColaReadyPriority() {
	if len(Estructura.colaReady) == 0 {
		return
	}

	// Sort the colaReady by priority, maintaining FIFO for same priority
	sort.SliceStable(Estructura.colaReady, func(i, j int) bool {
		return Estructura.colaReady[i].Prioridad < Estructura.colaReady[j].Prioridad
	})

	for {
		Estructura.hiloExecute = Estructura.colaReady[0]
		Estructura.colaReady = Estructura.colaReady[1:]
	}
}

/*
	func ManejarColaReadyCMN() {
		if len(Estructura.colaReady) == 0 {
			return
		}
		// Sort the colaReady by priority, maintaining FIFO for same priority
		sort.SliceStable(Estructura.colaReady, func(i, j int) bool {
			return Estructura.colaReady[i].Prioridad < Estructura.colaReady[j].Prioridad
		})

		// Create a map to simulate slices within the colaReady queue for each priority level
		priorityMap := make(map[int][]*commons.TCB)
		for _, tcb := range Estructura.colaReady {
			priorityMap[tcb.Prioridad] = append(priorityMap[tcb.Prioridad], tcb)
		}

		// Iterate over the priority levels and implement Round Robin scheduling
		for priority, queue := range priorityMap {
			for len(queue) > 0 {
				// Get the first thread in the queue
				Estructura.hiloExecute = queue[0]
				queue = queue[1:]

				// Simulate execution with a quantum
				time.Sleep(time.Duration(KConfig.Quantum) * time.Millisecond)

				// Check if the thread should be preempted or moved to the end of its queue
				if shouldPreempt(Estructura.hiloExecute) {
					queue = append(queue, Estructura.hiloExecute)
				}

				// Update the priority map
				priorityMap[priority] = queue
			}
		}
	}

// Helper function to determine if a thread should be preempted

	func shouldPreempt(tcb *commons.TCB) bool {
		// Implement the logic to determine if the thread should be preempted
		// For example, based on the remaining execution time or other criteria
		return false
	}
*/
func ManejarColaReadyCMN() {
	if len(Estructura.colaReady) == 0 {
		return
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
			Estructura.hiloExecute = queue[0]
			queue = queue[1:]

			quantumElapsed := false
			for !quantumElapsed {
				time.Sleep(time.Duration(KConfig.Quantum/10) * time.Millisecond)

				if hasHigherPriority() {
					// Si llega un hilo de mayor prioridad, desaloja el hilo actual
					priorityMap[priority] = append(priorityMap[priority], Estructura.hiloExecute)
					// Si el hilo actual es desalojado, debería volver a la cola de ready donde estaba
					return // Termina la ejecución para dar paso al hilo de mayor prioridad
				}

				quantumElapsed = checkQuantumElapsed()
			}

			// Si el hilo terminó su quantum, se reubica al final de la cola
			if !hasHigherPriority() {
				queue = append(queue, Estructura.hiloExecute)
			}

			priorityMap[priority] = queue
		}
	}
}

// Helper para verificar si hay un hilo de mayor prioridad
func hasHigherPriority() bool {
	// Lógica para comprobar si hay un hilo de mayor prioridad en colaReady
	for _, tcb := range Estructura.colaReady {
		if tcb.Prioridad < Estructura.hiloExecute.Prioridad {
			return true
		}
	}
	return false
}

// Helper para verificar si el quantum ha expirado
func checkQuantumElapsed() bool {
	// Lógica para determinar si el quantum ha expirado
	return true // Modifica esta parte para que funcione según el sistema
}
