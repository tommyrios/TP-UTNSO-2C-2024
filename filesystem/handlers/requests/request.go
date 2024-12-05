package requests

type RequestDumpMemory struct {
	Pid       int    `json:"pid"`
	Tid       int    `json:"tid"`
	Tamanio   int    `json:"tamanio"`
	Contenido []byte `json:"contenido"`
}
