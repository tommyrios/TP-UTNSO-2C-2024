package instruction_cycle

import (
	"log"
	"net/http"

	generalCPU "github.com/sisoputnfrba/tp-golang/cpu/general"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/instrucciones"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
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

	go EjecutarInstrucciones(pcbUsada)
}

func EjecutarInstrucciones(pcbUsada commons.PCB) {
	var despacho commons.DespachoProceso

	log.Printf("TID: %d - Solicito Contexto Ejecución", pcbUsada.Tid)

	*globals.Registros = pcbUsada.Registros
	*globals.Pid = pcbUsada.Pid
	globals.Registros.PC = uint32(pcbUsada.ProgramCounter)

	///////////////////////
	///CICLO INSTRUCCION///
	///////////////////////

	for {

		//Fetch: Recibe la instruccion
		Fetch()

		//Decode: Traduce
		Decode()

		//Execute: Puede ejecutar y seguir, ejecutar y hacer un salto con JNZ o Terminar su ejecucion tras hecha la instruccion

		continuarEjecucion, saltoJNZ := Execute(&despacho)

		if !saltoJNZ {
			globals.Registros.PC++
		}

		//Check Interruption
		if !continuarEjecucion || Interrupcion(&despacho) {
			//Log de por que se cortó la ejecucion
			break
		}

	}

	despacho.Pcb = pcbUsada
	despacho.Pcb.Registros = *globals.Registros
	despacho.Pcb.ProgramCounter = int(globals.Registros.PC)

	resp, err := commons.CodificarJSON(despacho)
	if err != nil {
		return
	}

	cliente.Post(globals.CConfig.IpKernel, globals.CConfig.PortKernel, "pcb", resp)

}

func Fetch() {
	resp, err := generalCPU.ObtenerInstruction()

	if err != nil || resp == nil {
		log.Fatal("Error al buscar instruccion en memoria")
		return
	}
	var respuestaInstruccion commons.GetRespuestaInstruccion
	commons.DecodificarJSON(resp.Body, &respuestaInstruccion)

	log.Printf("TID: %d - FETCH - Program Counter: %d", *globals.Pid, globals.Registros.PC)

}

func Decode() {
	//TRADUCCION CON MMU (LOGICO A FISICO)
}

func Execute(respuesta *commons.DespachoProceso) (bool, bool) {
	//Instrucciones que no requieren Decode:SET, SUM, SUB, JNZ, LOG.

	continuarEjecucion := true
	salto := false

	switch globals.Instruccion.CodigoInstruccion {
	case "SET":
		instrucciones.Set()
	case "SUM":
		instrucciones.Sum()
	case "SUB":
		instrucciones.Sub()
	case "JNZ":
		salto = instrucciones.Jnz()
	default:
		continuarEjecucion = false

	}

	return continuarEjecucion, salto
}

func Interrupcion(respuesta *commons.DespachoProceso) bool {
	return true
}
