package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/config"
	"github.com/sisoputnfrba/tp-golang/utils/servidor"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Iniciar configuracion
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	globals.KConfig = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.Config{}).(*globals.Config)

	if globals.KConfig == nil {
		log.Fatalln("Error al cargar la configuración")
	}

	// Iniciar conexión

	cliente.EnviarMensaje(globals.KConfig.IpMemory, globals.KConfig.PortMemory, "hola, soy kernel")

	cliente.EnviarMensaje(globals.KConfig.IpCpu, globals.KConfig.PortCpu, "hola, soy kernel")

	port := fmt.Sprintf(":%d", globals.KConfig.Port)

	log.Printf("El módulo memoria está a la escucha en el puerto %s", port)

	http.HandleFunc("GET /mensaje", servidor.RecibirMensaje)

	err = http.ListenAndServe(port, nil)

	if err != nil {
		panic(err)
	}
}
