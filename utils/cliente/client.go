package cliente

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Post(ip string, port int, ruta string, jsonData []byte) *http.Response {
	url := fmt.Sprintf("http://%s:%d/%s", ip, port, ruta)

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	log.Println("Respuesta POST:", string(body))

	return response
}

func Get(ip string, port int, ruta string) *http.Response {
	url := fmt.Sprintf("https://%s:%d/%s", ip, port, ruta)

	response, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	log.Println("Respuesta GET:", string(body))

	return response
}
