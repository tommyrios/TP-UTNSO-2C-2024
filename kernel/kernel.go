package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/handlers"
	configs "github.com/sisoputnfrba/tp-golang/utils/config"
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
	globals.ProcesoInicial(os.Args)

	//// Rutinas ////
	globals.CpuLibre <- true
	go globals.ManejarColaReady()

	mux := http.NewServeMux()
	http.HandleFunc("/syscall/thread_exit", handlers.HandleThreadExit)
	http.HandleFunc("/syscall/thread_join", handlers.HandleThreadJoin)
	http.HandleFunc("/syscall/thread_cancel", handlers.HandleThreadCancel)
	http.HandleFunc("/syscall/mutex_create", handlers.HandleMutexCreate)
	http.HandleFunc("/syscall/mutex_lock", handlers.HandleMutexLock)
	http.HandleFunc("/syscall/mutex_unlock", handlers.HandleMutexUnlock)
	http.HandleFunc("/syscall/dump_memory", handlers.HandleDumpMemory)
	http.HandleFunc("/syscall/io", handlers.HandleIO)

	port := fmt.Sprintf(":%d", globals.KConfig.Port)

	log.Printf("El módulo kernel está a la escucha en el puerto %s", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}
