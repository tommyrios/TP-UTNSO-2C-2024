package globals

type Config struct {
	Port             int    `json:"port"`
	IpMemory         string `json:"ip_memory"`
	PortMemory       int    `json:"port_memory"`
	MountDir         string `json:"mount_dir"`
	BlockSize        int    `json:"block_size"`
	BlockCount       int    `json:"block_count"`
	BlockAccessDelay int    `json:"block_access_delay"`
	LogLevel         string `json:"log_level"`
}

type Metadata struct {
	IndexBlock int `json:"index_block"`
	Size       int `json:"size"`
}

var FSConfig *Config

var FSMetadata *Metadata
