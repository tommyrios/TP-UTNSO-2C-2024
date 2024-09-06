package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/utils/servidor"
	"net/http"

	//"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/filesystem/globals"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/config"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Iniciar configuracion
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	globals.FSConfig = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.Config{}).(*globals.Config)

	if globals.FSConfig == nil {
		log.Fatalln("Error al cargar la configuración")
	}

	cliente.EnviarMensaje(globals.FSConfig.IpMemory, globals.FSConfig.PortMemory, "hola, soy filesystem")

	port := fmt.Sprintf(":%d", globals.FSConfig.Port)

	log.Printf("El módulo memoria está a la escucha en el puerto %s", port)

	http.HandleFunc("GET /mensaje", servidor.RecibirMensaje)

	err = http.ListenAndServe(port, nil)

	if err != nil {
		panic(err)
	}
}
