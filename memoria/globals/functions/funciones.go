package functions

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/handlers/requests"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
	"log"
	"net/http"
	"os"
)

// FUNCIONES CPU

func ObtenerRegistros(pid int, tid int) *commons.Registros {

	registros := globals.MemoriaSistema.TablaHilos[pid][tid]

	return registros
}

func ObtenerBaseLimite(pid int) (int, int) {

	base := globals.MemoriaSistema.TablaProcesos[pid].Base
	limite := globals.MemoriaSistema.TablaProcesos[pid].Limite

	return base, limite
}

func ActualizarRegistros(pid int, tid int, registrosActualizados *commons.Registros) error {
	if _, existeMap := globals.MemoriaSistema.TablaHilos[pid]; !existeMap {
		return fmt.Errorf("el PID %d no existe en la tabla de hilos", pid)
	}

	registrosAActualizar, existeMap := globals.MemoriaSistema.TablaHilos[pid][tid]

	if !existeMap || registrosAActualizar == nil {
		return fmt.Errorf("el TID %d no existe para el PID %d", tid, pid)
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

func ObtenerInstruccion(pid int, tid int, pc int) (requests.ResponseObtenerInstruccion, error) {
	instruccion := globals.MemoriaSistema.Pseudocodigos[pid][tid].Instrucciones[pc]

	responseInstruccion := requests.ResponseObtenerInstruccion{Instruccion: instruccion}

	if instruccion == "" {
		return responseInstruccion, errors.New("instrucción no encontrada")
	}

	return responseInstruccion, nil
}

func ObtenerTamanioMemoria(pid int) int {
	return globals.MemoriaSistema.TablaProcesos[pid].Limite - globals.MemoriaSistema.TablaProcesos[pid].Base
}

func LeerMemoria(direccion int, pid int) ([]byte, error) {

	if direccion < 0 || direccion+4 >= len(globals.MemoriaUsuario.Datos) {
		return nil, fmt.Errorf("dirección de memoria inválida")
	}
	for _, particion := range globals.MemoriaUsuario.Particiones {
		if direccion >= particion.Base && direccion+4 <= particion.Limite && particion.Pid == pid {
			return globals.MemoriaUsuario.Datos[direccion : direccion+4], nil
		}
	}

	return nil, fmt.Errorf("segmentation fault")
}

func EscribirMemoria(direccion int, pid int, datos []byte) error {

	if len(datos) == 0 {
		return fmt.Errorf("está vacío el contenido a escribir")
	}

	for _, particion := range globals.MemoriaUsuario.Particiones {
		if direccion >= particion.Base && direccion+4 <= particion.Limite && particion.Pid == pid {
			copy(globals.MemoriaUsuario.Datos[direccion:direccion+4], datos)
			break
		} else if particion.Pid == pid && direccion > particion.Limite {
			return fmt.Errorf("segmentation fault")
		}
	}
	value := binary.LittleEndian.Uint32(globals.MemoriaUsuario.Datos[direccion : direccion+4])

	log.Println(value)

	return nil
}

// FUNCIONES KERNEL

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
			particiones[indice-1].Limite = particiones[indice].Limite
			particiones = append(particiones[:indice], particiones[indice+1:]...)
			indice-- // Actualiza índice después de la consolidación
		}

		// Consolidar con partición siguiente si está libre
		if indice < len(particiones)-1 && particiones[indice+1].Libre {
			particiones[indice].Limite = particiones[indice+1].Limite
			particiones = append(particiones[:indice+1], particiones[indice+2:]...)
		}
	}

	return nil
}

func CrearHiloMemoria(pid int, tid int, pseudocodigo string) error {
	instrucciones, err := DesglosarPseudocodigo(pseudocodigo)

	if err != nil {
		log.Printf("Error al desglosar el pseudocódigo: %s\n", pseudocodigo)
		return err
	}

	if globals.MemoriaSistema.TablaHilos[pid] == nil {
		globals.MemoriaSistema.TablaHilos[pid] = make(map[int]*commons.Registros)
	}

	if globals.MemoriaSistema.Pseudocodigos[pid] == nil {
		globals.MemoriaSistema.Pseudocodigos[pid] = make(map[int]*globals.InstruccionesHilo)
	}

	globals.MemoriaSistema.Pseudocodigos[pid][tid] = &globals.InstruccionesHilo{Instrucciones: instrucciones}
	globals.MemoriaSistema.TablaHilos[pid][tid] = &commons.Registros{AX: 0, BX: 0, CX: 0, DX: 0, EX: 0, FX: 0, GX: 0, HX: 0, PC: 0}

	return nil
}

func DesglosarPseudocodigo(pseudocodigo string) ([]string, error) {
	archivo, err := os.Open("../the-last-of-c-pruebas/" + pseudocodigo)
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

func ObtenerContenidoMemoria(base, limite int) []byte {
	if base < 0 || limite >= len(globals.MemoriaUsuario.Datos) || base > limite {
		return nil
	}

	tamanio := limite - base + 1
	contenido := make([]byte, tamanio)
	copy(contenido, globals.MemoriaUsuario.Datos[base:limite+1])
	return contenido
}

// FUNCIONES MEMORIA (PARA PARTICIONAMIENTO)

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
	response, err := http.Post(fmt.Sprintf("http://%s:%d/compactacion", globals.MConfig.IpKernel, globals.MConfig.PortKernel), "application/json", nil)
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
