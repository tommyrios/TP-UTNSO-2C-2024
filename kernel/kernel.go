package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	configs "github.com/sisoputnfrba/tp-golang/utils/config"

	"log"
	"os"
	"path/filepath"
)

func main() {
	//// Configuración ////
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	globals.KConfig = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.Config{}).(*globals.Config)

	if globals.KConfig == nil {
		log.Fatalln("Error al cargar la configuración")
	}

	//// Logger ////
	configs.ConfigurarLogger("kernel")

	//// Proceso Inicial ////
	pseudocodigo := os.Args[1]
	tamanio, _ := strconv.Atoi(os.Args[2])
	prioridadHiloMain, _ := strconv.Atoi(os.Args[3])

	globals.CrearProceso(pseudocodigo, tamanio, prioridadHiloMain)

	log.Println(pseudocodigo, tamanio, prioridadHiloMain)

	globals.CrearProceso("pepito", 30, 0)
	//// Conexión ////
	mux := http.NewServeMux()
	mux.HandleFunc("POST /mensaje", commons.RecibirMensaje)
	//mux.HandleFunc("POST /process", globals.IniciarProceso)

	port := fmt.Sprintf(":%d", globals.KConfig.Port)

	log.Printf("El módulo kernel está a la escucha en el puerto %s", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}
