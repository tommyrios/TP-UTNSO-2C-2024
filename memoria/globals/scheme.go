package globals

import (
	"errors"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	request3 "github.com/sisoputnfrba/tp-golang/memoria/handlers/request"
	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"math"
	"net/http"
)

func AsignarParticionFija(pid int, tamanioProceso int) bool {
	tamanioParticion := MejorAjuste(tamanioProceso)

	switch MConfig.SearchAlgorithm {
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

	for _, tam := range MConfig.Partitions {
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

func AsignarParticionDinamica(pid int, tamanioProceso int) bool {

	tamanioParticion := tamanioProceso
	var err error

	switch MConfig.SearchAlgorithm {
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
			respuestaKernel, err2 := comunicarKernelCompactacion(pid, tamanioParticion) // Función que comunica al Kernel la necesidad de compactar
			if err2 != nil {
				return false
			}
			else {AsignarParticionDinamica(pid, tamanioParticion)}
		} else {
			// No hay espacio suficiente ni posibilidad de compactar
			return false
		}
	}

	return false
}

func comunicarKernelCompactacion(pid int, tamanioProceso int) (*http.Response, error) {
	request := request3.RequestProcesoMemoria{
		Pid:            pid,
		TamanioMemoria: tamanioProceso,
	}

	solicitudCodificada, err := commons.CodificarJSON(request)

	if err != nil {
		return nil, fmt.Errorf("Error al codificar la solicitud: %w", err)
	}

	response := cliente.Post(globals.MConfig.IpKernel, globals.MConfig.PortKernel, "/compactacion", solicitudCodificada)

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("el Kernel devolvió un error: %s", response.Status)
	}

	return response, nil
}

func asignarParticionFirstFit(pid int, tamanioParticion int) error {
	for i := 0; i <= len(MemoriaUsuario)-tamanioParticion; i++ {
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

	for i := 0; i <= len(MemoriaUsuario)-tamanioParticion; i++ {
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

	for i := 0; i <= len(MemoriaUsuario)-tamanioParticion; i++ {
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
		if MemoriaUsuario[i] != 0 { // 0 indica espacio libre
			return false
		}
	}
	return true
}

func asignarEspacio(pid, inicio, tamano int) {
	for i := inicio; i < inicio+tamano; i++ {
		MemoriaUsuario[i] = 1 // 1 indica espacio ocupado
	}

	MemoriaSistema.TablaProcesos[pid] = ContextoProceso{
		Base:   inicio,
		Limite: inicio + tamano - 1,
	}
}

func calcularDesperdicio(inicio, tamano int) int {
	espacioLibre := 0
	for i := inicio + tamano; i < len(MemoriaUsuario) && MemoriaUsuario[i] == 0; i++ {
		espacioLibre++
	}
	return espacioLibre
}

func espacioLibreTotal() int {
	espacioLibre := 0
	for _, byte := range MemoriaUsuario {
		if byte == 0 { // 0 indica espacio libre
			espacioLibre++
		}
	}
	return espacioLibre
}
