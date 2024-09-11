package commons

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Mensaje struct {
	Mensaje string `json:"mensaje"`
}

func RecibirMensaje(w http.ResponseWriter, r *http.Request) {
	log.Printf("Método: %s", r.Method)
	if r.Method == "POST" {
		// Leer el cuerpo de la solicitud
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Log para ver el cuerpo sin procesar
		log.Printf("Cuerpo de la solicitud: %s", string(body))

		// Procesar los datos JSON
		var data map[string]string
		if err := json.Unmarshal(body, &data); err != nil {
			log.Println("Error al deserializar JSON:", err)
			http.Error(w, "Datos JSON inválidos", http.StatusBadRequest)
			return
		}

		log.Println("Datos recibidos:", data)
	} else {
		log.Println("Método no permitido:", r.Method)
	}
}
