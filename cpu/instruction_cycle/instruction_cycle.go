package instruction_cycle

import (
	"log"
	"net/http"

	generalCPU "github.com/sisoputnfrba/tp-golang/cpu/general"
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

func EjecutarInstrucciones(pcbUsada *commons.PCB) {

	*globals.Registros = pcbUsada.Registros
	*globals.Pid = pcbUsada.Pid
	globals.Registros.PC = uint32(pcbUsada.ProgramCounter)

	///////////////////////
	///CICLO INSTRUCCION///
	///////////////////////

	for {

		//Fetch
		fetch()
		//Decode
		decode()
		//Execute
		execute()
		//Check Interruption
		if Interrupcion() {
			break
		}

	}

}

func fetch() {
	resp, err := generalCPU.ObtenerInstruction()

	if err != nil || resp == nil {
		log.Fatal("Error al buscar instruccion en memoria")
		return
	}

}

func decode() {

}

func execute() {

}

func Interrupcion() bool {
	return true
}
