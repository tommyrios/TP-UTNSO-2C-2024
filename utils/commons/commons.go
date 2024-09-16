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
	Pid     int `json:"pid"`
	Pc      int `json:"pc"`
	Quantum int `json:"quantum"`
}
