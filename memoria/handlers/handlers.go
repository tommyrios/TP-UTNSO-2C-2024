package handlers

import (
	"fmt"
	request2 "github.com/sisoputnfrba/tp-golang/kernel/handlers/request"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"net/http"
)

func HandleDevolverContexto(w http.ResponseWriter, r *http.Request) {
	var request commons.ContextoDeEjecucion

	err := commons.DecodificarJSON(r.Body, &request)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	var registros = ObtenerRegistros(request.Pid, request.Tid) // ¡¡¡Falta implementar esta función!!!!

	request.Registros = registros

	responseCodificada, err := commons.CodificarJSON(request)

	if err != nil {
		http.Error(w, "Error al codificar el JSON", http.StatusBadRequest)
		return
	}

	cliente.Post(globals.MConfig.IpCpu, globals.MConfig.PortCpu, "/contexto_de_ejecucion", responseCodificada)
}

func HandleActualizarContexto(w http.ResponseWriter, r *http.Request) {
	var request commons.ContextoDeEjecucion

	err := commons.DecodificarJSON(r.Body, &request)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	err = ActualizarRegistros(request.Pid, request.Tid, *request.Registros) // ¡¡¡Falta implementar esta función!!!!

	if err != nil {
		http.Error(w, "Error actualizando los registros", http.StatusInternalServerError)
		return
	}

	// Responder con éxito si se actualizaron correctamente
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Contexto actualizado correctamente"))
}

func ObtenerRegistros(pid int, tid int) *commons.Registros {
	// Buscar los registros guardados en memoria
	// TODO
	registros := &commons.Registros{}

	return registros
}

func ActualizarRegistros(pid int, tid int, registrosActualizados commons.Registros) error {
	// Obtener el puntero a los registros actuales desde la memoria
	var registrosAActualizar = ObtenerRegistros(pid, tid) // ¡¡¡Falta implementar esta función!!!!

	if registrosAActualizar == nil {
		return fmt.Errorf("Registros no encontrados para PID %d y TID %d", pid, tid)
	}

	registrosAActualizar.PC = registrosActualizados.PC
	registrosAActualizar.AX = registrosActualizados.AX
	registrosAActualizar.BX = registrosActualizados.BX
	registrosAActualizar.CX = registrosActualizados.CX
	registrosAActualizar.DX = registrosActualizados.DX
	registrosAActualizar.EX = registrosActualizados.EX
	registrosAActualizar.FX = registrosActualizados.FX
	registrosAActualizar.GX = registrosActualizados.GX
	registrosAActualizar.HX = registrosActualizados.HX

	return nil
}

func HandleSolicitarProceso(w http.ResponseWriter, r *http.Request) {
	var req request2.RequestProcessCreate

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Toda la lógica de verificar si hay espacio, etc.

	// Para este check solo devolvemos OK
	w.WriteHeader(http.StatusOK)
}
