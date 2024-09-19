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

	respuestaMemoria, err := SolicitarProcesoMemoria(requestCrearProceso.Pseudocodigo, requestCrearProceso.TamanioMemoria)

	log.Println(respuestaMemoria.StatusCode)

	if err != nil {
		http.Error(w, "Error al solicitar el proceso a memoria", http.StatusInternalServerError)
		return
	}

	if respuestaMemoria.StatusCode == http.StatusOK {
		PROCESS_CREATE(requestCrearProceso.Pseudocodigo, requestCrearProceso.TamanioMemoria)
	}

	//TODO: verificar que la cola new este vacia y hacer el pedido a memoria para inicializar el tidMain
	tcbMain := IniciarHilo(pcb.Pid, 0, 0)

	pcb.Tid = append(pcb.Tid, tcbMain.Tid)

	//agregar a cola new el pcb si no se inicializa el proceso a ready

	// TODO: Verificar que hay espacio en memoria para poner el proceso en READY
	// Repetir cuando un proceso finaliza si hay procesos esperando en NEW
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

func PROCESS_CREATE(instrucciones string, tamanio int) {
	pcb := commons.PCB{
		Pid:    commons.PidCounter,
		Estado: "NEW",
		Tid:    []int{},
	}

	commons.MutexPidCounter.Lock()
	commons.PidCounter++
	commons.MutexPidCounter.Unlock()

	commons.ColaNew.Mutex.Lock()
	AgregarACola(pcb, tcb, commons.ColaNew)
	commons.ColaNew.Mutex.Unlock()
}

func IniciarHilo(pid int, prioridad int, tid int) commons.TCB {
	return commons.TCB{
		Pid:       pid,
		Tid:       tid,
		Prioridad: prioridad,
	}
}

func AgregarACola(pcb commons.PCB, cola *commons.Colas) {
	cola.Procesos = append(cola.Procesos, pcb)
}

func SacarProceso(pcb commons.PCB, cola *commons.Colas) {
	for i, p := range cola.Procesos {
		if p.Pid == pcb.Pid {
			cola.Procesos = append(cola.Procesos[:i], cola.Procesos[i+1:]...)
			break
		}
	}
}
