package handlers

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/globals/functions"
	"github.com/sisoputnfrba/tp-golang/memoria/globals/schemes"
	"github.com/sisoputnfrba/tp-golang/memoria/handlers/requests"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
	"net/http"
	"time"
)

// ¡¡¡¡¡HANDLERS CPU!!!!!

func HandleDevolverContexto(w http.ResponseWriter, r *http.Request) {
	var request requests.RequestContexto

	time.Sleep(time.Duration(globals.MConfig.ResponseDelay) * time.Millisecond)

	err := commons.DecodificarJSON(r.Body, &request)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	registros := functions.ObtenerRegistros(request.Pid, request.Tid)

	var response requests.ResponseContexto
	response.Registros = registros
	response.Pid = request.Pid
	response.Tid = request.Tid

	base, limite := functions.ObtenerBaseLimite(request.Pid)

	response.Base = base
	response.Limite = limite

	responseCodificada, err := commons.CodificarJSON(response)

	if err != nil {
		http.Error(w, "Error al codificar el JSON", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseCodificada)
	if err != nil {
		return
	}
	log.Printf("## Contexto solicitado - (PID:TID) - (%d:%d)\n", request.Pid, request.Tid)

}

func HandleActualizarContexto(w http.ResponseWriter, r *http.Request) {
	var request requests.RequestActualizarContexto

	time.Sleep(time.Duration(globals.MConfig.ResponseDelay) * time.Millisecond)

	err := commons.DecodificarJSON(r.Body, &request)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	err = functions.ActualizarRegistros(request.Pid, request.Tid, request.Registros)

	if err != nil {
		http.Error(w, "Error actualizando los registros", http.StatusInternalServerError)
		log.Printf("Error al actualizar los registros - (PID:TID) - (%d:%d)\n", request.Pid, request.Tid)
		return
	}

	// Responder con éxito si se actualizaron correctamente
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("OK"))
	if err != nil {
		return
	}
	log.Printf("## Contexto actualizado - (PID:TID) - (%d:%d)\n", request.Pid, request.Tid)
}

func HandleEnviarInstruccion(w http.ResponseWriter, r *http.Request) {
	var request requests.RequestObtenerInstruccion

	time.Sleep(time.Duration(globals.MConfig.ResponseDelay) * time.Millisecond)

	err := commons.DecodificarJSON(r.Body, &request)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	instruccion, err := functions.ObtenerInstruccion(request.Pid, request.Tid, int(request.PC))

	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener la instrucción: %s", err.Error()), http.StatusInternalServerError)
		log.Printf("Error al obtener la instrucción - (PID:TID) - (%d:%d)\n", request.Pid, request.Tid)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(instruccion.Instruccion))
	if err != nil {
		return
	}
	log.Printf("## Obtener instrucción - (PID:TID) - (%d:%d) - Instrucción: %s\n", request.Pid, request.Tid, instruccion)
}

func HandleReadMemory(w http.ResponseWriter, r *http.Request) {
	var request requests.RequestReadMemory

	time.Sleep(time.Duration(globals.MConfig.ResponseDelay) * time.Millisecond)

	err := commons.DecodificarJSON(r.Body, &request)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	response, err := functions.LeerMemoria(request.Direccion, request.Pid)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error al leer la memoria: %s", err.Error()), http.StatusInternalServerError)
		log.Printf("Error al leer la memoria - (PID:TID) - (%d:%d)\n", request.Pid, request.Tid)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		return
	}
	log.Printf("## Lectura - (PID:TID) - (%d:%d) - Dir. Física: %d", request.Pid, request.Tid, request.Direccion)
}

func HandleWriteMemory(w http.ResponseWriter, r *http.Request) {
	var request requests.RequestWriteMemory

	time.Sleep(time.Duration(globals.MConfig.ResponseDelay) * time.Millisecond)

	err := commons.DecodificarJSON(r.Body, &request)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		log.Printf("Error al escribir la memoria - (PID:TID) - (%d:%d)\n", request.Pid, request.Tid)
		return
	}

	err = functions.EscribirMemoria(request.Direccion, request.Pid, request.Datos)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error al escribir la memoria: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("OK"))
	if err != nil {
		return
	}
	log.Printf("## Escritura - (PID:TID) - (%d:%d) - Dir. Física: %d", request.Pid, request.Tid, request.Direccion)
}

// ¡¡¡¡¡HANDLERS KERNEL!!!!!

