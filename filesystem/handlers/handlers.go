package handlers

//debo recibir el json de La creación de los archivos de DUMP de memoria se crearán al recibir desde Memoria
//la petición de creación de un nuevo archivo de DUMP.
//En la petición deberá venir el nombre del archivo, el tamaño y el contenido a grabar en el mismo
import (
	"github.com/sisoputnfrba/tp-golang/filesystem/globals"
	"github.com/sisoputnfrba/tp-golang/filesystem/globals/functions"
	"github.com/sisoputnfrba/tp-golang/filesystem/handlers/requests"
	"github.com/sisoputnfrba/tp-golang/filesystem/inicializacion"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
	"net/http"
	"path/filepath"
)

func CrearArchivoDump(w http.ResponseWriter, r *http.Request) {
	var archivo requests.RequestDumpMemory
	var nombreArchivo string
	//var puntero uint32
	err := commons.DecodificarJSON(r.Body, &archivo)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	//Paso1: Verificar si hay espacio disponible en el FS

	espacioDisponible, bloquesNecesarios := functions.VerificarEspacioDisponible(archivo.Tamanio)

	if !espacioDisponible {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al realizar el dump de memoria: No hay espacio disponible"))
		log.Printf("Error al realizar el dump de memoria - (PID:TID) - (%d:%d)", archivo.Pid, archivo.Tid)
	}

	//Paso2: Reservar espacio en el FS para el archivo
	bloqueIndice, bloquesReservados, err := functions.ReservarBloques(bloquesNecesarios)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al realizar el dump de memoria: No se pudieron reservar los bloques necesarios"))
		log.Printf("Error al realizar el dump de memoria - (PID:TID) - (%d:%d)", archivo.Pid, archivo.Tid)
	}

	//Paso3: Crear el archivo de metadata
	ruta := filepath.Join(globals.FSConfig.MountDir, "/files")
	err, nombreArchivo = inicializacion.CrearMetadata(ruta, archivo.Pid, archivo.Tid, archivo.Tamanio, bloqueIndice)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al realizar el dump de memoria: No se pudo crear el archivo de metadata"))
		log.Printf("Error al realizar el dump de memoria - (PID:TID) - (%d:%d)", archivo.Pid, archivo.Tid)
	}

	log.Printf("## Bloque asignado: %d - Archivo: %s - Bloques Libres: %d", bloqueIndice, nombreArchivo, functions.CuantosLibresHay())
	log.Printf("## Archivo Creado: %s - Tamaño: %d", nombreArchivo, archivo.Tamanio)

	//Paso4: Escribir el contenido en los bloques reservados
	err = functions.EscribirContenido(bloqueIndice, bloquesReservados, nombreArchivo, archivo.Contenido)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al realizar el dump de memoria: No se pudo escribir el contenido en los bloques reservados"))
		log.Printf("Error al realizar el dump de memoria - (PID:TID) - (%d:%d)", archivo.Pid, archivo.Tid)
	}

	log.Printf("## Fin de solicitud - Archivo: %s", nombreArchivo)

}
