package globals

import (
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"sync"
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
	procesos       map[int]*commons.PCB // Mapa de procesos activos
	colaNew        []*commons.PCB       // Cola de hilos nuevo
	colaReady      []*commons.TCB       // Cola de hilos listos para ejecución
	colaBloqueados []*commons.TCB       // Cola de hilos bloqueados
	hiloExecute    *commons.TCB         // Hilo en ejecución

	contadorPid int // PID autoincremental

	mutexContador  sync.Mutex
	mutexProcesos  sync.Mutex
	mutexColaNew   sync.Mutex
	mutexReady     sync.Mutex
	mutexBloqueado sync.Mutex
}

var Estructura = &Kernel{
	procesos:       make(map[int]*commons.PCB),
	colaReady:      []*commons.TCB{},
	colaBloqueados: []*commons.TCB{},
	hiloExecute:    nil,
	contadorPid:    1,
}

var KConfig *Config
