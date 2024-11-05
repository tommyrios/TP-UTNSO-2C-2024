package handlers

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
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

	registros := globals.ObtenerRegistros(request.Pid, request.Tid)

	var response request3.ResponseContexto
	response.Registros = registros
	response.Pid = request.Pid
	response.Tid = request.Tid

	base, limite := globals.ObtenerBaseLimite(request.Pid, request.Tid) // ¡¡¡Falta implementar esta función!!!!

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

	err = globals.ActualizarRegistros(request.Pid, request.Tid, request.Registros)

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

	instruccion, err := globals.ObtenerInstruccion(request.Pid, request.Tid, request.PC)

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

	response, err := globals.LeerMemoria(request.Byte)

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

	err = globals.EscribirMemoria(request.Byte, request.Pid, request.Datos)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error al escribir la memoria: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// ¡¡¡¡¡HANDLERS KERNEL!!!!!

func HandleSolicitarProceso(w http.ResponseWriter, r *http.Request) {
	var req request3.RequestProcesoMemoria

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	esquemaFijo := globals.MConfig.Scheme == "fijo"

	// Lógica de asignación de espacio
	if esquemaFijo {
		if !globals.AsignarParticionFija(req.Pid, req.TamanioMemoria) {
			http.Error(w, "No hay espacio en particiones fijas", http.StatusConflict)
			return
		}
	} else {
		if !globals.AsignarParticionDinamica(req.Pid, req.TamanioMemoria) {
			http.Error(w, "No hay espacio en particiones dinámicas, compactación requerida", http.StatusConflict)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
