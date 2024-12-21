package configs

import (
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"
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

func ConfigurarLogger(modulo string) {
	logFile, err := os.OpenFile(modulo+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	slog.SetLogLoggerLevel(slog.LevelDebug)
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
