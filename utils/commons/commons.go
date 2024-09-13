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
	Pid     int    `json:"pid"`
	Tid     []int  `json:"tid"`
	Estado  string `json:"estado"`
	Tamanio int    `json:"tamanio"`
	Cola    *Colas `json:"cola"`
}

type TCB struct {
	Tid       int `json:"tid"`
	Prioridad int `json:"prioridad"`
}

type Colas struct {
	mutex    sync.Mutex
	Procesos []PCB
}

var ColaNew *Colas
var ColaReady *Colas
var ColaBlocked *Colas

var PidCounter int = 1

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

func DecodificarJSON(r io.Reader, requestStruct interface{}) error {
	err := json.NewDecoder(r).Decode(requestStruct)
	if err != nil {
		log.Printf("Error al decodificar JSON: %s\n", err.Error())
	}
	return err
}

func CodificarJSON(w io.Writer, responseStruct interface{}) error {
	err := json.NewEncoder(w).Encode(responseStruct)
	if err != nil {
		log.Printf("Error al codificar JSON: %s\n", err.Error())
	}
	return err
}
