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

const (
	SET            = "SET"
	SUM            = "SUM"
	SUB            = "SUB"
	JNZ            = "JNZ"
	READ_MEM       = "READ_MEM"
	WRITE_MEM      = "WRITE_MEM"
	LOG            = "LOG"
	DUMP_MEMORY    = "DUMP_MEMORY"
	IO             = "IO"
	PROCESS_CREATE = "PROCESS_CREATE"
	THREAD_CREATE  = "THREAD_CREATE"
	THREAD_JOIN    = "THREAD_JOIN"
	THREAD_CANCEL  = "THREAD_CANCEL"
	MUTEX_CREATE   = "MUTEX_CREATE"
	MUTEX_LOCK     = "MUTEX_LOCK"
	MUTEX_UNLOCK   = "MUTEX_UNLOCK"
	THREAD_EXIT    = "THREAD_EXIT"
	PROCESS_EXIT   = "PROCESS_EXIT"
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
		instruccion := Fetch()

		//Decode: Traduce
		Decode(instruccion)

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

func Fetch() string {
	resp, err := generalCPU.ObtenerInstruction()

	if err != nil || resp == nil {
		log.Fatal("Error al buscar instruccion en memoria")
		return "ERROR"
	}
	var respuestaInstruccion commons.GetRespuestaInstruccion
	commons.DecodificarJSON(resp.Body, &respuestaInstruccion)

	log.Printf("TID: %d - FETCH - Program Counter: %d", *globals.Pid, globals.Registros.PC)

	return respuestaInstruccion.Instruccion

}

func Decode(instruccion string) {
	//TRADUCCION CON MMU (LOGICO A FISICO)
}

func Execute(respuesta *commons.DespachoProceso) (bool, bool) {
	//Instrucciones que no requieren Decode:SET, SUM, SUB, JNZ, LOG.

	continuarEjecucion := true
	salto := false

	switch globals.Instruccion.CodigoInstruccion {
	case SET:
		instrucciones.Set()
	case SUM:
		instrucciones.Sum()
	case SUB:
		instrucciones.Sub()
	case JNZ:
		salto = instrucciones.Jnz()
	case READ_MEM:
		//instrucciones.ReadMem()
	case WRITE_MEM:
		//instrucciones.WriteMem()
	case LOG:
		instrucciones.Log()
	case DUMP_MEMORY, IO, PROCESS_CREATE, THREAD_CREATE,
		THREAD_JOIN, THREAD_CANCEL, MUTEX_CREATE,
		MUTEX_LOCK, MUTEX_UNLOCK, THREAD_EXIT, PROCESS_EXIT:
		generalCPU.HandleSyscall(respuesta)
		continuarEjecucion = false
	default:
		continuarEjecucion = false
	}

	return continuarEjecucion, salto
}

func Interrupcion(respuesta *commons.DespachoProceso) bool {
	return true
}
