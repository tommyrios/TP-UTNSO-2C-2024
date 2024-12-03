package instruction_cycle

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/cpu/general/requests"
	"github.com/sisoputnfrba/tp-golang/cpu/instrucciones"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"log"
	"net/http"
	"strings"

	generalCPU "github.com/sisoputnfrba/tp-golang/cpu/general"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/interrupciones"
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
	log.Println("## Llega interrupcion al puerto Interrupt")
	var interrupcion globals.InterrupcionRecibida

	err := commons.DecodificarJSON(r.Body, &interrupcion)
	if err != nil {
		return
	}
}

// Similar a Recibir Mensaje
func Dispatch(w http.ResponseWriter, r *http.Request) {

	var req requests.RequestDispatch

	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		return
	}

	err = EjecutarInstruccion(req.Pid, req.Tid)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

/*func EjecutarInstrucciones(pcbUsada commons.PCB, tid int) {
	despacho := commons.DespachoProceso{Pcb: pcbUsada}
	tcbUsado := pcbUsada.Tid[tid]

	log.Printf("TID: %d - Solicito Contexto Ejecución", tid)

	*globals.Registros = tcbUsado.Registros
	*globals.Tid = tid
	*globals.Pid = pcbUsada.Pid
	globals.Registros.PC = uint32(pcbUsada.ProgramCounter)

	///////////////////////
	///CICLO INSTRUCCION///
	///////////////////////

	for {
		//Fetch: Recibe la instruccion
		instruccion := Fetch()

		//Decode: Traduce
		//Decode(instruccion)
		parts := strings.Split(instruccion, " ")
		opCode := parts[0]
		operands := parts[1:]

		instruccionStruct := globals.InstruccionStruct{Partes: parts, CodigoInstruccion: opCode, Operandos: operands}
		//Execute: Puede ejecutar y seguir, ejecutar y hacer un salto con JNZ o Terminar su ejecucion tras hecha la instruccion

		continuarEjecucion, saltoJNZ := Execute(&despacho, instruccionStruct)

		if !saltoJNZ {
			globals.Registros.PC++
		}

		//Check Interruption
		if !continuarEjecucion || Interrupcion(&despacho) {
			break
		}

	}

	log.Printf("## TID: %d - Actualizo Contexto Ejecución", *globals.Tid)

	resp, err := commons.CodificarJSON(despacho)
	if err != nil {
		return
	}

	cliente.Post(globals.CConfig.IpKernel, globals.CConfig.PortKernel, "pcb", resp)
}*/

func EjecutarInstruccion(pid int, tid int) error {
	contexto, err := solicitarContexto(pid, tid)

	if err != nil {
		return err
	}

	for {
		instruccionRecibida := Fetch(pid, tid, int(contexto.Registros.PC))

		instruccion := Decode(instruccionRecibida)

		Execute(instruccion, contexto.Registros)
	}
}

func solicitarContexto(pid int, tid int) (requests.ResponseContexto, error) {
	var reqContexto = requests.RequestContexto{Pid: pid, Tid: tid}

	var responseContexto requests.ResponseContexto

	reqCodificada, err := commons.CodificarJSON(reqContexto)

	if err != nil {
		return responseContexto, err
	}

	response, contexto := cliente.Post2(globals.CConfig.IpMemory, globals.CConfig.PortMemory, "contexto_de_ejecucion", reqCodificada)

	defer response.Body.Close()

	err = json.Unmarshal(contexto, &responseContexto)

	if err != nil {
		return responseContexto, err
	}

	return responseContexto, nil
}

func Fetch(pid int, tid int, pc int) string {
	reqPedidoInstruccion, err := commons.CodificarJSON(requests.RequestInstruccion{Pid: pid, Tid: tid, PC: pc})

	if err != nil {
		return ""
	}

	response, instruccion := cliente.Post2(globals.CConfig.IpMemory, globals.CConfig.PortMemory, "obtener_instruccion", reqPedidoInstruccion)

	defer response.Body.Close()

	return string(instruccion)
}

func Decode(instruccion string) globals.InstruccionStruct {
	partes := strings.Split(instruccion, " ")

	instruccionStruct := globals.InstruccionStruct{CodOperacion: partes[0], Operandos: partes[1:]}

	return instruccionStruct
}

func Execute(instruccion globals.InstruccionStruct, registros commons.Registros) {
	switch instruccion.CodOperacion {
	case SET:
		instrucciones.Set(instruccion.Operandos)
	case SUM:
		instrucciones.Sum(instruccion.Operandos)
	case SUB:
		instrucciones.Sub(instruccion.Operandos)
	case JNZ:
		salto = instrucciones.Jnz(instruccion.Operandos)
	case READ_MEM:
		instrucciones.ReadMem(instruccion.Operandos)
	case WRITE_MEM:
		instrucciones.WriteMem(instruccion.Operandos)
	case LOG:
		instrucciones.Log()
	case DUMP_MEMORY, IO, PROCESS_CREATE, THREAD_CREATE,
		THREAD_JOIN, THREAD_CANCEL, MUTEX_CREATE,
		MUTEX_LOCK, MUTEX_UNLOCK, THREAD_EXIT, PROCESS_EXIT:
		instrucciones.HandleSyscall(respuesta, &instruccion)
		continuarEjecucion = false
	default:
		continuarEjecucion = false
	}

}

/*func Execute(respuesta *commons.DespachoProceso, instruccion globals.InstruccionStruct) (bool, bool) {
	//Instrucciones que no requieren Decode:SET, SUM, SUB, JNZ, LOG.

	continuarEjecucion := true
	salto := false

	switch instruccion.CodigoInstruccion {
	case SET:
		instrucciones.Set()
	case SUM:
		instrucciones.Sum()
	case SUB:
		instrucciones.Sub()
	case JNZ:
		salto = instrucciones.Jnz()
	case READ_MEM:
		instrucciones.ReadMem()
	case WRITE_MEM:
		instrucciones.WriteMem()
	case LOG:
		instrucciones.Log()
	case DUMP_MEMORY, IO, PROCESS_CREATE, THREAD_CREATE,
		THREAD_JOIN, THREAD_CANCEL, MUTEX_CREATE,
		MUTEX_LOCK, MUTEX_UNLOCK, THREAD_EXIT, PROCESS_EXIT:
		instrucciones.HandleSyscall(respuesta, &instruccion)
		continuarEjecucion = false
	default:
		continuarEjecucion = false
	}

	log.Printf("## TID: %d - Ejecutando: %s - %v", *globals.Tid, globals.Instruccion.CodigoInstruccion, globals.Instruccion.Operandos)

	return continuarEjecucion, salto
}*/

func Interrupcion(respuesta *commons.DespachoProceso) bool {
	// Chequear si hubo una interrupción
	status, reason, tid := interrupciones.ObtenerYResetearInterrupcion()

	if status && tid == *globals.Tid {
		log.Printf("## TID: %d - Actualizo Contexto Ejecución", *globals.Tid)
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
