package globals

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type Config struct {
	Port               int    `json:"port"`
	IpMemory           string `json:"ip_memory"`
	PortMemory         int    `json:"port_memory"`
	IpCpu              string `json:"ip_cpu"`
	PortCpu            int    `json:"port_cpu"`
	SchedulerAlgorithm string `json:"scheduler_algorithm"`
	Quantum            int    `json:"quantum"`
	LogLevel           string `json:"log_level"`
}

type Kernel struct {
	Procesos       map[int]*commons.PCB // Mapa de Procesos activos
	ColaNew        []*commons.PCB       // Cola de hilos nuevo
	ColaReady      []*commons.TCB       // Cola de hilos listos para ejecución
	ColaBloqueados []*commons.TCB       // Cola de hilos bloqueados
	ColaExit       []*commons.TCB       // Cola de hilos finalizados
	HiloExecute    *commons.TCB         // Hilo en ejecución

	ContadorPid int // PID autoincremental
}

var Estructura = &Kernel{
	Procesos:       make(map[int]*commons.PCB),
	ColaNew:        []*commons.PCB{},
	ColaReady:      []*commons.TCB{},
	ColaBloqueados: []*commons.TCB{},
	ColaExit:       []*commons.TCB{},
	HiloExecute:    nil,
	ContadorPid:    1,
}

var KConfig *Config
var HilosReady = make(chan int)
var CpuLibre = make(chan bool, 1)
