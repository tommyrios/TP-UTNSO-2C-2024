package globals

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func ObtenerRegistros(pid int, tid int) commons.Registros {

	registros := MemoriaSistema.TablaHilos[pid][tid]

	return commons.Registros(registros)
}

func ObtenerBaseLimite(pid int, tid int) (int, int) {

	base := MemoriaSistema.TablaProcesos[pid].Base
	limite := MemoriaSistema.TablaProcesos[pid].Limite

	return base, limite
}

func ActualizarRegistros(pid int, tid int, registrosActualizados commons.Registros) error {

	registrosAActualizar := MemoriaSistema.TablaHilos[pid][tid]

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

	instruccion := MemoriaSistema.Pseudocodigos[pid][tid].Instrucciones[pc]

	if instruccion == "" {
		return "", fmt.Errorf("Instrucción no encontrada para PID %d, TID %d y PC %d", pid, tid, pc)
	}

	return instruccion, nil
}

func LeerMemoria(byteDireccion byte) ([]byte, error) {

	direccion := int(byteDireccion)

	if direccion < 0 || direccion+4 >= len(MemoriaUsuario) {
		return nil, fmt.Errorf("Dirección de memoria inválida")
	}

	return MemoriaUsuario[direccion : direccion+4], nil
}

func EscribirMemoria(byteDireccion byte, pid int, datos []byte) error {

	if len(datos) != 4 {
		return fmt.Errorf("Se deben proporcionar exactamente 4 bytes")
	}

	proceso := MemoriaSistema.TablaProcesos[pid]

	direccionFisica := proceso.Base + int(byteDireccion)

	if direccionFisica < 0 || direccionFisica+4 >= proceso.Limite {
		return fmt.Errorf("Segmentation fault")
	}

	copy(MemoriaUsuario[direccionFisica:direccionFisica+4], datos)

	return nil
}
