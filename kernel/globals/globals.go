package globals

type Config struct {
	IpMemoria          string `json:"ip_memory"`
	PuertoMemoria      int    `json:"port_memory"`
	IpCpu              string `json:"ip_cpu"`
	PuertoCpu          int    `json:"port_cpu"`
	SchedulerAlgorithm string `json:"scheduler_algorithm"`
	Quantum            int    `json:"quantum"`
	LogLevel           string `json:"log_level"`
}

var KConfig *Config
