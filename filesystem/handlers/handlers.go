package handlers

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/filesystem/functions"
	"github.com/sisoputnfrba/tp-golang/filesystem/handlers/requests"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

/*func CrearArchivoDump(w http.ResponseWriter, r *http.Request) {
	// Leer el cuerpo de la solicitud
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Aquí agregas el log para inspeccionar el JSON recibido
	slog.Info("Cuerpo recibido", "body", string(body))

	var archivo requests.Archivo

	err = json.Unmarshal(body, &archivo)
	if err != nil {
		http.Error(w, "Error al deserializar el JSON", http.StatusBadRequest)
		slog.Error("Error al deserializar el JSON", "Error", err)
		return
	}

	timestamp := time.Now().Format("15:04:05.000")
	timestamp = strings.Replace(timestamp, ".", ":", 1)

	resp := functions.CrearArchivo(archivo.Pid, archivo.Tid, timestamp, archivo.Tamanio, archivo.Contenido)

	if resp == 0 {
		http.Error(w, "Error al crear el archivo", http.StatusInternalServerError)
	} else if resp == 1 {
		// Responder con un mensaje de éxito
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Error al crear el archivo", http.StatusInternalServerError)
	}

}*/

func CrearArchivo(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Leer el cuerpo de la solicitud
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusInternalServerError)
			return
		}
	}(r.Body)

	var archivo requests.Archivo

	// Deserealizo el JSON a la estructura
	err = json.Unmarshal(body, &archivo)
	if err != nil {
		http.Error(w, "Error al deserializar el JSON", http.StatusBadRequest)
		slog.Error("Error al deserializar el JSON", "Error", err)
		return
	}

	timestamp := time.Now().Format("15:04:05.000")
	timestamp = strings.Replace(timestamp, ".", ":", 1)

	slog.Info("Se recibió el archivo ", "archivo", archivo)
	resp := functions.CrearArchivo(archivo.Pid, archivo.Tid, timestamp, archivo.Tamanio, archivo.Contenido)

	if resp == 0 {
		http.Error(w, "Error al crear el archivo", http.StatusInternalServerError)
	} else if resp == 1 {
		// Responder con un mensaje de éxito
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Error al crear el archivo", http.StatusInternalServerError)
	}

}
