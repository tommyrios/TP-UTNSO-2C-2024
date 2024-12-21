package functions

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/filesystem/globals"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func CrearArchivo(pid uint32, tid uint32, timestamp string, tamanio int, contenido []byte) int {

	nombreArchivo := strconv.Itoa(int(pid)) + "_" + strconv.Itoa(int(tid)) + "_" + timestamp
	slog.Debug(fmt.Sprintf("nombreArchivo: %s", nombreArchivo))
	err := CrearDirectorio(globals.FSConfig.MountDir + "/files")
	if err != nil {
		return -1
	}

	bloquesNecesarios := (tamanio + globals.FSConfig.BlockSize - 1) / globals.FSConfig.BlockSize
	fmt.Println("Bloques Necesarios:", bloquesNecesarios)

	if HayEspacioDisponible(bloquesNecesarios + 1) {

		directorioArchivo, errBloques := os.OpenFile(globals.FSConfig.MountDir+"/bloques.dat", os.O_RDWR, 0666)
		if errBloques != nil {
			slog.Error("Error al abrir el archivo de bloques", "Error", errBloques)
		}
		defer directorioArchivo.Close()

		bloqueIndice := buscarBloqueLibre()
		if bloqueIndice == -1 {
			return -1
		} else {

			MarcarBloqueOcupado(bloqueIndice)
			slog.Info("##", "Bloque INDICE asignado:", bloqueIndice, "Archivo:", nombreArchivo, "Bloques Libres:", cantidadBloquesLibres())
		}

		bloquesAsignados := make([]int, bloquesNecesarios)

		for i := 0; i < bloquesNecesarios; i++ {

			// busco el primer bloque libre
			posicionBloqueDatos := buscarBloqueLibre()

			MarcarBloqueOcupado(posicionBloqueDatos)

			bloquesAsignados[i] = posicionBloqueDatos

			slog.Info("##", "Bloque DATOS asignado:", posicionBloqueDatos, "Archivo:", nombreArchivo, "Bloques Libres:", cantidadBloquesLibres())

		}

		//Ver contenido del archivo
		VerContenidoBitmapHexa()

		// Creo el archivo de metadatos
		crearMetadata(nombreArchivo, tamanio, bloqueIndice)

		escribirEnBloqueIndice(directorioArchivo, bloqueIndice, bloquesAsignados, nombreArchivo)

		escribirEnBloquesDatos(directorioArchivo, contenido, bloquesAsignados, nombreArchivo)

		slog.Info("##", "Fin de solicitud - Archivo:", nombreArchivo)
		return 1
	} else {
		//ERROR
		slog.Info("NO HAY ESPACIO DISPONIBLE")
		slog.Info("##", "Fin de solicitud - Archivo:", nombreArchivo)
		return 0
	}
}

func CrearDirectorio(pathDirectorio string) error {
	if _, err := os.Stat(pathDirectorio); os.IsNotExist(err) { // Si no existe el directorio lo crea
		err := os.MkdirAll(pathDirectorio, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error al crear el directorio %s: %w", pathDirectorio, err)
		}
	}

	return nil
}

func crearMetadata(nombreArchivo string, pos int, tamanio int) {

	nombreArchivo = strings.ReplaceAll(nombreArchivo, ":", "-") // Cambia ":" por "-"

	pathArchivoMetadata := globals.FSConfig.MountDir + "/files/" + nombreArchivo + ".dmp"

	// Crea el archivo si no existe, si existe lo abre
	archivo, err := os.OpenFile(pathArchivoMetadata, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		slog.Error("Error al crear el archivo:", err)
		return
	}
	defer archivo.Close()

	// escribe en texto el index_block y el tamaño del archivo
	_, err = archivo.WriteString("index_block:" + strconv.Itoa(pos) + "\n" + "size:" + strconv.Itoa(tamanio))

	slog.Info("##", "Archivo Creado:", nombreArchivo, "Tamaño:", tamanio)
}

func CrearArchivoBloques(pathBloques string, tamanio int) {

	archivo, err := os.OpenFile(pathBloques, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		log.Fatal(err)
	}
	defer archivo.Close()

	err = archivo.Truncate(int64(tamanio)) // setea el tamaño del archivo
	if err != nil {
		log.Fatal(err)
	}
}

func CrearArchivoBitmap(pathBitmap string, tamanio int) {

	archivo, err := os.OpenFile(pathBitmap, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		log.Fatal(err)
	}
	defer archivo.Close()

	err = archivo.Truncate(int64(tamanio))
	if err != nil {
		log.Fatal(err)
	}

	inicializarBitmap(pathBitmap, tamanio)

}
