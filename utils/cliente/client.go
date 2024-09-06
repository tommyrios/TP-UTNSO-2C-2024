package cliente

import (
	"fmt"
	"io"
	_ "io"
	"net/http"
)

func EnviarMensaje(ip string, puerto int, mensajeTxt string) {
	cliente := &http.Client{}
	url := fmt.Sprintf("http://%s:%d/mensaje", ip, puerto)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return
	}

	q := req.URL.Query()
	q.Add("mensaje", mensajeTxt)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json")
	respuesta, err := cliente.Do(req)
	if err != nil {
		return
	}

	// Verificar el código de estado de la respuesta
	if respuesta.StatusCode != http.StatusOK {
		return
	}

	bodyBytes, err := io.ReadAll(respuesta.Body)
	if err != nil {
		return
	}

	fmt.Println(string(bodyBytes))
}
