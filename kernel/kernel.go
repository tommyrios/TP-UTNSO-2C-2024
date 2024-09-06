package main

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
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

	globals.KConfig = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.Config{}).(*globals.Config)

	if globals.KConfig == nil {
		log.Fatalln("Error al cargar la configuración")
	}

	cliente.EnviarMensaje(globals.KConfig.IpMemoria, globals.KConfig.PuertoMemoria, "hola")
}
