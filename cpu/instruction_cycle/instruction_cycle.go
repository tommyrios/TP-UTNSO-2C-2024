package instruction_cycle

import (
	"log"
	"net/http"
	"strings"

	generalCPU "github.com/sisoputnfrba/tp-golang/cpu/general"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/instrucciones"
	"github.com/sisoputnfrba/tp-golang/cpu/interrupciones"
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

func RecibirInterrupcion(w http.ResponseWriter, r *http.Request) {
	var interrupcion commons.InterrupcionRecibida

	err := commons.DecodificarJSON(r.Body, &interrupcion)
	if err != nil {
		return
	}

}

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

	tcbUsado := pcbUsada.Tid[0]

	log.Printf("TID: %d - Solicito Contexto Ejecución", tcbUsado.Tid)

	*globals.Registros = tcbUsado.Registros
	*globals.Tid = pcbUsada.Tid[0].Tid
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
			break
		}

	}

	despacho.Pcb = pcbUsada
	despacho.Pcb.Tid[0] = tcbUsado
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

	log.Printf("TID: %d - FETCH - Program Counter: %d", *globals.Tid, globals.Registros.PC)

	return respuestaInstruccion.Instruccion

}

func Decode(instruccion string) {
	// Separar la instrucción en partes: Opcode y operandos
	parts := strings.Split(instruccion, " ")
	opCode := parts[0]
	operands := parts[1:]

	globals.Instruccion.CodigoInstruccion = opCode
	globals.Instruccion.Operandos = operands
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
		instrucciones.HandleSyscall(respuesta)
		continuarEjecucion = false
	default:
		continuarEjecucion = false
	}

	return continuarEjecucion, salto
}

func Interrupcion(respuesta *commons.DespachoProceso) bool {
	// Chequear si hubo una interrupción
	status, reason, tid := interrupciones.ObtenerYResetearInterrupcion()

	if status && tid == *globals.Tid {
		// Si hay interrupción => actualizar contexto en mem
		tcbActual := commons.TCB{
			Pid:       *globals.Pid,
			Tid:       *globals.Tid,
			Registros: *globals.Registros,
		}

		_, err := generalCPU.NotifyMemory(tcbActual)
		if err != nil {
			log.Fatal("Error al actualizar el contexto en memoria")
		}

		// Notificar al Kernel el motivo de la interrupción
		respuesta.Reason = reason
		respuesta.Pcb.Tid[0].Registros = *globals.Registros
		respuesta.Pcb.ProgramCounter = int(globals.Registros.PC)

		_, err = generalCPU.NotifyKernel(respuesta, "/syscall/interruption")
		if err != nil {
			log.Fatal("Error al notificar la interrupción al Kernel")
		}

		return true
	}

	return false
}
