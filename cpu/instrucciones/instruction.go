package instrucciones

import (
	"encoding/binary"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/globals/requests"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
	"strconv"
)

func Set(operandos []string, registros *commons.Registros) {
	valor, _ := strconv.Atoi(operandos[1])
	globals.CambiarValorRegistros(operandos[0], uint32(valor), registros)
}

func Sum(operandos []string, registros *commons.Registros) {
	valorRegistroDestino := globals.ValorRegistros(operandos[0], registros)
	valorRegistroOrigen := globals.ValorRegistros(operandos[1], registros)

	globals.CambiarValorRegistros(operandos[0], valorRegistroOrigen+valorRegistroDestino, registros)
}

func Sub(operandos []string, registros *commons.Registros) {
	valorRegistroDestino := globals.ValorRegistros(operandos[0], registros)
	valorRegistroOrigen := globals.ValorRegistros(operandos[1], registros)

	globals.CambiarValorRegistros(operandos[0], valorRegistroDestino-valorRegistroOrigen, registros)
}

func Jnz(operandos []string, registros *commons.Registros) {
	if globals.ValorRegistros(operandos[0], registros) != 0 {
		valor, _ := strconv.Atoi(operandos[1])
		registros.PC = uint32(valor)
	} else {
		registros.PC++
	}
}

func Log(operandos []string, registros *commons.Registros) {
	log.Printf("## Valor del registro %s: %d\n", operandos[0], globals.ValorRegistros(operandos[0], registros))
}

func ReadMem(operandos []string, registros *commons.Registros, base int, limite int, pid int, tid int) int {
	registroDireccion := globals.ValorRegistros(operandos[1], registros)

	direccionFísica, err := globals.Mmu(int(registroDireccion), base, limite)

	if err == 1 {
		globals.DevolverPCB(pid, tid, "SEGMENTATION FAULT")

		return 1
	}

	reqReadMemory := requests.RequestReadMemory{Direccion: direccionFísica, Pid: pid, Tid: tid}

	reqCodificada, _ := commons.CodificarJSON(reqReadMemory)

	response, byteSolicitados := cliente.Post2(globals.CConfig.IpMemory, globals.CConfig.PortMemory, "read_memory", reqCodificada)

	defer response.Body.Close()

	nuevoValor := binary.LittleEndian.Uint32(byteSolicitados)

	globals.CambiarValorRegistros(operandos[0], nuevoValor, registros)

	log.Printf("## TID: %d - Acción: LEER - Dirección Física: %d.", tid, direccionFísica)

	return 0
}

func WriteMem(operandos []string, registros *commons.Registros, base int, limite int, pid int, tid int) int {
	registroDireccion := globals.ValorRegistros(operandos[0], registros)

	direccionFísica, err := globals.Mmu(int(registroDireccion), base, limite)

	if err == 1 {
		globals.DevolverPCB(pid, tid, "SEGMENTATION FAULT")

		return 1
	}

	reqWriteMemory := requests.RequestWriteMemory{Direccion: direccionFísica, Pid: pid, Tid: tid, Datos: []byte{byte(globals.ValorRegistros(operandos[1], registros))}}

	reqCodificada, _ := commons.CodificarJSON(reqWriteMemory)

	response, mensaje := cliente.Post2(globals.CConfig.IpMemory, globals.CConfig.PortMemory, "write_memory", reqCodificada)

	defer response.Body.Close()

	log.Printf("Respuesta de memoria a escribir memoria: %s\n", string(mensaje))

	log.Printf("## TID: %d - Acción: ESCRIBIR - Dirección Física: %d.", tid, direccionFísica)

	return 0
}

func Syscall(instruccion globals.InstruccionStruct, registros *commons.Registros, pid int, tid int) {
	var parametros requests.RequestSyscall

	var ruta string

	switch instruccion.CodOperacion {

	case "PROCESS_CREATE":
		tamMemoria, _ := strconv.Atoi(instruccion.Operandos[1])
		prioridadTid, _ := strconv.Atoi(instruccion.Operandos[2])

		parametros = requests.RequestSyscall{
			Pid:            pid,
			Tid:            tid,
			PseudoCodigo:   instruccion.Operandos[0],
			TamanioMemoria: tamMemoria,
			Prioridad:      prioridadTid,
		}
		ruta = "process_create"

	case "PROCESS_EXIT":
		parametros = requests.RequestSyscall{
			Pid: pid,
			Tid: tid,
		}
		ruta = "process_exit"

	case "THREAD_CREATE":
		prioridadTid, _ := strconv.Atoi(instruccion.Operandos[1])

		parametros = requests.RequestSyscall{
			Pid:          pid,
			Tid:          tid,
			PseudoCodigo: instruccion.Operandos[0],
			Prioridad:    prioridadTid,
		}
		ruta = "thread_create"

	case "THREAD_JOIN":
		tidParametro, _ := strconv.Atoi(instruccion.Operandos[0])

		parametros = requests.RequestSyscall{
			Pid:          pid,
			Tid:          tid,
			TidParametro: tidParametro,
		}
		ruta = "thread_join"

	case "THREAD_CANCEL":
		tidEliminar, _ := strconv.Atoi(instruccion.Operandos[0])

		parametros = requests.RequestSyscall{
			Pid:          pid,
			Tid:          tid,
			TidAEliminar: tidEliminar,
		}
		ruta = "thread_cancel"

	case "THREAD_EXIT":
		parametros = requests.RequestSyscall{
			Pid: pid,
			Tid: tid,
		}
		ruta = "thread_exit"

	case "MUTEX_CREATE":
		parametros = requests.RequestSyscall{
			Pid:         pid,
			Tid:         tid,
			NombreMutex: instruccion.Operandos[0],
		}
		ruta = "mutex_create"

	case "MUTEX_LOCK":
		parametros = requests.RequestSyscall{
			Pid:         pid,
			Tid:         tid,
			NombreMutex: instruccion.Operandos[0],
		}
		ruta = "mutex_lock"

	case "MUTEX_UNLOCK":
		parametros = requests.RequestSyscall{
			Pid:         pid,
			Tid:         tid,
			NombreMutex: instruccion.Operandos[0],
		}
		ruta = "mutex_unlock"

	case "DUMP_MEMORY":
		parametros = requests.RequestSyscall{
			Pid: pid,
			Tid: tid,
		}

		ruta = "dump_memory"

	case "IO":

		tiempo, _ := strconv.Atoi(instruccion.Operandos[0])

		parametros = requests.RequestSyscall{
			Pid:    pid,
			Tid:    tid,
			Tiempo: tiempo,
		}
		ruta = "handle_io"
	}

	requestBody, _ := commons.CodificarJSON(parametros)

	cliente.Post(globals.CConfig.IpKernel, globals.CConfig.PortKernel, ruta, requestBody)
}

func EnviarRegistrosActualizados(registros *commons.Registros, pid int, tid int) {
	reqRegistrosActualizados := requests.RequestActualizarRegistros{
		Pid:       pid,
		Tid:       tid,
		Registros: registros,
	}

	reqCodificada, err := commons.CodificarJSON(reqRegistrosActualizados)

	if err != nil {
		return
	}

	cliente.Post(globals.CConfig.IpMemory, globals.CConfig.PortMemory, "actualizar_contexto", reqCodificada)

	//defer response.Body.Close()

	log.Printf("## TID: %d - Actualizo Contexto Ejecución.", tid)
	//log.Printf("Respuesta de memoria a actualizar contexto: %s\n", string(mensaje))
}
