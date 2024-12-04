package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/instrucciones"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	configs "github.com/sisoputnfrba/tp-golang/utils/config"
)

func main() {
	//// Configuración  ////
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	//globals.Registros = new(commons.Registros)
	//globals.Pid = new(int)
	//globals.Tid = new(int)
	globals.CConfig = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.Config{}).(*globals.Config)

	if globals.CConfig == nil {
		log.Fatalln("Error al cargar la configuración")
	}

	//// Logger ////
	configs.ConfigurarLogger("cpu")

	//// Conexiones ////
	mux := http.NewServeMux()
	mux.HandleFunc("/dispatch", instrucciones.Dispatch)
	mux.HandleFunc("/interrupt", instrucciones.RecibirInterrupcion)

	port := fmt.Sprintf(":%d", globals.CConfig.Port)

	log.Printf("El módulo CPU está a la escucha en el puerto %s", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}
