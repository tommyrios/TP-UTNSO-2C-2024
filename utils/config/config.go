package configs

import (
	"encoding/json"
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