func HandleCrearProceso(w http.ResponseWriter, r *http.Request) {
	var procesoRequest requests.RequestProcesoMemoria

	err := commons.DecodificarJSON(r.Body, &procesoRequest)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	if procesoRequest.TamanioMemoria != -1 {
		// Lógica de asignación de espacio
		err = schemes.AsignarParticion(procesoRequest.Pid, procesoRequest.TamanioMemoria)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("El proceso %d no pudo ser inicializado por falta de particion\n", procesoRequest.Pid)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Printf("## Proceso creado -  PID: %d - Tamaño: %d", procesoRequest.Pid, procesoRequest.TamanioMemoria)
	} else {
		log.Printf("El proceso %d no pudo ser inicializado porque la cola new no está vacía\n", procesoRequest.Pid)
	}
}

func HandleFinalizarProceso(w http.ResponseWriter, r *http.Request) {
	var req requests.RequestFinalizarProceso

	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Liberamos la partición
	err = functions.LiberarProceso(req.Pid)

	if err != nil {
		http.Error(w, "Error al liberar la partición", http.StatusInternalServerError)
		log.Printf("El proceso %d no pudo ser destruido\n", req.Pid)
		return
	}

	tamanioProceso := functions.ObtenerTamanioMemoria(req.Pid)

	// Eliminar las estructuras correspondientes del proceso en la Memoria del Sistema
	delete(globals.MemoriaSistema.TablaProcesos, req.Pid)
	delete(globals.MemoriaSistema.TablaHilos, req.Pid)
	delete(globals.MemoriaSistema.Pseudocodigos, req.Pid)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("OK"))
	if err != nil {
		return
	}
	log.Printf("## Proceso destruido -  PID: %d - Tamaño: %d", req.Pid, tamanioProceso)
}

func HandleCrearHilo(w http.ResponseWriter, r *http.Request) {
	var req requests.RequestCrearHilo

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	err = functions.CrearHiloMemoria(req.Pid, req.Tid, req.Pseudocodigo)

	if err != nil {
		http.Error(w, "Error al crear el hilo", http.StatusInternalServerError)
		log.Printf("El hilo %d no pudo ser creado\n", req.Tid)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("## Hilo creado - (PID:TID) - (%d:%d)", req.Pid, req.Tid)
}

func HandleFinalizarHilo(w http.ResponseWriter, r *http.Request) {
	var req requests.RequestFinalizarHilo

	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	_, existe := globals.MemoriaSistema.TablaHilos[req.Pid][req.Tid]
	if !existe {
		http.Error(w, "Hilo no encontrado", http.StatusNotFound)
		log.Printf("El hilo %d no pudo ser destruido\n", req.Tid)
		return
	}

	// Eliminar las estructuras correspondientes del hilo en la Memoria del Sistema
	delete(globals.MemoriaSistema.TablaHilos[req.Pid], req.Tid)
	delete(globals.MemoriaSistema.Pseudocodigos[req.Pid], req.Tid)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("OK"))
	if err != nil {
		return
	}
	log.Printf("## Hilo destruido - (PID:TID) - (%d:%d)", req.Pid, req.Tid)
}

func HandleMemoryDump(w http.ResponseWriter, r *http.Request) {
	var req requests.RequestDumpMemory

	err := commons.DecodificarJSON(r.Body, &req)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}
	base, limite := functions.ObtenerBaseLimite(req.Pid)
	// Obtener el contenido de la memoria del proceso
	TamanioMemoriaProceso := functions.ObtenerTamanioMemoria(req.Pid)
	ContenidoProceso := functions.ObtenerContenidoMemoria(base, limite)

	// Solicitar al FileSystem la creación del archivo y escribir el contenido

	solicitud := requests.DumpMemoryFS{
		Pid:       uint32(req.Pid),
		Tid:       uint32(req.Tid),
		Tamanio:   TamanioMemoriaProceso,
		Contenido: ContenidoProceso,
	}

	log.Printf("Solicitud antes de codificar: %+v", solicitud)

	solicitudCodificada, err := commons.CodificarJSON(solicitud)

	if err != nil {
		http.Error(w, "Error al codificar JSON", http.StatusBadRequest)
	}

	response, mensaje := cliente.Post2(globals.MConfig.IpFileSystem, globals.MConfig.PortFileSystem, "memory_dump", solicitudCodificada)

	defer response.Body.Close()

	log.Printf("## Memory Dump solicitado - (PID:TID) - (%d:%d)", req.Pid, req.Tid)

	if response != nil && response.StatusCode == http.StatusOK {
		w.WriteHeader(http.StatusOK)

		_, err = w.Write([]byte("OK"))

		if err != nil {
			return
		}

	} else {
		w.WriteHeader(http.StatusInternalServerError)

		_, err = w.Write(mensaje)

		if err != nil {
			return
		}

		log.Printf("Error al realizar el dump de memoria - (PID:TID) - (%d:%d)", req.Pid, req.Tid)
	}
}
