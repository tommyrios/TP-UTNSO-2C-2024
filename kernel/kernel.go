package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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
	pseudocodigo := os.Args[1]
	tamanio, _ := strconv.Atoi(os.Args[2])
	prioridadHiloMain, _ := strconv.Atoi(os.Args[3])

	globals.CrearProceso(pseudocodigo, tamanio, prioridadHiloMain)

	log.Println(pseudocodigo, tamanio, prioridadHiloMain)

	//// Conexión ////
	mux := http.NewServeMux()
	http.HandleFunc("/syscall/process_create", handlers.HandleProcessCreate)
	http.HandleFunc("/syscall/thread_create", handlers.HandleThreadCreate)
	http.HandleFunc("/syscall/process_exit", handlers.HandleProcessExit)
	http.HandleFunc("/syscall/thread_exit", handlers.HandleThreadExit)
	http.HandleFunc("/syscall/mutex_create", handlers.HandleMutexCreate)
	http.HandleFunc("/syscall/mutex_lock", handlers.HandleMutexLock)
	http.HandleFunc("/syscall/mutex_unlock", handlers.HandleMutexUnlock)

	port := fmt.Sprintf(":%d", globals.KConfig.Port)

	log.Printf("El módulo kernel está a la escucha en el puerto %s", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}
