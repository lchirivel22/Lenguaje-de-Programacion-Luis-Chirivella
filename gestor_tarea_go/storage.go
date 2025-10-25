
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Almacenamiento[T any] struct {
	NombreArchivo string
}

func NuevoAlmacenamiento[T any](nombreArchivo string) *Almacenamiento[T] {
	return &Almacenamiento[T]{
		NombreArchivo: nombreArchivo,
	}
}

func (a *Almacenamiento[T]) Guardar(datos *T) error {
	contenidoJSON, err := json.MarshalIndent(datos, "", "    ")
	if err != nil {
		return fmt.Errorf("error al serializar datos: %w", err)
	}

	err = ioutil.WriteFile(a.NombreArchivo, contenidoJSON, 0644)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo %s: %w", a.NombreArchivo, err)
	}
	return nil
}

func (a *Almacenamiento[T]) Cargar(datos *T) error {
	_, err := os.Stat(a.NombreArchivo)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error al verificar el archivo %s: %w", a.NombreArchivo, err)
	}

	contenidoJSON, err := ioutil.ReadFile(a.NombreArchivo)
	if err != nil {
		return fmt.Errorf("error al leer el archivo %s: %w", a.NombreArchivo, err)
	}

	if len(contenidoJSON) == 0 {
		return nil
	}

	err = json.Unmarshal(contenidoJSON, datos)
	if err != nil {
		return fmt.Errorf("error al deserializar datos: %w", err)
	}
	return nil
}