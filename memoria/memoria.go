package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"github.com/sisoputnfrba/tp-golang/utils/config"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	//// Configuración ////
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	globals.MConfig = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.Config{}).(*globals.Config)

	if globals.MConfig == nil {
		log.Fatalln("Error al cargar la configuración")
	}

	//// Logger ////
	configs.ConfigurarLogger("memoria")

	//// Conexión ////
	mux := http.NewServeMux()
	mux.HandleFunc("/mensaje", commons.RecibirMensaje)

	port := fmt.Sprintf(":%d", globals.MConfig.Port)

	log.Printf("El módulo memoria está a la escucha en el puerto %s", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}
