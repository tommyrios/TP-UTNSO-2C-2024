package handlers

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/globals/functions"
	"github.com/sisoputnfrba/tp-golang/memoria/globals/schemes"
	request3 "github.com/sisoputnfrba/tp-golang/memoria/handlers/request"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"net/http"
	"time"
)

// ¡¡¡¡¡HANDLERS CPU!!!!!
//Agregar retardo en peticiones!!

func HandleDevolverContexto(w http.ResponseWriter, r *http.Request) {
	var request request3.RequestContexto

	time.Sleep(time.Duration(globals.MConfig.ResponseDelay) * time.Millisecond)

	err := commons.DecodificarJSON(r.Body, &request)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	registros := functions.ObtenerRegistros(request.Pid, request.Tid)

	var response request3.ResponseContexto
	response.Registros = registros
	response.Pid = request.Pid
	response.Tid = request.Tid

	base, limite := functions.ObtenerBaseLimite(request.Pid, request.Tid) // ¡¡¡Falta implementar esta función!!!!

	response.Base = base
	response.Limite = limite

	responseCodificada, err := commons.CodificarJSON(response)

	if err != nil {
		http.Error(w, "Error al codificar el JSON", http.StatusBadRequest)
		return
	}

	cliente.Post(globals.MConfig.IpCpu, globals.MConfig.PortCpu, "/contexto_de_ejecucion", responseCodificada)
}

func HandleActualizarContexto(w http.ResponseWriter, r *http.Request) {
	var request request3.RequestActualizarContexto

	time.Sleep(time.Duration(globals.MConfig.ResponseDelay) * time.Millisecond)

	err := commons.DecodificarJSON(r.Body, &request)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	err = functions.ActualizarRegistros(request.Pid, request.Tid, request.Registros)

	if err != nil {
		http.Error(w, "Error actualizando los registros", http.StatusInternalServerError)
		return
	}

	// Responder con éxito si se actualizaron correctamente
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func HandleObtenerInstruccion(w http.ResponseWriter, r *http.Request) {
	var request request3.RequestObtenerInstruccion

	time.Sleep(time.Duration(globals.MConfig.ResponseDelay) * time.Millisecond)

	err := commons.DecodificarJSON(r.Body, &request)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	instruccion, err := functions.ObtenerInstruccion(request.Pid, request.Tid, request.PC)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener la instrucción: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	responseCodificada, err := commons.CodificarJSON(instruccion)

	if err != nil {
		http.Error(w, "Error al codificar el JSON", http.StatusBadRequest)
		return
	}

	cliente.Post(globals.MConfig.IpCpu, globals.MConfig.PortCpu, "/instruccion", responseCodificada)
}

func HandleReadMemory(w http.ResponseWriter, r *http.Request) {
	var request request3.RequestMemory

	time.Sleep(time.Duration(globals.MConfig.ResponseDelay) * time.Millisecond)

	err := commons.DecodificarJSON(r.Body, &request)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	response, err := functions.LeerMemoria(request.Byte)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error al leer la memoria: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	responseCodificada, err := commons.CodificarJSON(response)

	if err != nil {
		http.Error(w, "Error al codificar el JSON", http.StatusBadRequest)
		return
	}

	cliente.Post(globals.MConfig.IpCpu, globals.MConfig.PortCpu, "/lectura_memoria", responseCodificada)
}

func HandleWriteMemory(w http.ResponseWriter, r *http.Request) {
	var request request3.RequestMemory

	time.Sleep(time.Duration(globals.MConfig.ResponseDelay) * time.Millisecond)

	err := commons.DecodificarJSON(r.Body, &request)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	err = functions.EscribirMemoria(request.Byte, request.Pid, request.Datos)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error al escribir la memoria: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// ¡¡¡¡¡HANDLERS KERNEL!!!!!

func HandleCrearHilo(w http.ResponseWriter, r *http.Request) {
	var req request3.RequestCrearHilo

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	functions.CrearHilo(req.Pid, req.Tid, req.Pseudocodigo)
}

func HandleSolicitarProceso(w http.ResponseWriter, r *http.Request) {
	var req request3.RequestProcesoMemoria

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Lógica de asignación de espacio
	err = schemes.AsignarParticion(req.Pid, req.TamanioMemoria)

	w.WriteHeader(http.StatusOK)
}

func HandleFinalizarProceso(w http.ResponseWriter, r *http.Request) {
	var req request3.RequestFinalizarProceso

	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	_, existe := globals.MemoriaSistema.TablaProcesos[req.Pid]
	if !existe {
		http.Error(w, "Proceso no encontrado", http.StatusNotFound)
		return
	}

	// Liberamos la partición
	err = functions.LiberarProceso(req.Pid)

	// Eliminar las estructuras correspondientes del proceso en la Memoria del Sistema
	delete(globals.MemoriaSistema.TablaProcesos, req.Pid)
	delete(globals.MemoriaSistema.TablaHilos, req.Pid)
	delete(globals.MemoriaSistema.Pseudocodigos, req.Pid)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func HandleFinalizarHilo(w http.ResponseWriter, r *http.Request) {
	var req request3.RequestFinalizarHilo

	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	_, existe := globals.MemoriaSistema.TablaHilos[req.Pid][req.Tid]
	if !existe {
		http.Error(w, "Hilo no encontrado", http.StatusNotFound)
		return
	}

	// Eliminar las estructuras correspondientes del hilo en la Memoria del Sistema
	delete(globals.MemoriaSistema.TablaHilos[req.Pid], req.Tid)
	delete(globals.MemoriaSistema.Pseudocodigos[req.Pid], req.Tid)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func HandleMemoryDump(w http.ResponseWriter, r *http.Request) {
	var req request3.RequestDumpMemory

	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}
	base, limite := functions.ObtenerBaseLimite(req.Pid, req.Tid)
	// Obtener el contenido de la memoria del proceso
	TamanioMemoriaProceso := functions.ObtenerTamanioMemoria(base, limite)
	ContenidoProceso := functions.ObtenerContenidoMemoria(base, limite)

	// Solicitar al FileSystem la creación del archivo y escribir el contenido

	solicitud := request3.DumpMemoryFS{
		Pid:       req.Pid,
		Tid:       req.Tid,
		Tamanio:   TamanioMemoriaProceso,
		Contenido: ContenidoProceso,
	}

	solicitudCodificada, err := commons.CodificarJSON(solicitud)

	if err != nil {
		http.Error(w, "Error al codificar JSON", http.StatusBadRequest)
	}

	response := cliente.Post(globals.MConfig.IpFileSystem, globals.MConfig.PortFileSystem, "/dump_memory", solicitudCodificada)

	if response != nil && response.StatusCode == http.StatusOK {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	} else {
		http.Error(w, "Error al solicitar el dump de memoria al FileSystem", http.StatusInternalServerError)
	}
}
