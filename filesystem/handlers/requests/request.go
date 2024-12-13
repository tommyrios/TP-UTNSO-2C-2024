package requests

type Archivo struct {
	Pid       uint32 `json:"pid"`
	Tid       uint32 `json:"tid"`
	Tamanio   int    `json:"tamanio"`
	Contenido []byte `json:"contenido"`
}

type Metadata struct {
	IndexBlock int `json:"index_block"`
	Size       int `json:"size"`
}
