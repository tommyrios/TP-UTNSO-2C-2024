package utils

import (
	"strconv"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

func ValorRegistros(registro string) uint32 {
	var valor uint32
	switch v := globals.Regis[registro].(type) {
	case *uint32:
		valor = *v
	case *uint8:
		valor = uint32(*v)
	}
	return valor
}

func ConvertirStringAEntero(s string) uint32 {
	valor, _ := strconv.ParseUint(s, 10, 32)
	return uint32(valor)
}
