package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
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

	globals.MConfig = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.Config{}).(*globals.Config)

	if globals.MConfig == nil {
		log.Fatalln("Error al cargar la configuración")
	}

	// Iniciar conexión

	port := fmt.Sprintf(":%d", globals.MConfig.Port)

	log.Printf("El módulo memoria está a la escucha en el puerto %s", port)

	http.HandleFunc("GET /mensaje", servidor.RecibirMensaje)

	err = http.ListenAndServe(port, nil)

	if err != nil {
		panic(err)
	}

	// cliente.EnviarMensaje(globals.MConfig.IpFileSystem, globals.MConfig.PortFileSystem, "hola, soy memoria")
	// cliente.EnviarMensaje(globals.MConfig.IpKernel, globals.MConfig.PortKernel, "hola, soy memoria")
	// cliente.EnviarMensaje(globals.MConfig.IpCpu, globals.MConfig.PortCpu, "hola, soy memoria")
}
