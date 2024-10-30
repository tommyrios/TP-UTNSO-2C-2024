package globals

import "github.com/sisoputnfrba/tp-golang/utils/commons"

type Config struct {
	Port            int    `json:"port"`
	MemorySize      int    `json:"memory_size"`
	InstructionPath string `json:"instruction_path"`
	ResponseDelay   int    `json:"response_delay"`
	IpKernel        string `json:"ip_kernel"`
	PortKernel      int    `json:"port_kernel"`
	IpCpu           string `json:"ip_cpu"`
	PortCpu         int    `json:"port_cpu"`
	IpFileSystem    string `json:"ip_filesystem"`
	PortFileSystem  int    `json:"port_filesystem"`
	Scheme          string `json:"scheme"`
	SearchAlgorithm string `json:"search_algorithm"`
	Partitions      []int  `json:"partitions"`
	LogLevel        string `json:"log_level"`
}

var MConfig *Config

var MemoriaUsuario []byte

type ContextoProceso struct {
	Base   int // Dirección base de la memoria del proceso
	Limite int // Límite (tamaño de memoria asignada)
}

type ContextoHilo commons.Registros //typedef de C!!

type InstruccionesHilo struct {
	Instrucciones []string // Instrucciones leídas del pseudocódigo
}

var MemoriaSistema = MemSistema{
	TablaProcesos: make(map[int]ContextoProceso),
	TablaHilos:    make(map[int]map[int]ContextoHilo),
	Pseudocodigos: make(map[int]map[int]InstruccionesHilo),
}

type MemSistema struct {
	TablaProcesos map[int]ContextoProceso           // Tabla de procesos (PID -> Contexto de proceso)
	TablaHilos    map[int]map[int]ContextoHilo      // Tabla de hilos (PID -> TID -> Contexto de hilo)
	Pseudocodigos map[int]map[int]InstruccionesHilo // Pseudocódigos (PID -> TID -> Código)
}
