package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/instruction_cycle"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	configs "github.com/sisoputnfrba/tp-golang/utils/config"
)

func main() {
	//// Configuración  ////
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	globals.Registros = new(commons.Registros)
	globals.Pid = new(int)
	globals.CConfig = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.Config{}).(*globals.Config)

	if globals.CConfig == nil {
		log.Fatalln("Error al cargar la configuración")
	}

	//// Logger ////
	configs.ConfigurarLogger("cpu")

	//// Conexiones ////
	mux := http.NewServeMux()
	mux.HandleFunc("POST /dispatch", instruction_cycle.Ejecutar)

	port := fmt.Sprintf(":%d", globals.CConfig.Port)

	log.Printf("El módulo memoria está a la escucha en el puerto %s", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}
