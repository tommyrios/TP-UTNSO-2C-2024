package handlers

import (
	"errors"
	"fmt"
	request2 "github.com/sisoputnfrba/tp-golang/kernel/handlers/request"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	request3 "github.com/sisoputnfrba/tp-golang/memoria/handlers/request"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"math"
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
	var req request2.RequestProcessCreate

	err := commons.DecodificarJSON(r.Body, &req)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	esquemaFijo := globals.MConfig.Scheme == "fijo"

	// Lógica de asignación de espacio
	if esquemaFijo {
		if !asignarParticionFija(req.Pid, req.TamanioMemoria) {
			http.Error(w, "No hay espacio en particiones fijas", http.StatusConflict)
			return
		}
	} else {
		if !asignarParticionDinamica(req.Pid, req.TamanioMemoria) {
			http.Error(w, "No hay espacio en particiones dinámicas, compactación requerida", http.StatusConflict)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func asignarParticionFija(pid int, tamanioProceso int) bool {
	tamanioParticion := MejorAjuste(tamanioProceso)

	switch globals.MConfig.SearchAlgorithm {
	case "FIRST":
		err := asignarParticionFirstFit(pid, tamanioParticion)
		if err == nil {
			return true
		}
	case "BEST":
		err := asignarParticionBestFit(pid, tamanioParticion)
		if err == nil {
			return true
		}

	case "WORST":
		err := asignarParticionWorstFit(pid, tamanioParticion)
		if err == nil {
			return true
		}
	}

	return false
}

func MejorAjuste(x int) int {
	mejorTamanio := -1
	menorDesperdicio := math.MaxInt32

	for _, tam := range globals.MConfig.Partitions {
		if tam >= x {
			desperdicio := tam - x
			if desperdicio < menorDesperdicio {
				menorDesperdicio = desperdicio
				mejorTamanio = tam
			}
		}
	}

	return mejorTamanio
}

func asignarParticionDinamica(pid int, tamanioProceso int) bool {

	tamanioParticion := tamanioProceso
	var err error

	switch globals.MConfig.SearchAlgorithm {
	case "FIRST":
		err := asignarParticionFirstFit(pid, tamanioParticion)
		if err == nil {
			return true
		}

	case "BEST":
		err := asignarParticionBestFit(pid, tamanioParticion)
		if err == nil {
			return true
		}

	case "WORST":
		err := asignarParticionWorstFit(pid, tamanioParticion)
		if err == nil {
			return true
		}
	}

	if err != nil {
		// Si el espacio libre total permite asignar el proceso, necesitamos compactar
		if espacioLibreTotal() >= tamanioParticion {
			// notificarNecesidadDeCompactacion() // Función que comunica al Kernel la necesidad de compactar
			return false
		} else {
			// No hay espacio suficiente ni posibilidad de compactar
			return false
		}
	}

	return false
}

func asignarParticionFirstFit(pid int, tamanioParticion int) error {
	for i := 0; i <= len(globals.MemoriaUsuario)-tamanioParticion; i++ {
		if esEspacioLibre(i, tamanioParticion) {
			asignarEspacio(pid, i, tamanioParticion)
			return nil
		}
	}
	return errors.New("no hay espacio contiguo suficiente en memoria")
}

func asignarParticionBestFit(pid int, tamanioParticion int) error {
	mejorInicio := -1
	menorDesperdicio := math.MaxInt32

	for i := 0; i <= len(globals.MemoriaUsuario)-tamanioParticion; i++ {
		if esEspacioLibre(i, tamanioParticion) {
			desperdicio := calcularDesperdicio(i, tamanioParticion)
			if desperdicio < menorDesperdicio {
				menorDesperdicio = desperdicio
				mejorInicio = i
			}
		}
	}

	if mejorInicio != -1 {
		asignarEspacio(pid, mejorInicio, tamanioParticion)
		return nil
	}
	return errors.New("no hay espacio contiguo suficiente en memoria")
}

func asignarParticionWorstFit(pid int, tamanioParticion int) error {
	peorInicio := -1
	mayorEspacioLibre := -1

	for i := 0; i <= len(globals.MemoriaUsuario)-tamanioParticion; i++ {
		if esEspacioLibre(i, tamanioParticion) {
			espacioLibre := calcularDesperdicio(i, tamanioParticion)
			if espacioLibre > mayorEspacioLibre {
				mayorEspacioLibre = espacioLibre
				peorInicio = i
			}
		}
	}

	if peorInicio != -1 {
		asignarEspacio(pid, peorInicio, tamanioParticion)
		return nil
	}
	return errors.New("no hay espacio contiguo suficiente en memoria")
}

func esEspacioLibre(inicio, tamano int) bool {
	for i := inicio; i < inicio+tamano; i++ {
		if globals.MemoriaUsuario[i] != 0 { // 0 indica espacio libre
			return false
		}
	}
	return true
}

func asignarEspacio(pid, inicio, tamano int) {
	for i := inicio; i < inicio+tamano; i++ {
		globals.MemoriaUsuario[i] = 1 // 1 indica espacio ocupado
	}

	globals.MemoriaSistema.TablaProcesos[pid] = ContextoProceso{
		Base:   inicio,
		Limite: inicio + tamano - 1,
	}
}

func calcularDesperdicio(inicio, tamano int) int {
	espacioLibre := 0
	for i := inicio + tamano; i < len(globals.MemoriaUsuario) && globals.MemoriaUsuario[i] == 0; i++ {
		espacioLibre++
	}
	return espacioLibre
}

func espacioLibreTotal() int {
	espacioLibre := 0
	for _, byte := range globals.MemoriaUsuario {
		if byte == 0 { // 0 indica espacio libre
			espacioLibre++
		}
	}
	return espacioLibre
}
