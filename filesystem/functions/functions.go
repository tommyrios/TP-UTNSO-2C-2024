package functions

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/filesystem/globals"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"
)

func HayEspacioDisponible(bloquesNecesarios int) bool {
	contador := 0
	for i := 0; i < len(globals.Bitmap); i++ {
		if globals.Bitmap[i] == 0 { // si hay bloqe libre aumenta el contador
			contador++
		}
		if contador == bloquesNecesarios {
			return true
		}
	}
	return false
}

func buscarBloqueLibre() int {
	for i := 0; i < len(globals.Bitmap); i++ {
		if globals.Bitmap[i] == 0 {
			return i
		}
	}
	return -1
}

func cantidadBloquesLibres() int {
	contador := 0
	for i := 0; i < len(globals.Bitmap); i++ {
		if globals.Bitmap[i] == 0 {
			contador++
		}
	}
	return contador
}

func ActualizarBitmap() {
	archivoBitmap, err := os.OpenFile(globals.FSConfig.MountDir+"/bitmap.dat", os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer archivoBitmap.Close()

	_, err = archivoBitmap.WriteAt(globals.Bitmap, 0) // Escribe el bitmap en el archivo bitmap.dat desde la posición 0
	if err != nil {
		log.Fatal(err)
	}
}

func MarcarBloqueOcupado(posicion int) {
	globals.Bitmap[posicion] = 1
	ActualizarBitmap() // para que se actualice el bitmap en el archivo bitmap.dat
}

func ExisteArchivo(path string) bool {
	_, err := os.Stat(path)    // verifica si existe el archivo, si no existe devuelve error
	return !os.IsNotExist(err) // si no existe devuelve true
}

func inicializarBitmap(pathBitmap string, tamanio int) {
	archivo, err := os.OpenFile(pathBitmap, os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer archivo.Close()

	// Inicializa el archivo con ceros
	globals.Bitmap = make([]byte, tamanio)
	_, err = archivo.Write(globals.Bitmap) // Escribe el bitmap en el archivo bitmap.dat
	if err != nil {
		log.Fatal(err)
	}
}

func CargarBitmap(path string, tamanio int) {
	archivo, err := os.Open(path) // Abre el archivo bitmap.dat
	if err != nil {
		log.Fatal(err)
	}

	defer archivo.Close() // al finalizar Cargar_bitmap() cierra el archivo

	globals.Bitmap = make([]byte, tamanio)
	_, err = archivo.Read(globals.Bitmap)
	if err != nil {
		log.Fatal(err)
	}
}

func VerContenidoBitmapHexa() {
	path := globals.FSConfig.MountDir + "/bitmap.dat"
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// Convertir el contenido a formato hexadecimal y separarlo en pares de dígitos
	hexContent := hex.EncodeToString(file)
	var formattedHex strings.Builder
	for i := 0; i < len(hexContent); i += 2 {
		formattedHex.WriteString(hexContent[i:i+2] + " ")
	}

	// Imprimir el contenido en formato hexadecimal
	fmt.Printf("Contenido de bitmap.dat:\n%s\n", strings.TrimSpace(formattedHex.String()))

}

func escribirEnBloqueIndice(path *os.File, posicionBloqueIndice int, bloquesAsignados []int, nombreArchivo string) {
	// [0] [1,2,3]
	offset := int64(posicionBloqueIndice * globals.FSConfig.BlockSize)

	_, errOffset := path.Seek(offset, 0)

	if errOffset != nil {
		slog.Error("Error al buscar el offset", "Error", errOffset)
		return
	}

	bytesAEscribir := make([]byte, 4*len(bloquesAsignados))

	for i, bloque := range bloquesAsignados {
		// empieza a escribir desde el indice i*4 es decir si i = 0 -> empieza a escribir desde la posicion bytesAEscribir[0] hasta bytesAEscribir[3] = 4 y asi sucesivamente
		binary.LittleEndian.PutUint32(bytesAEscribir[i*4:], uint32(bloque))

	}
	// al final te queda un array de bytes de 4 bytes * el tamaño de los bloquesAsignados
	// por ejemplo se tiene 2 bloques asignados [1,2] => bytesAEscribir = [0,0,0,1,0,0,0,2]
	_, err := path.Write(bytesAEscribir)
	if err != nil {
		slog.Error("Error al escribir en el archivo", "Error", err)
		return
	}

	slog.Info("##", "Acceso Bloque - Archivo:", nombreArchivo, "Tipo Bloque:", "INDICE", "Bloque File System:", posicionBloqueIndice)

	// hay que esperar el tiempo delayBlock en milisegundos ante cada acceso
	time.Sleep(time.Duration(globals.FSConfig.BlockAccessDelay) * time.Millisecond)

}

func escribirEnBloquesDatos(path *os.File, contenido []byte, bloquesAsignados []int, nombreArchivo string) {

	// divide el contenido en subarrays del tamaño  block_size
	//4 bytes
	//[12,23,64,25,84,35]
	//[[12,23,64,25], [84,35]]
	contenidoSubArrays := DividirContenido(contenido)
	i := 0
	// [4,5]
	// [ [0,1,2,3,4,5,6,7] , [8,9,10,11,12,13,14,15] , [16,17,18,19,20,21] ]
	for _, bloque := range bloquesAsignados {

		offset := int64(bloque * globals.FSConfig.BlockSize)
		_, err := path.Seek(offset, 0)
		if err != nil {
			slog.Error("Error al buscar el offset", "Error", err)
			return
		}

		_, err = path.Write(contenidoSubArrays[i])
		if err != nil {
			slog.Error("Error al escribir los DATOS en el archivo", "Error", err)
			return
		}
		i++
		slog.Info("##", "Acceso Bloque - Archivo:", nombreArchivo, "Tipo Bloque:", "DATO", "Bloque File System:", bloque)
		time.Sleep(time.Duration(globals.FSConfig.BlockAccessDelay) * time.Millisecond)
	}

}
func DividirContenido(contenido []byte) [][]byte {
	/*
		tamañoBloque = 8
		i = 0
		contenido = [0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21] bytes
		append(arrayDeSubArrays, contenido[0:8]) // devuelve desde la posicion 0 hasta la 7 ya que el 8 no lo incluye
		arrayDeSubArrays = [[0,1,2,3,4,5,6,7]]
		i = 8
		append(arrayDeSubArrays, contenido[8:16])
		arrayDeSubArrays = [[0,1,2,3,4,5,6,7], [8,9,10,11,12,13,14,15]]
		i = 16
		16+8 < 21? NO -> else -> rrayDeSubArrays = append(arrayDeSubArrays, contenido[i:])
		append(arrayDeSubArrays, contenido[16:] // toma a partir de la posicion 16 hasta el final
		arrayDeSubArrays = [[0,1,2,3,4,5,6,7], [8,9,10,11,12,13,14,15], [16,17,18,19,20,21]]
		arrayDeSubArrays tamaño = [ [8] [8] [6] ]
	*/
	arrayDeSubArrays := make([][]byte, 0)
	for i := 0; i < len(contenido); i += globals.FSConfig.BlockSize {
		if i+globals.FSConfig.BlockSize < len(contenido) {
			arrayDeSubArrays = append(arrayDeSubArrays, contenido[i:i+globals.FSConfig.BlockSize])
		} else {
			arrayDeSubArrays = append(arrayDeSubArrays, contenido[i:])
		}
	}
	return arrayDeSubArrays
}
