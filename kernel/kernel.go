package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	configs "github.com/sisoputnfrba/tp-golang/utils/config"
	"net/http"

	"log"
	"os"
	"path/filepath"
)

func main() {
	//// Configuración ////
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	globals.KConfig = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.Config{}).(*globals.Config)

	if globals.KConfig == nil {
		log.Fatalln("Error al cargar la configuración")
	}

	//// Logger ////
	configs.ConfigurarLogger("kernel")

	//// Conexión ////
	mux := http.NewServeMux()
	mux.HandleFunc("/mensaje", commons.RecibirMensaje)
	mux.HandleFunc("/process", globals.IniciarProceso)

	port := fmt.Sprintf(":%d", globals.KConfig.Port)

	log.Printf("El módulo kernel está a la escucha en el puerto %s", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}
