package handlers

import (
	"github.com/sisoputnfrba/tp-golang/filesystem/functions"
	"github.com/sisoputnfrba/tp-golang/filesystem/handlers/requests"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"net/http"
	"strings"
	"time"
)

func CrearArchivo(w http.ResponseWriter, r *http.Request) {
	var archivo requests.Archivo

	err := commons.DecodificarJSON(r.Body, &archivo)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	timestamp := time.Now().Format("15:04:05.000")
	timestamp = strings.Replace(timestamp, ".", ":", 1)

	//slog.Info("Se recibió el archivo: ", archivo)

	resp := functions.CrearArchivo(archivo.Pid, archivo.Tid, timestamp, archivo.Tamanio, archivo.Contenido)

	if resp == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("No hay más espacio disponible"))
	} else if resp == 1 {
		// Responder con un mensaje de éxito
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Error al crear el archivo", http.StatusInternalServerError)
	}

}
