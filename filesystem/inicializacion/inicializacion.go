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
		slog.Info("creando bitmap.dat")
		functions.CrearArchivoBitmap(ruta+"/bitmap.dat", globals.FSConfig.BlockCount)
	} else {
		slog.Info("bitmap.dat ya creado")
		slog.Info("cargando bitmap.dat")
		functions.CargarBitmap(ruta+"/bitmap.dat", globals.FSConfig.BlockCount)
	}

	// BLOQUES
	if !functions.ExisteArchivo(ruta + "/bloques.dat") {
		slog.Info("creando bloques.dat")
		functions.CrearArchivoBloques(ruta+"/bloques.dat", globals.FSConfig.BlockCount*globals.FSConfig.BlockSize)
	} else {
		slog.Info("bloques.dat ya creado")
	}

	return nil
}
