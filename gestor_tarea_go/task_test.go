package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGenerarID(t *testing.T) {
	tareas := Tareas{}
	if id := tareas.GenerarID(); id != 1 {
		t.Errorf("Esperado ID 1 para lista vacia, obtenido %d", id)
	}

	tareas = Tareas{
		{ID: 1, Titulo: "Tarea 1"},
		{ID: 5, Titulo: "Tarea 2"},
	}
	if id := tareas.GenerarID(); id != 6 {
		t.Errorf("Esperado ID 6, obtenido %d", id)
	}
}

func TestAgregar(t *testing.T) {
	tareas := Tareas{}
	err := tareas.Agregar("Nueva tarea")
	if err != nil {
		t.Fatalf("Error inesperado al agregar tarea: %v", err)
	}
	if len(tareas) != 1 {
		t.Errorf("Esperado 1 tarea, obtenido %d", len(tareas))
	}
	if tareas[0].Titulo != "Nueva tarea" {
		t.Errorf("Titulo de tarea incorrecto")
	}
	if tareas[0].Estado != PorHacer {
		t.Errorf("Estado de tarea incorrecto, esperado PorHacer")
	}
	if tareas[0].ID != 1 {
		t.Errorf("ID de tarea incorrecto, esperado 1")
	}
}

func TestEliminar(t *testing.T) {
	tareas := Tareas{
		{ID: 1, Titulo: "Tarea 1"},
		{ID: 2, Titulo: "Tarea 2"},
	}

	err := tareas.Eliminar(1)
	if err != nil {
		t.Fatalf("Error inesperado al eliminar tarea: %v", err)
	}
	if len(tareas) != 1 {
		t.Errorf("Esperado 1 tarea, obtenido %d", len(tareas))
	}
	if tareas[0].ID != 2 {
		t.Errorf("La tarea incorrecta fue eliminada")
	}

	err = tareas.Eliminar(99)
	if err == nil {
		t.Error("Esperado error al eliminar tarea inexistente, obtenido nil")
	}
	if err.Error() != "ID de tarea inválido" {
		t.Errorf("Mensaje de error incorrecto: %s", err.Error())
	}
}

func TestAlternarEstado(t *testing.T) {
	tareas := Tareas{
		{ID: 1, Titulo: "Tarea 1", Estado: PorHacer},
	}

	err := tareas.AlternarEstado(1)
	if err != nil {
		t.Fatalf("Error inesperado al alternar estado: %v", err)
	}
	if tareas[0].Estado != Hecha {
		t.Errorf("Esperado estado Hecha, obtenido %s", tareas[0].Estado.String())
	}
	if tareas[0].CompletadaEn == nil {
		t.Errorf("CompletadaEn no deberia ser nil")
	}

	err = tareas.AlternarEstado(1)
	if err != nil {
		t.Fatalf("Error inesperado al alternar estado: %v", err)
	}
	if tareas[0].Estado != PorHacer {
		t.Errorf("Esperado estado PorHacer, obtenido %s", tareas[0].Estado.String())
	}
	if tareas[0].CompletadaEn != nil {
		t.Errorf("CompletadaEn deberia ser nil")
	}

	err = tareas.AlternarEstado(99)
	if err == nil {
		t.Error("Esperado error al alternar estado de tarea inexistente, obtenido nil")
	}
	if err.Error() != "ID de tarea inválido" {
		t.Errorf("Mensaje de error incorrecto: %s", err.Error())
	}
}

func TestEditar(t *testing.T) {
	tareas := Tareas{
		{ID: 1, Titulo: "Tarea Vieja", Estado: PorHacer},
	}

	err := tareas.Editar(1, "Tarea Nueva")
	if err != nil {
		t.Fatalf("Error inesperado al editar tarea: %v", err)
	}
	if tareas[0].Titulo != "Tarea Nueva" {
		t.Errorf("Titulo de tarea no actualizado correctamente")
	}
	if tareas[0].EditadaEn == nil {
		t.Errorf("EditadaEn no deberia ser nil")
	}

	err = tareas.Editar(99, "Otro titulo")
	if err == nil {
		t.Error("Esperado error al editar tarea inexistente, obtenido nil")
	}
	if err.Error() != "ID de tarea inválido" {
		t.Errorf("Mensaje de error incorrecto: %s", err.Error())
	}
}

