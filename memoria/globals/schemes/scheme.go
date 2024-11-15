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

func AsignarParticion(pid int, tamanioProceso int) error {
	indice, err := buscarParticion(tamanioProceso)
	if err != nil {
		return err
	}

	particion := &globals.MemoriaUsuario.Particiones[indice]
	particion.Libre = false
	particion.Pid = pid

	// Ajusta la partición si sobra espacio (solo en dinámicas)
	if globals.MConfig.Scheme == "DINAMICAS" {
		espacioSobrante := particion.Limite - particion.Base - tamanioProceso
		if espacioSobrante > 0 {
			nuevaParticion := globals.Particion{
				Base:   particion.Base + tamanioProceso,
				Limite: particion.Limite,
				Libre:  true,
				Pid:    -1,
			}
			particion.Limite = particion.Base + tamanioProceso - 1
			globals.MemoriaUsuario.Particiones = append(globals.MemoriaUsuario.Particiones[:indice+1], append([]globals.Particion{nuevaParticion}, globals.MemoriaUsuario.Particiones[indice+1:]...)...)
		}
	}

	return nil
}

func buscarParticion(tamanioProceso int) (int, error) {
	estrategia := globals.MConfig.SearchAlgorithm
	mejorIndice := -1
	mejorValor := math.MaxInt32
	peorValor := -1

	for i, particion := range globals.MemoriaUsuario.Particiones {
		if particion.Libre && particion.Limite-particion.Base >= tamanioProceso {
			espacioLibre := particion.Limite - particion.Base

			switch estrategia {
			case "FIRST":
				return i, nil // Devuelve la primera partición encontrada
			case "BEST":
				if espacioLibre < mejorValor {
					mejorValor = espacioLibre
					mejorIndice = i
				}
			case "WORST":
				if espacioLibre > peorValor {
					peorValor = espacioLibre
					mejorIndice = i
				}
			}
		}
	}

	if mejorIndice != -1 {
		return mejorIndice, nil
	}

	// Caso de que haya espacio total pero no contiguo
	if globals.MConfig.Scheme == "DINAMICAS" {
		if functions.EspacioLibreTotal() >= tamanioProceso {
			if functions.SolicitarCompactacion() {
				// Espera a que el Kernel confirme que se puede compactar
				CompactacionCond.L.Lock()
				CompactacionCond.Wait() // Espera hasta que Kernel confirme que puede compactar
				CompactacionCond.L.Unlock()

				// Realiza la compactación
				compactarMemoria()

				// Notifica al Kernel que la compactación ha finalizado
				functions.NotificarFinalizacionCompactacion()

				particionIndice, _ := buscarParticion(tamanioProceso)
				return particionIndice, nil
			}
		}

		return -1, errors.New("no hay espacio suficiente en memoria")
	}

	return -1, errors.New("no se encontró una partición adecuada")
}

func compactarMemoria() {
	nuevaPosicion := 0
	var nuevasParticiones []globals.Particion

	for _, particion := range globals.MemoriaUsuario.Particiones {
		if !particion.Libre {
			tamanio := particion.Limite - particion.Base + 1

			// Mover datos al nuevo espacio
			copy(globals.MemoriaUsuario.Datos[nuevaPosicion:], globals.MemoriaUsuario.Datos[particion.Base:particion.Limite+1])

			// Crear una partición actualizada
			nuevaParticion := globals.Particion{
				Base:   nuevaPosicion,
				Limite: nuevaPosicion + tamanio - 1,
				Libre:  false,
				Pid:    particion.Pid,
			}
			nuevasParticiones = append(nuevasParticiones, nuevaParticion)
			nuevaPosicion += tamanio
		}
	}

	// Crear una partición libre con el resto de la memoria
	if nuevaPosicion < len(globals.MemoriaUsuario.Datos) {
		nuevaParticionLibre := globals.Particion{
			Base:   nuevaPosicion,
			Limite: len(globals.MemoriaUsuario.Datos) - 1,
			Libre:  true,
			Pid:    -1,
		}
		nuevasParticiones = append(nuevasParticiones, nuevaParticionLibre)
	}

	// Actualizar las particiones y limpiar los datos no usados
	globals.MemoriaUsuario.Particiones = nuevasParticiones
	for i := nuevaPosicion; i < len(globals.MemoriaUsuario.Datos); i++ {
		globals.MemoriaUsuario.Datos[i] = 0
	}
}
