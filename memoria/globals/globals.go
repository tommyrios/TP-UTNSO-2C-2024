package globals

type Partition struct{}

type Config struct {
	Port            int         `json:"port"`
	MemorySize      int         `json:"memory_size"`
	InstructionPath string      `json:"instruction_path"`
	ResponseDelay   int         `json:"response_delay"`
	IpKernel        string      `json:"ip_kernel"`
	PortKernel      int         `json:"port_kernel"`
	IpCpu           string      `json:"ip_cpu"`
	PortCpu         int         `json:"port_cpu"`
	IpFileSystem    string      `json:"ip_filesystem"`
	PortFileSystem  int         `json:"port_filesystem"`
	Scheme          string      `json:"scheme"`
	SearchAlgorithm string      `json:"search_algorithm"`
	Partitions      []Partition `json:"partitions"`
	LogLevel        string      `json:"log_level"`
}

var MConfig *Config
