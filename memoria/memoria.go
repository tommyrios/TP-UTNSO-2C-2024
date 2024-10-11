package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/handlers"
	"github.com/sisoputnfrba/tp-golang/utils/config"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	//// Configuración ////
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	globals.MConfig = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.Config{}).(*globals.Config)

	if globals.MConfig == nil {
		log.Fatalln("Error al cargar la configuración")
	}

	//// Logger ////
	configs.ConfigurarLogger("memoria")

	//// Inicialización ////
	globals.InicializarMemoria()

	//// Conexión ////
	mux := http.NewServeMux()
	mux.HandleFunc("/contexto_de_ejecucion", handlers.HandleDevolverContexto)
	mux.HandleFunc("/actualizar_contexto", handlers.HandleActualizarContexto)
	mux.HandleFunc("/process", handlers.HandleSolicitudProceso)

	port := fmt.Sprintf(":%d", globals.MConfig.Port)

	log.Printf("El módulo memoria está a la escucha en el puerto %s", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}
