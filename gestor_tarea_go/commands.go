package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type BanderasComando struct {
	Agregar  string
	Eliminar int
	Editar   string
	Alternar int
	Listar   bool
	Ayuda    bool
}

func NuevoBanderasComando() *BanderasComando {
	bc := &BanderasComando{}

	flag.StringVar(&bc.Agregar, "agregar", "", "Añade una tarea nueva. Ejemplo: -agregar \"Titulo de ejemplo\"")
	flag.IntVar(&bc.Eliminar, "eliminar", -1, "Quita una tarea usando su ID. Ejemplo: -eliminar <ID>")
	flag.StringVar(&bc.Editar, "editar", "", "Modifica el titulo de una tarea por ID. Ejemplo: -editar \"ID:Nuevo Titulo\"")
	flag.IntVar(&bc.Alternar, "alternar", -1, "Alterna el estado de una tarea (Por hacer/Hecha) por su ID. Ejemplo: -alternar <ID>")
	flag.BoolVar(&bc.Listar, "listar", false, "Muestra la lista de tareas.")
	flag.BoolVar(&bc.Ayuda, "ayuda", false, "Muestra la ayuda disponible.")

	flag.Parse()

	return bc
}

func (bc *BanderasComando) Ejecutar(tareas *Tareas) error {
	if bc.Ayuda {
		fmt.Println("Administrador de Tareas - Go")
		fmt.Println("Opciones disponibles:")
		flag.PrintDefaults()
		return nil
	}

	switch {
	case bc.Listar:
		tareas.Imprimir()
		return nil
	case bc.Agregar != "":
		return tareas.Agregar(bc.Agregar)
	case bc.Eliminar != -1:
		return tareas.Eliminar(bc.Eliminar)
	case bc.Alternar != -1:
		return tareas.AlternarEstado(bc.Alternar)
	case bc.Editar != "":
		partes := strings.SplitN(bc.Editar, ":", 2)
		if len(partes) != 2 {
			return errors.New("formato incorrecto para editar. Ejemplo: \"ID:Nuevo Titulo\"")
		}
		idStr, nuevoTitulo := partes[0], partes[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return errors.New("ID invalido para editar")
		}
		return tareas.Editar(id, nuevoTitulo)
	default:
		fmt.Println("Comando invalido. Usa -ayuda para ver los comandos disponibles.")
		return nil
	}
}
