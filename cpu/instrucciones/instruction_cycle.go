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
)

var interrupcion = globals.InterrupcionStruct{Presencia: false}

func RecibirInterrupcion(w http.ResponseWriter, r *http.Request) {
	log.Println("## Llega interrupcion al puerto Interrupt")

	var interrupcionRecibida requests.RequestInterrupcion

	err := commons.DecodificarJSON(r.Body, &interrupcion)

	interrupcion.Pid = interrupcionRecibida.Pid
	interrupcion.Tid = interrupcionRecibida.Tid
	interrupcion.Presencia = true

	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
}

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
	return
}

func EjecutarInstruccion(pid int, tid int) error {
	contexto, err := solicitarContexto(pid, tid)

	var ultimaInstruccion string

	if err != nil {
		return err
	}

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

		if CheckInterrupt(pid, tid) == true {
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

	log.Printf("## TID: %d - Solicito Contexto Ejecución.\n", tid)

	return responseContexto, nil
}

func Fetch(pid int, tid int, pc int) string {
	reqPedidoInstruccion, err := commons.CodificarJSON(requests.RequestInstruccion{Pid: pid, Tid: tid, PC: pc})

	if err != nil {
		return ""
	}

	response, instruccion := cliente.Post2(globals.CConfig.IpMemory, globals.CConfig.PortMemory, "obtener_instruccion", reqPedidoInstruccion)

	defer response.Body.Close()

	log.Printf("## TID: %d - FETCH - Program Counter: %d.", tid, pc)

	return string(instruccion)
}

func Decode(instruccion string) globals.InstruccionStruct {
	partes := strings.Split(instruccion, " ")

	instruccionStruct := globals.InstruccionStruct{CodOperacion: partes[0], Operandos: partes[1:]}

	return instruccionStruct
}

func Execute(instruccion globals.InstruccionStruct, registros *commons.Registros, base int, limite int, pid int, tid int) int {

	log.Printf("## TID: %d - Ejecutando: %s - %s.", tid, instruccion.CodOperacion, instruccion.Operandos)

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
		Syscall(instruccion, registros, pid, tid)
		globals.DevolverPCB(pid, tid, "PROCESS_EXIT")
		return 1

	case "THREAD_EXIT":
		Syscall(instruccion, registros, pid, tid)
		globals.DevolverPCB(pid, tid, "THREAD_EXIT")
		return 1

	case "DUMP_MEMORY", "IO", "PROCESS_CREATE", "THREAD_CREATE",
		"THREAD_JOIN", "THREAD_CANCEL", "MUTEX_CREATE",
		"MUTEX_LOCK", "MUTEX_UNLOCK":
		Syscall(instruccion, registros, pid, tid)
		globals.DevolverPCB(pid, tid, "SYSCALL")

		return 1
	}

	return 4
}

func CheckInterrupt(pid int, tid int) bool {
	if interrupcion.Presencia == true {
		interrupcion.Presencia = false
		if interrupcion.Pid == pid && interrupcion.Tid == tid {
			globals.DevolverPCB(pid, tid, "INTERRUPT")
			return true
		}
	}

	return false
}
