package cliente

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func Post(ip string, port int, ruta string, jsonData []byte) *http.Response {
	url := fmt.Sprintf("http://%s:%d/%s", ip, port, ruta)

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	//log.Println("Respuesta POST:", response.Status)

	return response
}

func Post2(ip string, port int, ruta string, jsonData []byte) (*http.Response, []byte) {
	url := fmt.Sprintf("http://%s:%d/%s", ip, port, ruta)

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		panic(err)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil
	}

	//log.Println("Respuesta POST:", response.Status)

	return response, bodyBytes
}
