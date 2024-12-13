package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sisoputnfrba/tp-golang/filesystem/globals"
	"github.com/sisoputnfrba/tp-golang/filesystem/handlers"
	"github.com/sisoputnfrba/tp-golang/filesystem/inicializacion"

	configs "github.com/sisoputnfrba/tp-golang/utils/config"
)

func main() {
	//// Configuracion ////
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	globals.FSConfig = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.Config{}).(*globals.Config)

	if globals.FSConfig == nil {
		log.Fatalln("Error al cargar la configuración")
	}

	//// Logger ////
	configs.ConfigurarLogger("filesystem")

	//// Inicialización ////

	err = inicializacion.IniciarFS(globals.FSConfig.MountDir)

	if err != nil {
		log.Fatalf("Error al inicializar el File System: %v", err)
	}
	log.Println("Inicialización del File System completada.")

	//// Conexión ////
	mux := http.NewServeMux()
	mux.HandleFunc("/memory_dump", handlers.CrearArchivo)

	port := fmt.Sprintf(":%d", globals.FSConfig.Port)

	log.Printf("El módulo filesystem está a la escucha en el puerto %s", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}
