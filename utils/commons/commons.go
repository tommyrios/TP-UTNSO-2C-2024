package commons

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Mensaje struct {
	Mensaje string `json:"mensaje"`
}

type PCB struct {
	Pid           int     `json:"pid"`
	Tid           []TCB   `json:"tid"`
	Mutex         []Mutex `json:"mutex"`
	ContadorHilos int     `json:"contador_hilos"`
	Estado        string  `json:"estado"`
	Tamanio       int     `json:"tamanio"`
	PseudoCodigo  string  `json:"pseudocodigo"`
	PrioridadTID0 int     `json:"prioridadtid0"`
}

type TCB struct {
	Pid           int    `json:"pid"`
	Tid           int    `json:"tid"`
	Estado        string `json:"estado"`
	Prioridad     int    `json:"prioridad"`
	Instrucciones string `json:"instrucciones"`
	Mutex         Mutex  `json:"mutex"`
}

type Mutex struct {
	Nombre          string `json:"nombre"`
	Valor           int    `json:"valor"`
	HilosBloqueados []*TCB `json:"hilos_bloqueados"`
}

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
