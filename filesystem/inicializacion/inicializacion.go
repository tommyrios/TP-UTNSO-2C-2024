package inicializacion

import (
	"github.com/sisoputnfrba/tp-golang/filesystem/functions"
	"github.com/sisoputnfrba/tp-golang/filesystem/globals"
	"log/slog"
)

func IniciarFS(ruta string) error {
	err := functions.CrearDirectorio(ruta)
	if err != nil {
		slog.Error("Error creando MOUNT_DIR", "error", err)
	}

	// BITMAP
	if !functions.ExisteArchivo(ruta + "/bitmap.dat") {
		slog.Debug("Creando bitmap.dat")
		functions.CrearArchivoBitmap(ruta+"/bitmap.dat", globals.FSConfig.BlockCount)
	} else {
		slog.Debug("bitmap.dat ya creado")
		slog.Debug("cargando bitmap.dat")
		functions.CargarBitmap(ruta+"/bitmap.dat", globals.FSConfig.BlockCount)
	}

	// BLOQUES
	if !functions.ExisteArchivo(ruta + "/bloques.dat") {
		slog.Debug("creando bloques.dat")
		functions.CrearArchivoBloques(ruta+"/bloques.dat", globals.FSConfig.BlockCount*globals.FSConfig.BlockSize)
	} else {
		slog.Debug("bloques.dat ya creado")
	}

	return nil
}
