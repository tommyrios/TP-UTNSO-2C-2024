package main

import (
	//"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/config"
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
}