func TestAlmacenamientoGuardarCargar(t *testing.T) {
	nombreArchivo := "test_tareas.json"
	defer os.Remove(nombreArchivo)

	almacenamiento := NuevoAlmacenamiento[Tareas](nombreArchivo)

	tareasOriginales := Tareas{
		{ID: 1, Titulo: "Tarea Prueba 1", Estado: PorHacer, CreadaEn: time.Now()},
		{ID: 2, Titulo: "Tarea Prueba 2", Estado: Hecha, CreadaEn: time.Now(), CompletadaEn: func() *time.Time { t := time.Now(); return &t }()},
	}

	err := almacenamiento.Guardar(&tareasOriginales)
	if err != nil {
		t.Fatalf("Error al guardar tareas: %v", err)
	}

	var tareasCargadas Tareas
	err = almacenamiento.Cargar(&tareasCargadas)
	if err != nil {
		t.Fatalf("Error al cargar tareas: %v", err)
	}

	if len(tareasOriginales) != len(tareasCargadas) {
		t.Fatalf("Longitud de tareas cargadas incorrecta. Esperado %d, obtenido %d", len(tareasOriginales), len(tareasCargadas))
	}
	for i := range tareasOriginales {
		if tareasOriginales[i].ID != tareasCargadas[i].ID ||
			tareasOriginales[i].Titulo != tareasCargadas[i].Titulo ||
			tareasOriginales[i].Estado != tareasCargadas[i].Estado ||
			!tareasOriginales[i].CreadaEn.Equal(tareasCargadas[i].CreadaEn) {
			t.Errorf("Tarea %d difiere.\nEsperado: %+v\nObtenido: %+v", i, tareasOriginales[i], tareasCargadas[i])
		}

		if (tareasOriginales[i].CompletadaEn == nil && tareasCargadas[i].CompletadaEn != nil) ||
			(tareasOriginales[i].CompletadaEn != nil && tareasCargadas[i].CompletadaEn == nil) {
			t.Errorf("CompletadaEn difiere para la tarea %d", i)
		}
		if tareasOriginales[i].CompletadaEn != nil && tareasCargadas[i].CompletadaEn != nil {
			if !tareasOriginales[i].CompletadaEn.Equal(*tareasCargadas[i].CompletadaEn) {
				t.Errorf("CompletadaEn no coincide para la tarea %d", i)
			}
		}

		if (tareasOriginales[i].EditadaEn == nil && tareasCargadas[i].EditadaEn != nil) ||
			(tareasOriginales[i].EditadaEn != nil && tareasCargadas[i].EditadaEn == nil) {
			t.Errorf("EditadaEn difiere para la tarea %d", i)
		}
		if tareasOriginales[i].EditadaEn != nil && tareasCargadas[i].EditadaEn != nil {
			if !tareasOriginales[i].EditadaEn.Equal(*tareasCargadas[i].EditadaEn) {
				t.Errorf("EditadaEn no coincide para la tarea %d", i)
			}
		}
	}
}

func TestCargarArchivoVacio(t *testing.T) {
	nombreArchivo := "test_vacio.json"
	defer os.Remove(nombreArchivo)
	err := ioutil.WriteFile(nombreArchivo, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Error al crear archivo vacio: %v", err)
	}

	almacenamiento := NuevoAlmacenamiento[Tareas](nombreArchivo)
	var tareasCargadas Tareas
	err = almacenamiento.Cargar(&tareasCargadas)
	if err != nil {
		t.Fatalf("Error al cargar archivo vacio: %v", err)
	}
	if len(tareasCargadas) != 0 {
		t.Errorf("Esperado 0 tareas al cargar archivo vacio, obtenido %d", len(tareasCargadas))
	}
}

func TestCargarArchivoInexistente(t *testing.T) {
	nombreArchivo := "test_inexistente.json"
	os.Remove(nombreArchivo)

	almacenamiento := NuevoAlmacenamiento[Tareas](nombreArchivo)
	var tareasCargadas Tareas
	err := almacenamiento.Cargar(&tareasCargadas)
	if err != nil {
		t.Fatalf("Error al cargar archivo inexistente: %v", err)
	}
	if len(tareasCargadas) != 0 {
		t.Errorf("Esperado 0 tareas al cargar archivo inexistente, obtenido %d", len(tareasCargadas))
	}
}

func TestCargarJSONInvalido(t *testing.T) {
	nombreArchivo := "test_invalido.json"
	defer os.Remove(nombreArchivo)
	err := ioutil.WriteFile(nombreArchivo, []byte("{esto no es json valido"), 0644)
	if err != nil {
		t.Fatalf("Error al crear archivo JSON invalido: %v", err)
	}

	almacenamiento := NuevoAlmacenamiento[Tareas](nombreArchivo)
	var tareasCargadas Tareas
	err = almacenamiento.Cargar(&tareasCargadas)
	if err == nil {
		t.Error("Esperado error al cargar JSON invalido, obtenido nil")
	}
	if !strings.Contains(err.Error(), "error al deserializar datos") {
		t.Errorf("Mensaje de error incorrecto para JSON invalido: %s", err.Error())
	}
}

func TestAlmacenamientoGuardarLeerPermisos(t *testing.T) {
	nombreArchivo := "test_permisos.json"
	defer os.Remove(nombreArchivo)

	almacenamiento := NuevoAlmacenamiento[Tareas](nombreArchivo)
	tareasOriginales := Tareas{{ID: 1, Titulo: "Permiso", Estado: PorHacer, CreadaEn: time.Now()}}

	err := almacenamiento.Guardar(&tareasOriginales)
	if err != nil {
		t.Fatalf("Error al guardar tareas con permisos: %v", err)
	}

	_, err = os.Stat(nombreArchivo)
	if err != nil {
		t.Fatalf("El archivo no fue creado o no se puede acceder a el: %v", err)
	}
}

func TestStringEstadoTarea(t *testing.T) {
	tests := []struct {
		estado   EstadoTarea
		esperado string
	}{
		{PorHacer, "Por hacer"},
		{Hecha, "Hecha"},
		{EstadoTarea(99), "Desconocido"},
	}

	for _, tc := range tests {
		obtenido := tc.estado.String()
		if obtenido != tc.esperado {
			t.Errorf("Para estado %d, esperado '%s', obtenido '%s'", tc.estado, tc.esperado, obtenido)
		}
	}
}
