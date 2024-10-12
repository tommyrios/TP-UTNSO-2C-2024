package globals

func InicializarMemoria() {
	switch MConfig.Scheme {
	case "FIJAS":
		particionesFijas()
	case "VARIABLES":
		particionesVariables()
	}
}

func particionesFijas() {
	//TODO
}

func particionesVariables() {
	//TODO
}
