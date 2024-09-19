package globals

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

type RequestProceso struct {
	Pseudocodigo   string `json:"pseudocodigo"`
	TamanioMemoria int    `json:"tamanio_memoria"`
}

var KConfig *Config
