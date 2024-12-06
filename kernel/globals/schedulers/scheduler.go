package schedulers

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/threads"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
	"sort"
)

var mu = globals.Estructura.MtxReady

func ManejarColaReady() {
	switch globals.KConfig.SchedulerAlgorithm {
	case "FIFO":
		go ManejarColaReadyFIFO()
	case "PRIORIDADES", "CMN":
		go ManejarColaReadyPriority()
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
					handlers.Interrupt("INTERRUPCION", globals.Estructura.HiloExecute.Pid, globals.Estructura.HiloExecute.Tid)
				}
			} else {
				globals.Estructura.HiloExecute = globals.Estructura.ColaReady[0]
			}
			mu.Unlock()
		}
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
		PrintColaReady()
		mu.Unlock()

		ExecuteThread(hiloAEjecutar.Pid, hiloAEjecutar.Tid)
	}
}

func PrintColaReady() {
	fmt.Println("Estado actual de ColaReady:")
	for _, hilo := range globals.Estructura.ColaReady {
		fmt.Printf("Pid: %d, Tid: %d, Prioridad: %d, Estado: %s\n", hilo.Pid, hilo.Tid, hilo.Prioridad, hilo.Estado)
	}
	fmt.Println("---------------")
	fmt.Println("Estado actual de ColaBlock:")
	for _, hilo := range globals.Estructura.ColaBloqueados {
		fmt.Printf("Pid: %d, Tid: %d, Prioridad: %d, Estado: %s\n", hilo.Pid, hilo.Tid, hilo.Prioridad, hilo.Estado)
	}
	fmt.Println("---------------")
	fmt.Println("Estado actual de ColaNew:")
	for _, proceso := range globals.Estructura.ColaNew {
		fmt.Printf("Pid: %d, Estado: %s\n", proceso.Pid, proceso.Estado)
	}
	fmt.Println("---------------")
}

func ExecuteThread(pid int, tid int) {
	if _, err := handlers.Dispatch(pid, tid); err != nil {
		log.Printf("Error al enviar el PID y TID %d al CPU.", pid)
		threads.FinalizarHilo(pid, tid)
	}
}
