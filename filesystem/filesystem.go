package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/filesystem/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"github.com/sisoputnfrba/tp-golang/utils/config"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	//// Conexión ////
	mux := http.NewServeMux()
	mux.HandleFunc("/mensaje", commons.RecibirMensaje)

	port := fmt.Sprintf(":%d", globals.FSConfig.Port)

	log.Printf("El módulo filesystem está a la escucha en el puerto %s", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}
