package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/servidor"
	"log"
	"net/http"
	"os"
	"path/filepath"

	//"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/config"
)

func main() {
	// Iniciar configuracion
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	globals.CConfig = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.Config{}).(*globals.Config)

	if globals.CConfig == nil {
		log.Fatalln("Error al cargar la configuración")
	}

	configs.ConfigurarLogger()

	cliente.EnviarMensaje(globals.CConfig.IpMemory, globals.CConfig.PortMemory, "hola, soy cpu")

	port := fmt.Sprintf(":%d", globals.CConfig.Port)

	log.Printf("El módulo memoria está a la escucha en el puerto %s", port)

	http.HandleFunc("GET /mensaje", servidor.RecibirMensaje)

	err = http.ListenAndServe(port, nil)

	if err != nil {
		panic(err)
	}

	//cliente.EnviarMensaje(globals.CConfig.IpKernel, globals.CConfig.PortKernel, "hola, soy cpu")
}
