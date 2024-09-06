package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"log"
	"net/http"
	"os"
	"path/filepath"
	//"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/config"
	"github.com/sisoputnfrba/tp-golang/utils/servidor"
)

func main() {
	// Iniciar configuracion
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	globals.MConfig = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.Config{}).(*globals.Config)

	if globals.MConfig == nil {
		log.Fatalln("Error al cargar la configuración")
	}

	mux := http.NewServeMux()

	port := fmt.Sprintf(":%d", globals.MConfig.Port)

	http.HandleFunc("GET /mensaje", servidor.RecibirMensaje)

	err2 := http.ListenAndServe(port, mux)
	
	if err2 != nil {
		panic(err2)
	}
}
