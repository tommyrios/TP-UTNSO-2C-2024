package globals

import (
	"errors"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

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

type ContextoProceso struct {
	Base   int // Dirección base de la memoria del proceso
	Limite int // Límite (tamaño de memoria asignada)
}

type ContextoHilo commons.Registros //typedef de C!!

type InstruccionesHilo struct {
	Instrucciones []string // Instrucciones leídas del pseudocódigo
}

type MemSistema struct {
	TablaProcesos map[int]ContextoProceso           // Tabla de procesos (PID -> Contexto de proceso)
	TablaHilos    map[int]map[int]ContextoHilo      // Tabla de hilos (PID -> TID -> Contexto de hilo)
	Pseudocodigos map[int]map[int]InstruccionesHilo // Pseudocódigos (PID -> TID -> Código)
}

var MemoriaSistema = MemSistema{
	TablaProcesos: make(map[int]ContextoProceso),
	TablaHilos:    make(map[int]map[int]ContextoHilo),
	Pseudocodigos: make(map[int]map[int]InstruccionesHilo),
}

type MemUsuario struct {
	datos       []byte
	particiones []Particion
}

var MemoriaUsuario = MemUsuario{
	datos:       make([]byte, MConfig.MemorySize),
	particiones: []Particion{},
}

type Particion struct {
	base   int
	limite int
	//agregar ocupado o no
}

func inicializarMemoria() {
	if MConfig.Scheme == "FIJAS" {
		base := 0
		for _, tamaño := range MConfig.Partitions {
			if base+tamaño > MConfig.MemorySize {
				errors.New("error: Particiones fijas exceden el tamaño total de memoria")
			}

			// Crear una nueva partición y añadirla a la lista
			nuevaParticion := Particion{
				base:   base,
				limite: tamaño,
			}
			MemoriaUsuario.particiones = append(MemoriaUsuario.particiones, nuevaParticion)
			base += tamaño
		}

		if base != MConfig.MemorySize {
			fmt.Println("Advertencia: No se ha utilizado la memoria completa en particiones fijas.")
		}
	}

	fmt.Println("Memoria inicializada con éxito.")
}
