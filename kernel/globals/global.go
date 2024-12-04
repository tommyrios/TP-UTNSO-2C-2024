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
	ColaIO         []*request.RequestIO // Cola de hilos esperando por IO
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
	ColaIO:         []*request.RequestIO{},
	HiloExecute:    nil,
	ContadorPid:    1,
	MtxReady:       &sync.Mutex{},
}

var KConfig *Config
var MutexPlanificacion sync.Mutex
var Planificar = make(chan bool, 1) // Verificar si no hay que poner
var IO chan int
