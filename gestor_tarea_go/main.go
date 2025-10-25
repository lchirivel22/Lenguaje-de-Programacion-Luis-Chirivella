package main

import (
	"fmt"
	"os"
)

func main() {
	var tareas Tareas
	almacenamiento := NuevoAlmacenamiento[Tareas]("tareas.json")

	err := almacenamiento.Cargar(&tareas)
	if err != nil {
		fmt.Printf("Error al cargar tareas: %v\n", err)
		os.Exit(1)
	}

	comandos := NuevoBanderasComando()
	err = comandos.Ejecutar(&tareas)
	if err != nil {
		fmt.Printf("Error al ejecutar comando: %v\n", err)
		os.Exit(1)
	}

	err = almacenamiento.Guardar(&tareas)
	if err != nil {
		fmt.Printf("Error al guardar tareas: %v\n", err)
		os.Exit(1)
	}
}