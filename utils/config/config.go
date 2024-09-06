package configs

import (
	"encoding/json"
	"io"
	"log"
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

func ConfigurarLogger() {
	logFile, err := os.OpenFile("tp0.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
