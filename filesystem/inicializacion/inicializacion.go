package inicializacion

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/sisoputnfrba/tp-golang/filesystem/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func IniciarFileSystem(mountDir string) error {
	err := os.MkdirAll(mountDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	RutaBitmap := filepath.Join(mountDir, "bitmap.dat")
	RutaBloques := filepath.Join(mountDir, "bloques.dat")

	//Crear/Verificar bitmap
	if _, err := os.Stat(RutaBitmap); os.IsNotExist(err) {
		log.Printf("Archivo %s no encontrado. Creando uno nuevo.", RutaBitmap)
		if err := crearBitmap(RutaBitmap); err != nil {
			panic(err)
		}
	} else {
		log.Printf("Archivo encontrado: %s.", RutaBitmap)
	}

	//Crear/Verificar bloques
	if _, err := os.Stat(RutaBloques); os.IsNotExist(err) {
		log.Printf("Archivo %s no encontrado. Creando uno nuevo.", RutaBloques)
		if err := crearBloques(RutaBloques); err != nil {
			panic(err)
		}
	} else {
		log.Printf("Archivo encontrado: %s.", RutaBloques)
	}

	return nil
}

func crearBitmap(ruta string) error {
	tamañoBitmap := (globals.FSConfig.BlockCount + 7) / 8
	archivo, err := os.Create(ruta)
	if err != nil {
		return err
	}
	defer archivo.Close()

	_, err = archivo.Write(make([]byte, tamañoBitmap))
	if err != nil {
		return err
	}

	return nil
}

func crearBloques(ruta string) error {
	tamañoBloques := globals.FSConfig.BlockSize * globals.FSConfig.BlockCount
	archivo, err := os.Create(ruta)
	if err != nil {
		return err
	}
	defer archivo.Close()

	_, err = archivo.Write(make([]byte, tamañoBloques))
	if err != nil {
		return err
	}

	return nil
}

func CrearMetadata(ruta string, pid int, tid int, tamaño int, indexBlock int) (error, string) {
	timestamp := time.Now().Format("150405.000")
	nombreArchivo := fmt.Sprintf("%d-%d-%s.dmp", pid, tid, timestamp)

	rutaArchivo := filepath.Join(ruta, nombreArchivo)

	err := os.MkdirAll(ruta, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creando directorio %s: %w", ruta, err), ""
	}

	archivo, err := os.Create(rutaArchivo)
	if err != nil {
		return fmt.Errorf("error creando archivo de metadata %s: %w", rutaArchivo, err), ""
	}
	defer archivo.Close()

	metadata := globals.Metadata{
		IndexBlock: indexBlock,
		Size:       tamaño,
	}

	json, err := commons.CodificarJSON(metadata)
	if err != nil {
		return fmt.Errorf("error codificando metadata a JSON: %w", err), ""
	}

	_, err = archivo.Write(json)
	if err != nil {
		return fmt.Errorf("error escribiendo metadata en el archivo: %w", err), ""
	}

	return nil, nombreArchivo
}
