package schemes

import (
	"errors"
	_ "github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/globals/functions"
	"math"
	"sync"
)

var CompactacionCond = sync.NewCond(&sync.Mutex{})

func AsignarParticionFija(pid int, tamanioProceso int) bool {
	tamanioParticion := functions.MejorAjuste(tamanioProceso)

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

func AsignarParticionDinamica(pid int, tamanioProceso int) bool {

	tamanioParticion := tamanioProceso

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

	// Caso de que haya espacio total pero no contiguo
	if functions.EspacioLibreTotal() >= tamanioParticion {
		if functions.SolicitarCompactacion() {
			// Espera a que el Kernel confirme que se puede compactar
			CompactacionCond.L.Lock()
			CompactacionCond.Wait() // Espera hasta que Kernel confirme que puede compactar
			CompactacionCond.L.Unlock()

			// Realiza la compactación
			compactarMemoria()

			// Notifica al Kernel que la compactación ha finalizado
			functions.NotificarFinalizacionCompactacion()
			return true
		}
	}

	// Caso de que no hay espacio total disponible
	return false
}

func compactarMemoria() {
	nuevaPosicion := 0

	for pid, proceso := range globals.MemoriaSistema.TablaProcesos {
		base := proceso.Base
		limite := proceso.Limite
		tamanio := limite - base + 1

		// Mover el proceso a la nueva posición
		copy(globals.MemoriaUsuario[nuevaPosicion:], globals.MemoriaUsuario[base:limite+1])

		// Actualizar la tabla de procesos con la nueva posición
		globals.MemoriaSistema.TablaProcesos[pid] = globals.ContextoProceso{
			Base:   nuevaPosicion,
			Limite: nuevaPosicion + tamanio - 1,
		}

		// Limpiar la memoria antigua
		for i := base; i <= limite; i++ {
			globals.MemoriaUsuario[i] = 0
		}

		// Actualizar la nueva posición
		nuevaPosicion += tamanio
	}
}

func asignarParticionFirstFit(pid int, tamanioParticion int) error {
	for i := 0; i <= len(globals.MemoriaUsuario)-tamanioParticion; i++ {
		if functions.EsEspacioLibre(i, tamanioParticion) {
			functions.AsignarEspacio(pid, i, tamanioParticion)
			return nil
		}
	}
	return errors.New("no hay espacio contiguo suficiente en memoria")
}

func asignarParticionBestFit(pid int, tamanioParticion int) error {
	mejorInicio := -1
	menorDesperdicio := math.MaxInt32

	for i := 0; i <= len(globals.MemoriaUsuario)-tamanioParticion; i++ {
		if functions.EsEspacioLibre(i, tamanioParticion) {
			desperdicio := functions.CalcularDesperdicio(i, tamanioParticion)
			if desperdicio < menorDesperdicio {
				menorDesperdicio = desperdicio
				mejorInicio = i
			}
		}
	}

	if mejorInicio != -1 {
		functions.AsignarEspacio(pid, mejorInicio, tamanioParticion)
		return nil
	}
	return errors.New("no hay espacio contiguo suficiente en memoria")
}

func asignarParticionWorstFit(pid int, tamanioParticion int) error {
	peorInicio := -1
	mayorEspacioLibre := -1

	for i := 0; i <= len(globals.MemoriaUsuario)-tamanioParticion; i++ {
		if functions.EsEspacioLibre(i, tamanioParticion) {
			espacioLibre := functions.CalcularDesperdicio(i, tamanioParticion)
			if espacioLibre > mayorEspacioLibre {
				mayorEspacioLibre = espacioLibre
				peorInicio = i
			}
		}
	}

	if peorInicio != -1 {
		functions.AsignarEspacio(pid, peorInicio, tamanioParticion)
		return nil
	}
	return errors.New("no hay espacio contiguo suficiente en memoria")
}
