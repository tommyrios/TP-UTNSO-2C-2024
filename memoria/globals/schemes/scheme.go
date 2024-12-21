package schemes

import (
	"errors"
	"fmt"
	_ "github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/globals/functions"
	"log"
	"log/slog"
	"math"
)

func AsignarParticion(pid int, tamanioProceso int) error {
	indice, err := buscarParticion(tamanioProceso)

	if err != nil {
		return err
	}

	particion := globals.MemoriaUsuario.Particiones[indice]
	particion.Libre = false
	particion.Pid = pid

	if globals.MConfig.Scheme == "FIJAS" {
		globals.MemoriaSistema.TablaProcesos[pid] = &globals.ContextoProceso{Base: particion.Base, Limite: particion.Limite}
		//printParticiones() // DEBUG
	}

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
			globals.MemoriaUsuario.Particiones = append(globals.MemoriaUsuario.Particiones[:indice+1], append([]*globals.Particion{&nuevaParticion}, globals.MemoriaUsuario.Particiones[indice+1:]...)...)
		}
		globals.MemoriaSistema.TablaProcesos[pid] = &globals.ContextoProceso{Base: particion.Base, Limite: particion.Limite}
		// printParticiones // DEBUG
	}

	return nil
}

func buscarIndiceParticion(tamanioProceso int) (int, error) {
	estrategia := globals.MConfig.SearchAlgorithm
	mejorIndice := -1
	mejorValor := math.MaxInt32
	peorValor := -1

	for i, particion := range globals.MemoriaUsuario.Particiones {
		if particion.Libre && particion.Limite-particion.Base+1 >= tamanioProceso {
			espacioLibre := particion.Limite - particion.Base + 1
			switch estrategia {
			case "FIRST":
				return i, nil
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

	return mejorIndice, nil
}

func buscarParticion(tamanioProceso int) (int, error) {
	mejorIndice, _ := buscarIndiceParticion(tamanioProceso)

	if mejorIndice != -1 {
		return mejorIndice, nil
	}

	if globals.MConfig.Scheme == "DINAMICAS" {
		if functions.EspacioLibreTotal() >= tamanioProceso {
			if functions.SolicitarCompactacion() {
				compactarMemoria()

				functions.NotificarFinalizacionCompactacion()

				particionIndice, _ := buscarIndiceParticion(tamanioProceso)

				return particionIndice, nil
			}
		}

		return -1, errors.New("no hay espacio suficiente en memoria")
	}

	return -1, errors.New("no se encontró una partición adecuada")
}

func compactarMemoria() {
	nuevaPosicion := 0
	var nuevasParticiones []*globals.Particion

	for _, particion := range globals.MemoriaUsuario.Particiones {
		if !particion.Libre {
			tamanio := particion.Limite - particion.Base + 1

			copy(globals.MemoriaUsuario.Datos[nuevaPosicion:], globals.MemoriaUsuario.Datos[particion.Base:particion.Limite+1])

			nuevaParticion := globals.Particion{
				Base:   nuevaPosicion,
				Limite: nuevaPosicion + tamanio - 1,
				Libre:  false,
				Pid:    particion.Pid,
			}

			globals.MemoriaSistema.TablaProcesos[particion.Pid] = &globals.ContextoProceso{Base: nuevaPosicion, Limite: nuevaPosicion + tamanio - 1}

			nuevasParticiones = append(nuevasParticiones, &nuevaParticion)

			nuevaPosicion += tamanio
		}
	}

	if nuevaPosicion < len(globals.MemoriaUsuario.Datos) {
		nuevaParticionLibre := globals.Particion{
			Base:   nuevaPosicion,
			Limite: len(globals.MemoriaUsuario.Datos) - 1,
			Libre:  true,
			Pid:    -1,
		}

		nuevasParticiones = append(nuevasParticiones, &nuevaParticionLibre)
	}

	globals.MemoriaUsuario.Particiones = nuevasParticiones

	slog.Debug(fmt.Sprintf("Compactacion finalizada."))
}

func printParticiones() {
	for _, particion := range globals.MemoriaUsuario.Particiones {
		log.Printf("Particion - Base: %d, Limite: %d, Libre: %t, Pid: %d\n", particion.Base, particion.Limite, particion.Libre, particion.Pid)
	}
}
