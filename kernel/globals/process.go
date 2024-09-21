package globals

import (
	"log"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/utils/cliente"
	"github.com/sisoputnfrba/tp-golang/utils/commons"
)

func IniciarProceso(w http.ResponseWriter, r *http.Request) {

	var requestCrearProceso RequestProceso

	err := commons.DecodificarJSON(r.Body, &requestCrearProceso)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	if len(commons.ColaNew.Procesos) == 0 {
		log.Println("Cola NEW está vacía, solicitando memoria.")

		// Solicitar espacio en memoria
		respuestaMemoria, err := SolicitarProcesoMemoria(requestCrearProceso.Pseudocodigo, requestCrearProceso.TamanioMemoria)
		if err != nil {
			http.Error(w, "Error al solicitar proceso en memoria", http.StatusInternalServerError)
			return
		}

		// Si la memoria aceptó el proceso, crearlo y pasarlo a READY
		if respuestaMemoria.StatusCode == http.StatusOK {
			pcb := crearProceso(requestCrearProceso.Pseudocodigo, requestCrearProceso.TamanioMemoria)

			commons.ColaReady.Mutex.Lock()
			AgregarProcesoACola(pcb, commons.ColaReady)
			commons.ColaReady.Mutex.Unlock()

			go IniciarHilo(pcb, 0, 0) // Crear el hilo TID 0
			pcb.ContadorHilos++
			// TOCHECK: ¿Debería ser la prioridad del hilo main 0 o la que se recibe por parametro?
			// TODO: Agregar hilo a la cola de hilos del proceso

			log.Println("Proceso creado y movido a READY")

		} else {
			log.Println("Memoria no tiene suficiente espacio. Proceso en espera.")
		}
	} else {
		// Si hay otros procesos, solo encolar el nuevo
		crearProceso(requestCrearProceso.Pseudocodigo, requestCrearProceso.TamanioMemoria)
		log.Println(commons.ColaNew.Procesos[0])
		log.Println("Proceso encolado. Cola NEW no está vacía.")
	}
}

func EliminarProceso(w http.ResponseWriter, r *http.Request) {

	// TODO: Implementar y una vez que termine un proceso, agarrar el primero de la cola NEW (si lo hay) y pasarlo a READY

}

func SolicitarProcesoMemoria(path string, tamanio int) (*http.Response, error) {
	request := RequestProceso{
		Pseudocodigo:   path,
		TamanioMemoria: tamanio,
	}

	solicitudCodificada, err := commons.CodificarJSON(request)

	if err != nil {
		return nil, err
	}

	return cliente.Post(KConfig.IpMemory, KConfig.PortMemory, "process", solicitudCodificada), nil
}

func crearProceso(instrucciones string, tamanio int) commons.PCB {
	pcb := commons.PCB{
		Pid:     commons.PidCounter,
		Estado:  "NEW",
		Tid:     []commons.TCB{},
		Tamanio: tamanio,
		// TODO: VER TEMA ARCHIVO DE INSTRUCCIONES
	}

	commons.MutexPidCounter.Lock()
	commons.PidCounter++
	commons.MutexPidCounter.Unlock()

	commons.ColaNew.Mutex.Lock()
	AgregarProcesoACola(pcb, commons.ColaNew)
	commons.ColaNew.Mutex.Unlock()

	return pcb
}

func IniciarHilo(pcb commons.PCB, prioridad int, tid int) commons.TCB {
	tcb := commons.TCB{
		Pid:       pcb.Pid,
		Tid:       tid,
		Prioridad: prioridad,
	}

	pcb.Tid = append(pcb.Tid, tcb) // Chequear después

	commons.ColaNew.Mutex.Lock()
	commons.ColaReady.Hilos = append(commons.ColaReady.Hilos, tcb)
	commons.ColaNew.Mutex.Unlock()

	return tcb
}

func AgregarProcesoACola(pcb commons.PCB, cola *commons.Colas) {
	cola.Procesos = append(cola.Procesos, pcb)
}

func AgregarHiloACola(tcb commons.TCB, cola *commons.Colas) {
	cola.Hilos = append(cola.Hilos, tcb)
}

func SacarProceso(pcb commons.PCB, cola *commons.Colas) {
	for i, p := range cola.Procesos {
		if p.Pid == pcb.Pid {
			cola.Procesos = append(cola.Procesos[:i], cola.Procesos[i+1:]...)
			break
		}
	}
}
