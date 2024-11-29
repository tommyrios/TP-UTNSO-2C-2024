package handlers

//debo recibir el json de La creación de los archivos de DUMP de memoria se crearán al recibir desde Memoria
//la petición de creación de un nuevo archivo de DUMP.
//En la petición deberá venir el nombre del archivo, el tamaño y el contenido a grabar en el mismo
import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sisoputnfrba/tp-golang/filesystem/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

type RequestDumpMemory struct {
	Pid       int    `json:"pid"`
	Tid       int    `json:"tid"`
	Tamanio   int    `json:"tamanio"`
	Contenido string `json:"contenido"`
}

func CrearArchivoDump(w http.ResponseWriter, r *http.Request) {
	var archivo RequestDumpMemory
	//var puntero uint32
	err := commons.DecodificarJSON(r.Body, &archivo)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	//Paso1: Verificar si hay espacio disponible en el FS

	espacioDisponible, bloquesNecesarios := VerificarEspacioDisponible(archivo.Tamanio)

	if !espacioDisponible {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al realizar el dump de memoria"))
		log.Printf("Error al realizar el dump de memoria - (PID:TID) - (%d:%d)", archivo.Pid, archivo.Tid)
	}

	//Paso2: Reservar espacio en el FS para el archivo
	bloquesReservados, err := ReservarBloques(bloquesNecesarios)

}

func VerificarEspacioDisponible(tamanio int) (bool, int) {

	bloquesNecesarios := (tamanio + globals.FSConfig.BlockSize - 1) / globals.FSConfig.BlockSize
	bloquesLibres := cuantosLibresHay()
	return bloquesLibres >= bloquesNecesarios, bloquesNecesarios
}

func cuantosLibresHay() int {
	bitmap := LeerBitmap()
	cantidadBloquesLibres := 0
	for _, byte := range bitmap {
		for i := 0; i < 8; i++ {
			if byte&(1<<i) == 0 {
				cantidadBloquesLibres++
			}
		}
	}
	return cantidadBloquesLibres
}

func LeerBitmap() []byte {
	ruta := filepath.Join(globals.FSConfig.MountDir, "bitmap.dat")
	archivo, err := os.ReadFile(ruta)
	if err != nil {
		panic(err)
	}
	return archivo
}

func ActualizarBitmap(ruta string, bitmap []byte) error {
	archivo, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer archivo.Close()

	_, err = archivo.Write(bitmap)
	if err != nil {
		return err
	}

	return nil
}

func ReservarBloques(tamanio int) ([]int, error) {
	var bloquesReservados []int

	rutaBit := filepath.Join(globals.FSConfig.MountDir, "bitmap.dat")
	bitmap := LeerBitmap()

	for i := 0; i < len(bitmap)*8; i++ {
		byteIndex := i / 8
		bitIndex := i % 8

		if (bitmap[byteIndex] & (1 << bitIndex)) == 0 {
			bloquesReservados = append(bloquesReservados, i)
			// Marcar el bloque como ocupado
			bitmap[byteIndex] |= (1 << bitIndex)
		}

		if len(bloquesReservados) == tamanio {
			break
		}
	}

	err := ActualizarBitmap(rutaBit, bitmap)
	if err != nil {
		return nil, err
	}

	return bloquesReservados, nil

}
