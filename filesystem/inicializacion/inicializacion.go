package inicializacion

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sisoputnfrba/tp-golang/filesystem/globals"
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
