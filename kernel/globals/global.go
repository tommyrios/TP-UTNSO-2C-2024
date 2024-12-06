package globals

import (
	"github.com/sisoputnfrba/tp-golang/kernel/handlers/request"
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
	Procesos       map[int]*commons.PCB // Mapa de Procesos activos
	ColaNew        []*commons.PCB       // Cola de hilos nuevo
	ColaReady      []*commons.TCB       // Cola de hilos listos para ejecución
	ColaBloqueados []*commons.TCB       // Cola de hilos bloqueados
	ColaExit       []*commons.TCB       // Cola de hilos finalizados
	ColaIO         []*commons.TCB       // Cola de hilos esperando por IO
	HiloExecute    *commons.TCB         // Hilo en ejecución
	ContadorPid    int                  // PID autoincremental
	MtxReady       *sync.Mutex
}

var Estructura = &Kernel{
	Procesos:       make(map[int]*commons.PCB),
	ColaNew:        []*commons.PCB{},
	ColaReady:      []*commons.TCB{},
	ColaBloqueados: []*commons.TCB{},
	ColaExit:       []*commons.TCB{},
	ColaIO:         []*commons.TCB{},
	HiloExecute:    nil,
	ContadorPid:    1,
	MtxReady:       &sync.Mutex{},
}

var IOChannel chan request.RequestIO

var KConfig *Config

var CpuLibre = make(chan bool, 1)

var Planificar = make(chan bool)
