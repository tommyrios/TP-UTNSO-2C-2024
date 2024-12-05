package functions

import (
	"bytes"
	"encoding/binary"
	"github.com/sisoputnfrba/tp-golang/filesystem/globals"
	"log"
	"os"
	"path/filepath"
	"time"
)

func VerificarEspacioDisponible(tamanio int) (bool, int) {

	bloquesNecesarios := (tamanio + globals.FSConfig.BlockSize - 1) / globals.FSConfig.BlockSize
	bloquesLibres := CuantosLibresHay()
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
				bitmap[byteIndex] |= 1 << bitIndex
			} else {
				bloquesDatos = append(bloquesDatos, i)
				bitmap[byteIndex] |= 1 << bitIndex
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

func EscribirContenido(bloqueIndice int, bloquesReservados []int, nombreArchivo string, contenido []byte) error {
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
	partes := DividirEnBloques(contenido, globals.FSConfig.BlockSize)
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

		time.Sleep(time.Duration(globals.FSConfig.BlockAccessDelay) * time.Millisecond)

	}

	return nil
}

func CuantosLibresHay() int {
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
