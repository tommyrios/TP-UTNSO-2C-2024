package configs

import (
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
)

func IniciarConfiguracion(filePath string, moduloConfig interface{}) interface{} {
	configFile, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)

	if err := jsonParser.Decode(&moduloConfig); err != nil {
		log.Fatal(err.Error())
	}

	return moduloConfig
}

func ConfigurarLogger(modulo string, logLevel string) {
	logFile, err := os.OpenFile(modulo+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "INFO":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case "WARN":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case "ERROR":
		slog.SetLogLoggerLevel(slog.LevelError)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
