package commons

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
)

type Mensaje struct {
	Mensaje string `json:"mensaje"`
}

type PCB struct {
	Pid           int          `json:"pid"`
	Tid           []int        `json:"tid"`
	Mutex         []sync.Mutex `json:"mutex"`
	ContadorHilos int          `json:"contador_hilos"`
	Estado        string       `json:"estado"`
	Tamanio       int          `json:"tamanio"`
}

type TCB struct {
	Pid       int `json:"pid"`
	Tid       int `json:"tid"`
	Prioridad int `json:"prioridad"`
}

type Colas struct {
	Mutex    sync.Mutex
	Procesos []PCB
	Hilos    []TCB
}

var ColaNew = &Colas{
	Procesos: []PCB{},
	Hilos:    []TCB{},
}
var ColaReady = &Colas{
	Procesos: []PCB{},
	Hilos:    []TCB{},
}
var ColaBlocked = &Colas{
	Procesos: []PCB{},
	Hilos:    []TCB{},
}

var PidCounter int = 1
var MutexPidCounter sync.Mutex

// w es el cuerpo de la respuesta y r es el cuerpo de la solicitud
func RecibirMensaje(w http.ResponseWriter, r *http.Request) {
	log.Printf("Método: %s", r.Method)

	var mensaje Mensaje

	if r.Body == nil {
		http.Error(w, "Cuerpo de solicitud vacío", http.StatusBadRequest)
		return
	}

	err := DecodificarJSON(r.Body, &mensaje)

	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Mensaje recibido %+v\n", mensaje)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Mensaje recibido"))
}

// r es el cuerpo de la solicitud y requestStruct es la estructura a la que se decodificará el JSON
func DecodificarJSON(r io.Reader, requestStruct interface{}) error {
	err := json.NewDecoder(r).Decode(requestStruct)
	if err != nil {
		log.Printf("Error al decodificar JSON: %s\n", err.Error())
	}
	return err
}

// w es el cuerpo de la respuesta y responseStruct es la estructura que se codificará en JSON
func CodificarJSON(responseStruct interface{}) ([]byte, error) {
	requestCodificada, err := json.Marshal(responseStruct)
	if err != nil {
		log.Printf("Error al codificar JSON: %s\n", err.Error())
	}
	return requestCodificada, err
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

type PCB struct {
	Pid            int       `json:"pid"`
	ProgramCounter int       `json:"ProgramCounter"`
	Quantum        int       `json:"quantum"`
	Registros      Registros `json:"registros"`
}

type PedidoInstruccion struct {
	Pid int    `json:"pid"`
	PC  uint32 `json:"pc"`
}
