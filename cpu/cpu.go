package main

import (
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"log"
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

	log.Println(globals.CConfig.IpMemory)
}
