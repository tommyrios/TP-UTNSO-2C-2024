package handlers

//debo recibir el json de La creación de los archivos de DUMP de memoria se crearán al recibir desde Memoria
//la petición de creación de un nuevo archivo de DUMP.
//En la petición deberá venir el nombre del archivo, el tamaño y el contenido a grabar en el mismo
import (
	"bytes"
	"encoding/binary"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sisoputnfrba/tp-golang/filesystem/globals"
	"github.com/sisoputnfrba/tp-golang/filesystem/inicializacion"
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
	var nombreArchivo string
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
		w.Write([]byte("Error al realizar el dump de memoria: No hay espacio disponible"))
		log.Printf("Error al realizar el dump de memoria - (PID:TID) - (%d:%d)", archivo.Pid, archivo.Tid)
	}

	//Paso2: Reservar espacio en el FS para el archivo
	bloqueIndice, bloquesReservados, err := ReservarBloques(bloquesNecesarios)
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

	log.Printf("## Bloque asignado: %d - Archivo: %s - Bloques Libres: %d", bloqueIndice, nombreArchivo, cuantosLibresHay())
	log.Printf("## Archivo Creado: %s - Tamaño: %d", nombreArchivo, archivo.Tamanio)

	//Paso4: Escribir el contenido en los bloques reservados
	err = EscribirContenido(bloqueIndice, bloquesReservados, nombreArchivo, archivo.Contenido)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al realizar el dump de memoria: No se pudo escribir el contenido en los bloques reservados"))
		log.Printf("Error al realizar el dump de memoria - (PID:TID) - (%d:%d)", archivo.Pid, archivo.Tid)
	}

	log.Printf("## Fin de solicitud - Archivo: %s", nombreArchivo)

}

func VerificarEspacioDisponible(tamanio int) (bool, int) {

	bloquesNecesarios := (tamanio + globals.FSConfig.BlockSize - 1) / globals.FSConfig.BlockSize
	bloquesLibres := cuantosLibresHay()
	return bloquesLibres >= bloquesNecesarios, bloquesNecesarios
}

func ReservarBloques(tamanio int) (int, []int, error) {
	var bloquesDatos []int

	rutaBit := filepath.Join(globals.FSConfig.MountDir, "bitmap.dat")
	bitmap := LeerBitmap()

	var bloqueIndice int = -1

	for i := 0; i < len(bitmap)*8; i++ {
		byteIndex := i / 8
		bitIndex := i % 8

		if (bitmap[byteIndex] & (1 << bitIndex)) == 0 {
			if bloqueIndice == -1 {
				bloqueIndice = i // Reservar el primer bloque como índice
				bitmap[byteIndex] |= (1 << bitIndex)
			} else {
				bloquesDatos = append(bloquesDatos, i)
				bitmap[byteIndex] |= (1 << bitIndex)
			}
		}

		if len(bloquesDatos) == tamanio {
			break
		}
	}

	err := ActualizarBitmap(rutaBit, bitmap)
	if err != nil {
		return -1, nil, err
	}

	return bloqueIndice, bloquesDatos, nil
}

func EscribirContenido(bloqueIndice int, bloquesReservados []int, nombreArchivo string, contenido string) error {
	ruta := filepath.Join(globals.FSConfig.MountDir, "bloques.dat")
	archivo, err := os.OpenFile(ruta, os.O_WRONLY, 0644)
	if err != nil {
	}

	defer archivo.Close()

	//PUNTERO
	puntero := SliceToBytes(bloquesReservados)
	offsetIndice := bloqueIndice * globals.FSConfig.BlockSize
	log.Printf("## Acceso Bloque - Archivo: %s - Tipo Bloque: INDICE - Bloque File System %d", nombreArchivo, bloqueIndice)
	_, err = archivo.WriteAt(puntero, int64(offsetIndice))
	if err != nil {
		return err
	}

	//BLOQUES
	partes := DividirEnBloques([]byte(contenido), globals.FSConfig.BlockSize)
	for i, bloque := range bloquesReservados {
		if i >= len(partes) {
			break
		}
		offsetBloque := bloque * globals.FSConfig.BlockSize
		log.Printf("## Acceso Bloque - Archivo: %s - Tipo Bloque: DATOS - Bloque File System %d", nombreArchivo, bloque)
		_, err := archivo.WriteAt(partes[i], int64(offsetBloque))
		if err != nil {
			return err
		}

		// Simular retardo
		time.Sleep(100 * time.Millisecond) // Ajusta el tiempo según sea necesario
	}

	return nil
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

func DividirEnBloques(contenido []byte, size int) [][]byte {
	var bloques [][]byte
	for size < len(contenido) {
		bloques = append(bloques, contenido[:size])
		contenido = contenido[size:]
	}
	bloques = append(bloques, contenido) // Agregar el resto
	return bloques
}

func SliceToBytes(ints []int) []byte {
	buf := new(bytes.Buffer)
	for _, n := range ints {
		binary.Write(buf, binary.LittleEndian, int32(n)) // Convierte cada entero a un int32 en formato LittleEndian
	}
	return buf.Bytes()
}
