package instrucciones

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/globals/requests"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
	"net/http"
	"strings"
	"time"
)

func Dispatch(w http.ResponseWriter, r *http.Request) {

	var req requests.RequestDispatch

	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		return
	}

	err = EjecutarInstruccion(req.Pid, req.Tid, req.Quantum, req.Scheduler)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func EjecutarInstruccion(pid int, tid int, quantum int, scheduler string) error {
	contexto, err := solicitarContexto(pid, tid)

	var ultimaInstruccion string

	if err != nil {
		return err
	}

	inicio := time.Now()

	for {
		instruccionRecibida := Fetch(pid, tid, int(contexto.Registros.PC))

		instruccion := Decode(instruccionRecibida)

		if Execute(instruccion, contexto.Registros, contexto.Base, contexto.Limite, pid, tid) == 1 {
			contexto.Registros.PC++
			ultimaInstruccion = instruccion.CodOperacion
			break
		}

		if instruccion.CodOperacion != "JNZ" {
			contexto.Registros.PC++
		}

		tiempo := int(time.Since(inicio).Milliseconds())

		if checkQuantum(pid, tid, scheduler, quantum, tiempo) {
			break
		}
	}

	if ultimaInstruccion == "THREAD_EXIT" || ultimaInstruccion == "PROCESS_EXIT" {
		return nil
	}

	EnviarRegistrosActualizados(contexto.Registros, pid, tid)

	return nil
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

	log.Printf("## PID %d - TID: %d - Solicito Contexto Ejecución.\n", pid, tid)

	return responseContexto, nil
}

func Fetch(pid int, tid int, pc int) string {
	reqPedidoInstruccion, err := commons.CodificarJSON(requests.RequestInstruccion{Pid: pid, Tid: tid, PC: pc})

	if err != nil {
		return ""
	}

	response, instruccion := cliente.Post2(globals.CConfig.IpMemory, globals.CConfig.PortMemory, "obtener_instruccion", reqPedidoInstruccion)

	defer response.Body.Close()

	log.Printf("## PID: %d TID: %d - FETCH - Program Counter: %d.", pid, tid, pc)

	return string(instruccion)
}

func Decode(instruccion string) globals.InstruccionStruct {
	partes := strings.Split(instruccion, " ")

	instruccionStruct := globals.InstruccionStruct{CodOperacion: partes[0], Operandos: partes[1:]}

	return instruccionStruct
}

func Execute(instruccion globals.InstruccionStruct, registros *commons.Registros, base int, limite int, pid int, tid int) int {

	log.Printf("## PID: %d TID: %d - Ejecutando: %s - %s.", pid, tid, instruccion.CodOperacion, instruccion.Operandos)

	switch instruccion.CodOperacion {
	case "SET":
		Set(instruccion.Operandos, registros)
		return 0
	case "SUM":
		Sum(instruccion.Operandos, registros)
		return 0
	case "SUB":
		Sub(instruccion.Operandos, registros)
		return 0
	case "JNZ":
		Jnz(instruccion.Operandos, registros)
		return 0
	case "LOG":
		Log(instruccion.Operandos, registros)
		return 0
	case "READ_MEM":
		return ReadMem(instruccion.Operandos, registros, base, limite, pid, tid)

	case "WRITE_MEM":
		return WriteMem(instruccion.Operandos, registros, base, limite, pid, tid)

	case "PROCESS_EXIT":
		Syscall(instruccion, pid, tid)
		globals.DevolverPCB(pid, tid, "PROCESS_EXIT")
		return 1

	case "THREAD_EXIT":
		Syscall(instruccion, pid, tid)
		globals.DevolverPCB(pid, tid, "THREAD_EXIT")
		return 1

	case "DUMP_MEMORY":
		Syscall(instruccion, pid, tid)
		globals.DevolverPCB(pid, tid, "MEMORY_DUMP")
		return 1

	case "IO", "THREAD_JOIN", "MUTEX_LOCK":
		Syscall(instruccion, pid, tid)
		globals.DevolverPCB(pid, tid, "SYSCALL")
		return 1

	case "PROCESS_CREATE", "THREAD_CREATE", "THREAD_CANCEL", "MUTEX_CREATE", "MUTEX_UNLOCK":
		Syscall(instruccion, pid, tid)
		return 0
	}

	return 4
}

func checkQuantum(pid int, tid int, scheduler string, quantum int, tiempo int) bool {
	if scheduler == "CMN" {
		if tiempo > quantum {
			log.Printf("END OF QUANTUM - (PID:TID) - (%d:%d)", pid, tid)
			globals.DevolverPCB(pid, tid, "END_OF_QUANTUM")
			return true
		}
	}

	return false
}
