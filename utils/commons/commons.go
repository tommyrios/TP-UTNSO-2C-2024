package commons

import (
	"encoding/json"
	"io"
	"log"
)

type PCB struct {
	Pid               int     `json:"pid"`
	Tid               []*TCB  `json:"tid"`
	Mutex             []Mutex `json:"mutex"`
	ContadorHilos     int     `json:"contador_hilos"`
	Estado            string  `json:"estado"`
	Tamanio           int     `json:"tamanio"`
	PseudoCodigoHilo0 string  `json:"pseudocodigo_hilo_0"`
	PrioridadTID0     int     `json:"prioridadtid0"`
	ProgramCounter    int     `json:"program_counter"`
}

type TCB struct {
	Pid             int       `json:"pid"`
	Tid             int       `json:"tid"`
	Estado          string    `json:"estado"`
	Prioridad       int       `json:"prioridad"`
	Pseudocodigo    string    `json:"pseudocodigo"`
	Mutex           Mutex     `json:"mutex"`
	Registros       Registros `json:"registros"`
	ProgramCounter  int       `json:"program_counter"`
	TcbADesbloquear []*TCB    `json:"tcb_en_espera"`
}

type Mutex struct {
	Nombre          string `json:"nombre"`
	Valor           int    `json:"valor"`
	HilosBloqueados []*TCB `json:"hilos_bloqueados"`
}

type Registros struct {
	PC uint32 `json:"pc"`
	AX uint32 `json:"ax"`
	BX uint32 `json:"bx"`
	CX uint32 `json:"cx"`
	DX uint32 `json:"dx"`
	EX uint32 `json:"ex"`
	FX uint32 `json:"fx"`
	GX uint32 `json:"gx"`
	HX uint32 `json:"hx"`
}

func DecodificarJSON(r io.Reader, requestStruct interface{}) error {
	err := json.NewDecoder(r).Decode(requestStruct)
	if err != nil {
		log.Printf("Error al decodificar JSON: %s\n", err.Error())
	}
	return err
}

func CodificarJSON(responseStruct interface{}) ([]byte, error) {
	requestCodificada, err := json.Marshal(responseStruct)
	if err != nil {
		log.Printf("Error al codificar JSON: %s\n", err.Error())
	}
	return requestCodificada, err
}
