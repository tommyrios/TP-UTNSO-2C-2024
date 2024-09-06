package servidor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func RecibirMensaje(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	mensaje := queryParams.Get("mensaje")

	respuesta, err := json.Marshal(fmt.Sprintf("Este es el mensaje: %s", mensaje))

	if err != nil {
		http.Error(w, "Error al codificar los datos como JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuesta)

	log.Println(string(respuesta))
}
