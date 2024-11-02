package mmu

import (
	"fmt"
	"log"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

func TraducirDireccion(base, desplazamiento, limite uint32) (uint32, error) {
	if desplazamiento >= limite {
		return 0, fmt.Errorf("SEGMENTATION FAULT")
	}
	return desplazamiento + base, nil
}

func LeerMemoria(direccionLogica uint32) (uint32, error) {
	direccionFisica, err := TraducirDireccion(globals.Registros.Base, direccionLogica, globals.Registros.Limite)
	if err != nil {
		return 0, err
	}

	log.Printf("## TID: %d - Acción: LEER - Dirección Física: %d", *globals.Tid, direccionFisica)
	return 0, nil
}

func EscribirMemoria(direccionLogica uint32, valor uint32) error {
	direccionFisica, err := TraducirDireccion(globals.Registros.Base, direccionLogica, globals.Registros.Limite)
	if err != nil {
		return err
	}
	log.Printf("## TID: %d - Acción: ESCRIBIR - Dirección Física: %d", *globals.Tid, direccionFisica)
	return nil
}
