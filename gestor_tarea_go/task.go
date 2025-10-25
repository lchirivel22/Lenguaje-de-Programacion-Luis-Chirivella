package main

import (
	"errors"
	"fmt"
	"time"
)

type EstadoTarea int

const (
	PorHacer EstadoTarea = iota
	Hecha
)

func (e EstadoTarea) String() string {
	switch e {
	case PorHacer:
		return "Por hacer"
	case Hecha:
		return "Hecha"
	default:
		return "Desconocido"
	}
}

type Tarea struct {
	ID           int         `json:"id"`
	Titulo       string      `json:"titulo"`
	Estado       EstadoTarea `json:"estado"`
	CreadaEn     time.Time   `json:"creada_en"`
	CompletadaEn *time.Time  `json:"completada_en,omitempty"`
	EditadaEn    *time.Time  `json:"editada_en,omitempty"`
}

type Tareas []Tarea

func (t *Tareas) GenerarID() int {
	if len(*t) == 0 {
		return 1
	}
	ultimoID := (*t)[len(*t)-1].ID
	return ultimoID + 1
}

func (t *Tareas) Agregar(titulo string) error {
	nuevaTarea := Tarea{
		ID:       t.GenerarID(),
		Titulo:   titulo,
		Estado:   PorHacer,
		CreadaEn: time.Now(),
	}
	*t = append(*t, nuevaTarea)
	return nil
}

func (t *Tareas) validarIndice(id int) error {
	encontrado := false
	for _, tarea := range *t {
		if tarea.ID == id {
			encontrado = true
			break
		}
	}
	if !encontrado {
		return errors.New("ID de tarea inválido")
	}
	return nil
}

func (t *Tareas) Eliminar(id int) error {
	if err := t.validarIndice(id); err != nil {
		return err
	}

	indiceAEliminar := -1
	for i, tarea := range *t {
		if tarea.ID == id {
			indiceAEliminar = i
			break
		}
	}

	if indiceAEliminar == -1 {
		return errors.New("tarea no encontrada")
	}

	*t = append((*t)[:indiceAEliminar], (*t)[indiceAEliminar+1:]...)
	return nil
}

func (t *Tareas) AlternarEstado(id int) error {
	if err := t.validarIndice(id); err != nil {
		return err
	}

	for i := range *t {
		if (*t)[i].ID == id {
			if (*t)[i].Estado == PorHacer {
				(*t)[i].Estado = Hecha
				ahora := time.Now()
				(*t)[i].CompletadaEn = &ahora
			} else {
				(*t)[i].Estado = PorHacer
				(*t)[i].CompletadaEn = nil
			}
			return nil
		}
	}
	return errors.New("tarea no encontrada")
}

func (t *Tareas) Editar(id int, nuevoTitulo string) error {
	if err := t.validarIndice(id); err != nil {
		return err
	}

	for i := range *t {
		if (*t)[i].ID == id {
			(*t)[i].Titulo = nuevoTitulo
			ahora := time.Now()
			(*t)[i].EditadaEn = &ahora
			return nil
		}
	}
	return errors.New("tarea no encontrada")
}

func (t *Tareas) Imprimir() {
	if len(*t) == 0 {
		fmt.Println("--- No existen tareas actualmente ---")
		return
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("%-5s | %-25s | %-12s | %-20s\n", "ID", "Título", "Estado", "Última edición")
	fmt.Println("--------------------------------------------------------------------------------")

	for _, tarea := range *t {
		estadoStr := tarea.Estado.String()

		editadaEnStr := "N/E"
		if tarea.EditadaEn != nil {
			editadaEnStr = tarea.EditadaEn.Format("2006-01-02 15:04")
		}

		fmt.Printf("%-5d | %-25s | %-12s | %-20s\n", tarea.ID, tarea.Titulo, estadoStr, editadaEnStr)
	}
	fmt.Println("--------------------------------------------------------------------------------")
}
