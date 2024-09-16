package instruction_cycle

import (
	//"log"
	"net/http"
	//"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

// Similar a Recibir Mensaje
func Ejecutar(w http.ResponseWriter, r *http.Request) {

	var pcbUsada commons.PCB

	err := commons.DecodificarJSON(r.Body, &pcbUsada)
	if err != nil {
		return
	}

	/////////////////////
	///CARGAR CONTEXTO///
	/////////////////////

	///////////////////////
	///CICLO INSTRUCCION///
	///////////////////////

}
