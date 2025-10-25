use std::{
    fs::{self, File},
    io::{self, Write},
    path::Path,
};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, Clone)]
enum Estado {
    Pendiente,
    EnProceso,
    Completada,
}

impl std::fmt::Display for Estado {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Estado::Pendiente => write!(f, "Pendiente"),
            Estado::EnProceso => write!(f, "En Proceso"),
            Estado::Completada => write!(f, "Completada"),
        }
    }
}

#[derive(Debug, Serialize, Deserialize, Clone)]
struct Tarea {
    id: usize,
    descripcion: String,
    estado: Estado,
}

const RUTA_ARCHIVO: &str = "tareas.json";

fn cargar_tareas() -> Vec<Tarea> {
    if !Path::new(RUTA_ARCHIVO).exists() {
        return Vec::new();
    }
    let contenido = match fs::read_to_string(RUTA_ARCHIVO) {
        Ok(c) => c,
        Err(_) => {
            eprintln!("Error al leer el archivo.");
            return Vec::new();
        }
    };
    match serde_json::from_str(&contenido) {
        Ok(tareas) => tareas,
        Err(_) => {
            eprintln!("Error al parsear JSON.");
            Vec::new()
        }
    }
}

fn guardar_tareas(tareas: &[Tarea]) -> io::Result<()> {
    let json_data = serde_json::to_string_pretty(tareas)?;
    let mut file = File::create(RUTA_ARCHIVO)?;
    file.write_all(json_data.as_bytes())?;
    Ok(())
}

fn listar_tareas(tareas: &[Tarea]) {
    if tareas.is_empty() {
        println!("\nNo hay tareas.");
        return;
    }
    println!("\nLista de Tareas:");
    for tarea in tareas {
        println!("[ID: {}] {} - Estado: {}", tarea.id, tarea.descripcion, tarea.estado);
    }
}

fn añadir_tarea(tareas: &mut Vec<Tarea>, descripcion: String) {
    let nuevo_id = tareas.iter().map(|t| t.id).max().unwrap_or(0) + 1;
    let nueva_tarea = Tarea {
        id: nuevo_id,
        descripcion,
        estado: Estado::Pendiente,
    };
    tareas.push(nueva_tarea);
    println!("Tarea añadida con ID: {}", nuevo_id);
}

fn editar_descripcion(tareas: &mut Vec<Tarea>, id: usize, nueva_desc: String) {
    if let Some(tarea) = tareas.iter_mut().find(|t| t.id == id) {
        tarea.descripcion = nueva_desc;
        println!("Tarea ID {} descripción actualizada.", id);
    } else {
        println!("No se encontró la tarea con ID: {}.", id);
    }
}

fn cambiar_estado(tareas: &mut Vec<Tarea>, id: usize, nuevo_estado: &str) {
    if let Some(tarea) = tareas.iter_mut().find(|t| t.id == id) {
        tarea.estado = match nuevo_estado.to_lowercase().as_str() {
            "pendiente" => Estado::Pendiente,
            "proceso" => Estado::EnProceso,
            "completada" | "completa" => Estado::Completada,
            _ => {
                println!("Estado inválido. Use 'pendiente', 'proceso' o 'completada'.");
                return;
            }
        };
        println!("Tarea ID {} estado cambiado a {}.", id, tarea.estado);
    } else {
        println!("No se encontró la tarea con ID: {}.", id);
    }
}

fn eliminar_tarea(tareas: &mut Vec<Tarea>, id: usize) {
    let inicial_len = tareas.len();
    tareas.retain(|t| t.id != id);
    if tareas.len() < inicial_len {
        println!("Tarea ID {} eliminada.", id);
    } else {
        println!("No se encontró la tarea con ID: {}.", id);
    }
}

fn mostrar_ayuda() {
    println!("\nComandos Disponibles:");
    println!("-------------------------------------------------------------");
    println!("cargo run -- listar                 -> Muestra todas las tareas.");
    println!("cargo run -- añadir <descripcion>   -> Añade una tarea nueva.");
    println!("cargo run -- editar <id> <nueva_desc> -> Edita la descripción de una tarea.");
    println!("cargo run -- estado <id> <estado>   -> Cambia el estado (pendiente|proceso|completada).");
    println!("cargo run -- eliminar <id>          -> Elimina una tarea por su ID.");
    println!("cargo run -- ayuda                  -> Muestra esta ayuda.");
    println!("-------------------------------------------------------------\n");
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut tareas = cargar_tareas();
    let args: Vec<String> = std::env::args().collect();

    if args.len() < 2 {
        mostrar_ayuda();
        return Ok(());
    }

    let comando = &args[1];

    match comando.as_str() {
        "listar" => listar_tareas(&tareas),

        "añadir" => {
            if args.len() < 3 {
                println!("Uso: cargo run -- añadir <descripcion>");
            } else {
                let descripcion = args[2..].join(" ");
                añadir_tarea(&mut tareas, descripcion);
                guardar_tareas(&tareas)?;
            }
        }

        "editar" => {
            if args.len() < 4 {
                println!("Uso: cargo run -- editar <id> <nueva_descripcion>");
            } else {
                let id: usize = match args[2].parse() {
                    Ok(n) => n,
                    Err(_) => {
                        println!("ID inválido. Debe ser un número.");
                        return Ok(());
                    }
                };
                let nueva_desc = args[3..].join(" ");
                editar_descripcion(&mut tareas, id, nueva_desc);
                guardar_tareas(&tareas)?;
            }
        }

        "estado" => {
            if args.len() < 4 {
                println!("Uso: cargo run -- estado <id> <nuevo_estado>");
            } else {
                let id: usize = match args[2].parse() {
                    Ok(n) => n,
                    Err(_) => {
                        println!("ID inválido. Debe ser un número.");
                        return Ok(());
                    }
                };
                let nuevo_estado = &args[3];
                cambiar_estado(&mut tareas, id, nuevo_estado);
                guardar_tareas(&tareas)?;
            }
        }

        "eliminar" => {
            if args.len() < 3 {
                println!("Uso: cargo run -- eliminar <id>");
            } else {
                let id: usize = match args[2].parse() {
                    Ok(n) => n,
                    Err(_) => {
                        println!("ID inválido. Debe ser un número.");
                        return Ok(());
                    }
                };
                eliminar_tarea(&mut tareas, id);
                guardar_tareas(&tareas)?;
            }
        }

        "ayuda" | _ => mostrar_ayuda(),
    }

    Ok(())
}