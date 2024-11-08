package schedulers

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/queues"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/threads"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
	"sort"
	"time"
)

func ManejarColaReady() {
	for {
		globals.CpuLibre <- true
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

func ManejarHiloRunning() {
	for {
		select {
		case <-globals.CpuLibre:
			hiloAEjecutar := globals.Estructura.ColaReady[0]
			pcbHilo := queues.BuscarPCBEnColas(hiloAEjecutar.Pid)

			// Lo asignamos al hilo en ejecución
			globals.Estructura.HiloExecute = hiloAEjecutar

			// Lo eliminamos de la cola de ready
			globals.Estructura.ColaReady = globals.Estructura.ColaReady[1:]

			executeThread(pcbHilo, hiloAEjecutar.Tid)
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
			pasarHiloAEjecutar()
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

func pasarHiloAEjecutar() {
	// Agarramos el primer hilo de la cola de ready
	hiloAEjecutar := globals.Estructura.ColaReady[0]
	pcbHilo := queues.BuscarPCBEnColas(hiloAEjecutar.Pid)

	// Lo asignamos al hilo en ejecución
	globals.Estructura.HiloExecute = hiloAEjecutar

	// Lo eliminamos de la cola de ready
	globals.Estructura.ColaReady = globals.Estructura.ColaReady[1:]

	executeThread(pcbHilo, hiloAEjecutar.Tid)
}

func executeThread(pcb *commons.PCB, tid int) {
	if _, err := handlers.Dispatch(pcb, tid); err != nil {
		log.Printf("Error al enviar el PCB %d al CPU.", pcb.Pid)
		threads.FinalizarHilo(pcb.Pid, tid)
	}
}
