package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/instrucciones"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	configs "github.com/sisoputnfrba/tp-golang/utils/config"
)

func main() {
	//// Configuración ////
	path, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	globals.CConfig = configs.IniciarConfiguracion(filepath.Join(path, "config.json"), &globals.Config{}).(*globals.Config)

	if globals.CConfig == nil {
		slog.Debug(fmt.Sprintf("Error al cargar la configuración"))
	}

	//// Logger ////
	configs.ConfigurarLogger("cpu", globals.CConfig.LogLevel)

	//// Conexiones ////
	mux := http.NewServeMux()
	mux.HandleFunc("/dispatch", instrucciones.Dispatch)

	port := fmt.Sprintf(":%d", globals.CConfig.Port)

	slog.Info(fmt.Sprintf("El módulo CPU está a la escucha en el puerto %s", port))

	err = http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}
