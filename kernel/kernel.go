package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/processes"
	"github.com/sisoputnfrba/tp-golang/kernel/globals/schedulers"
	"log/slog"
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
		slog.Debug(fmt.Sprintf("Error al cargar la configuración"))
	}

	//// Logger ////
	configs.ConfigurarLogger("kernel")

	//// Cola Ready ////
	schedulers.ManejarColaReady()

	//// Hilo Execute ////
	go schedulers.ManejarHiloRunning()

	//// Proceso Inicial ////
	processes.ProcesoInicial(os.Args)

	//go processes.CrearProcesoNew()

	//// Servidor ////
	mux := http.NewServeMux()
	mux.HandleFunc("/process_create", handlers.HandleProcessCreate)
	mux.HandleFunc("/thread_create", handlers.HandleThreadCreate)
	mux.HandleFunc("/process_exit", handlers.HandleProcessExit)
	mux.HandleFunc("/thread_exit", handlers.HandleThreadExit)
	mux.HandleFunc("/thread_join", handlers.HandleThreadJoin)
	mux.HandleFunc("/thread_cancel", handlers.HandleThreadCancel)
	mux.HandleFunc("/mutex_create", handlers.HandleMutexCreate)
	mux.HandleFunc("/mutex_lock", handlers.HandleMutexLock)
	mux.HandleFunc("/mutex_unlock", handlers.HandleMutexUnlock)
	mux.HandleFunc("/dump_memory", handlers.HandleDumpMemory)
	mux.HandleFunc("/handle_io", handlers.HandleIO)
	mux.HandleFunc("/compactacion", handlers.HandleCompactacion)
	mux.HandleFunc("/compactacion_finalizada", handlers.HandleCompactacionFinalizada)
	mux.HandleFunc("/pcb", handlers.HandleDesalojoCpu)

	port := fmt.Sprintf(":%d", globals.KConfig.Port)

	slog.Info(fmt.Sprintf("El módulo kernel está a la escucha en el puerto %s", port))

	globals.CpuLibre <- true

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}
