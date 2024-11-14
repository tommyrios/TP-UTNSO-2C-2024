package functions

import (
	"errors"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"math"
	"net/http"
)

type MemUsuario globals.MemUsuario

func ObtenerRegistros(pid int, tid int) commons.Registros {

	registros := globals.MemoriaSistema.TablaHilos[pid][tid]

	return commons.Registros(registros)
}

func ObtenerBaseLimite(pid int, tid int) (int, int) {

	base := globals.MemoriaSistema.TablaProcesos[pid].Base
	limite := globals.MemoriaSistema.TablaProcesos[pid].Limite

	return base, limite
}

func ObtenerTamanioMemoria(base int, limite int) int {
	return limite - base
}

func ActualizarRegistros(pid int, tid int, registrosActualizados commons.Registros) error {

	registrosAActualizar := globals.MemoriaSistema.TablaHilos[pid][tid]

	/* Chequear error de nil

	if registrosAActualizar.PC == nil {
		return fmt.Errorf("Registros no encontrados para PID %d y TID %d", pid, tid)
	}

	*/

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

func ObtenerInstruccion(pid int, tid int, pc uint32) (string, error) {

	instruccion := globals.MemoriaSistema.Pseudocodigos[pid][tid].Instrucciones[pc]

	if instruccion == "" {
		return "", fmt.Errorf("Instrucción no encontrada para PID %d, TID %d y PC %d", pid, tid, pc)
	}

	return instruccion, nil
}

func LiberarProceso(pid int) error {
	indice := -1
	particiones := globals.MemoriaUsuario.Particiones
	// Buscar la partición correspondiente al PID
	for i, particion := range particiones {
		if particion.Pid == pid && !particion.Libre {
			indice = i
			break
		}
	}

	if indice == -1 {
		return errors.New("Proceso no encontrado o ya está liberado")
	}

	// Liberar la partición
	particiones[indice].Pid = -1
	particiones[indice].Libre = true

	if globals.MConfig.Scheme == "DINAMICAS" {
		// Consolidar con partición anterior si está libre
		if indice > 0 && particiones[indice-1].Libre {
			particiones[indice-1].Limite += particiones[indice].Limite
			particiones = append(particiones[:indice], particiones[indice+1:]...)
			indice-- // Actualiza índice después de la consolidación
		}

		// Consolidar con partición siguiente si está libre
		if indice < len(particiones)-1 && particiones[indice+1].Libre {
			particiones[indice].Limite += particiones[indice+1].Limite
			particiones = append(particiones[:indice+1], particiones[indice+2:]...)
		}
	}

	return nil
}

func LeerMemoria(byteDireccion byte) ([]byte, error) {

	direccion := int(byteDireccion)

	if direccion < 0 || direccion+4 >= len(globals.MemoriaUsuario.Datos) {
		return nil, fmt.Errorf("Dirección de memoria inválida")
	}

	//verificar segmentation fault
	return globals.MemoriaUsuario.Datos[direccion : direccion+4], nil
}

func EscribirMemoria(byteDireccion byte, pid int, datos []byte) error {

	if len(datos) != 4 {
		return fmt.Errorf("Se deben proporcionar exactamente 4 bytes")
	}

	proceso := globals.MemoriaSistema.TablaProcesos[pid]

	direccionFisica := proceso.Base + int(byteDireccion)

	if direccionFisica < 0 || direccionFisica+4 >= proceso.Limite {
		return fmt.Errorf("Segmentation fault")
	}

	copy(globals.MemoriaUsuario.Datos[direccionFisica:direccionFisica+4], datos)

	return nil
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

func EsEspacioLibre(inicio, tamano int) bool {
	for i := inicio; i < inicio+tamano; i++ {
		if globals.MemoriaUsuario.Datos[i] != 0 { // 0 indica espacio libre
			return false
		}
	}
	return true
}

func AsignarEspacio(pid, inicio, tamano int) {
	for i := inicio; i < inicio+tamano; i++ {
		globals.MemoriaUsuario.Datos[i] = 1 // 1 indica espacio ocupado
	}

	globals.MemoriaSistema.TablaProcesos[pid] = globals.ContextoProceso{
		Base:   inicio,
		Limite: inicio + tamano - 1,
	}
}

func CalcularDesperdicio(inicio, tamano int) int {
	espacioLibre := 0
	for i := inicio + tamano; i < len(globals.MemoriaUsuario.Datos) && globals.MemoriaUsuario.Datos[i] == 0; i++ {
		espacioLibre++
	}
	return espacioLibre
}

func EspacioLibreTotal() int {
	espacioLibre := 0
	for _, byte := range globals.MemoriaUsuario.Datos {
		if byte == 0 { // 0 indica espacio libre
			espacioLibre++
		}
	}
	return espacioLibre
}

func SolicitarCompactacion() bool {
	// Enviar solicitud HTTP al Kernel para compactación
	response, err := http.Post(fmt.Sprintf("http://%s:%d/compactar", globals.MConfig.IpKernel, globals.MConfig.PortKernel), "application/json", nil)
	if err != nil || response.StatusCode != http.StatusOK {
		return false // Falló la solicitud o el Kernel no aprobó la compactación
	}
	return true
}

func NotificarFinalizacionCompactacion() {
	// Notificar al Kernel que la compactación ha finalizado
	http.Post(fmt.Sprintf("http://%s:%d/compactacion_finalizada", globals.MConfig.IpKernel, globals.MConfig.PortKernel), "application/json", nil)
}

func ObtenerContenidoMemoria(base, limite int) []byte {
	if base < 0 || limite >= len(globals.MemoriaUsuario.Datos) || base > limite {
		return nil
	}

	tamanio := limite - base + 1
	contenido := make([]byte, tamanio)
	copy(contenido, globals.MemoriaUsuario.Datos[base:limite+1])

	return contenido
}
