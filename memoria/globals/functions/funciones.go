package functions

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
	"net/http"
	"os"
)

type MemUsuario globals.MemUsuario

func ObtenerRegistros(pid int, tid int) globals.ContextoHilo {

	return *globals.MemoriaSistema.TablaHilos[pid][tid]
}

func ObtenerBaseLimite(pid int) (int, int) {

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
		return "", fmt.Errorf("instrucción no encontrada para PID %d, TID %d y PC %d", pid, tid, pc)
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
		return errors.New("proceso no encontrado o ya está liberado")
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
		return nil, fmt.Errorf("dirección de memoria inválida")
	}

	//verificar segmentation fault
	return globals.MemoriaUsuario.Datos[direccion : direccion+4], nil
}

func EscribirMemoria(byteDireccion byte, pid int, datos []byte) error {

	if len(datos) != 4 {
		return fmt.Errorf("se deben proporcionar exactamente 4 bytes")
	}

	proceso := globals.MemoriaSistema.TablaProcesos[pid]

	direccionFisica := proceso.Base + int(byteDireccion)

	if direccionFisica < 0 || direccionFisica+4 >= proceso.Limite {
		return fmt.Errorf("segmentation fault")
	}

	copy(globals.MemoriaUsuario.Datos[direccionFisica:direccionFisica+4], datos)

	return nil
}

func EspacioLibreTotal() int {
	espacioLibre := 0
	particiones := globals.MemoriaUsuario.Particiones

	for _, particion := range particiones {
		if particion.Libre { // 0 indica espacio libre
			espacioLibre += 1 + (particion.Limite - particion.Base)
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
	log.Println("Compactacion solicitada.")
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

func CrearHilo(pid int, tid int, pseudocodigo string) {
	// Crear hilo con pseudocódigo y agregarlo a la tabla de hilos
	instrucciones, err := DesglosarPseudocodigo(pseudocodigo)

	if err != nil {
		log.Printf("Error al desglosar el pseudocódigo: %s\n")
	}

	globals.MemoriaSistema.Pseudocodigos[pid][tid] = &globals.InstruccionesHilo{Instrucciones: instrucciones}
	globals.MemoriaSistema.TablaHilos[pid][tid] = &globals.ContextoHilo{AX: 0, BX: 0, CX: 0, DX: 0, EX: 0, FX: 0, GX: 0, HX: 0, PC: 0}
}

func DesglosarPseudocodigo(pseudocodigo string) ([]string, error) {
	archivo, err := os.Open("/kernel/" + pseudocodigo)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo: %w", err)
	}
	defer archivo.Close()

	var lineas []string

	// Crear un scanner para leer línea por línea
	scanner := bufio.NewScanner(archivo)
	for scanner.Scan() {
		lineas = append(lineas, scanner.Text())
	}

	// Comprobar si hubo errores durante la lectura
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error al leer el archivo: %w", err)
	}

	return lineas, nil
}
