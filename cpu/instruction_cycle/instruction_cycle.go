package instruction_cycle

import (
	//"log"
	"net/http"
	//"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

// Similar a Recibir Mensaje
func Ejecutar(w http.ResponseWriter, r *http.Request) {

	var pcbUsada commons.PCB

	err := commons.DecodificarJSON(r.Body, &pcbUsada)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Pcb OK"))

	go EjecutarInstrucciones(&pcbUsada)
}

/////////////////////
///CARGAR CONTEXTO///
/////////////////////

func EjecutarInstrucciones(pcbUsada *commons.PCB) {

	*globals.Registros = pcbUsada.Registros
	*globals.Pid = pcbUsada.Pid
	globals.Registros.PC = uint32(pcbUsada.ProgramCounter)

	///////////////////////
	///CICLO INSTRUCCION///
	///////////////////////

	for {

		//Fetch

		//Decode

		//Execute

		//Check Interruption

	}

}
