package schedulers

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/threads"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
	"sort"
	"time"
)

var mu = globals.Estructura.MtxReady

func ManejarColaReady() {
	switch globals.KConfig.SchedulerAlgorithm {
	case "FIFO":
		go ManejarColaReadyFIFO()
	case "CMN":
		go ManejarColaReadyCMN()
	case "PRIORIDADES":
		go ManejarColaReadyPriority()
	}
}

func ManejarHiloRunning() {
	for {
		<-globals.CpuLibre

		mu.Lock()
		if len(globals.Estructura.ColaReady) == 0 {
			mu.Unlock()
			globals.CpuLibre <- true
			continue
		}

		hiloAEjecutar := globals.Estructura.ColaReady[0]

		globals.Estructura.HiloExecute = hiloAEjecutar

		hiloAEjecutar.Estado = "EXEC"

		if len(globals.Estructura.ColaReady) > 1 {
			globals.Estructura.ColaReady = globals.Estructura.ColaReady[1:]
		} else {
			globals.Estructura.ColaReady = []*commons.TCB{}
		}

		mu.Unlock()

		executeThread(hiloAEjecutar.Pid, hiloAEjecutar.Tid)
	}
}

func ManejarColaReadyFIFO() {
	for {
		<-globals.Planificar
	}
}

func ManejarColaReadyPriority() {
	for {
		<-globals.Planificar
		if len(globals.Estructura.ColaReady) != 0 {
			mu.Lock()
			sort.SliceStable(globals.Estructura.ColaReady, func(i, j int) bool {
				return globals.Estructura.ColaReady[i].Prioridad < globals.Estructura.ColaReady[j].Prioridad
			})
			if globals.Estructura.HiloExecute != nil {
				if globals.Estructura.ColaReady[0].Prioridad < globals.Estructura.HiloExecute.Prioridad {
					handlers.Interrupt("Desalojo por prioridad", globals.Estructura.HiloExecute.Pid, globals.Estructura.HiloExecute.Tid)
				}
			} else {
				globals.Estructura.HiloExecute = globals.Estructura.ColaReady[0]
			}
			mu.Unlock()
		}
	}
}

func ManejarColaReadyCMN() {
	for {
		select {
		case <-globals.Planificar:

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

func executeThread(pid int, tid int) {
	if _, err := handlers.Dispatch(pid, tid); err != nil {
		log.Printf("Error al enviar el PID y TID %d al CPU.", pid)
		threads.FinalizarHilo(pid, tid)
	}
}
