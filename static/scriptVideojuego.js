function verVideojuego(id) {
    window.location.href = `/templates/videojuego_editar.html?id=${id}`;
}

function eliminarVideojuego(id) {
    if (!confirm("¿Estás seguro de eliminar este videojuego?")) return;

    fetch('/videojuegos/eliminar', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: new URLSearchParams({ id })
    })
    .then(res => {
        if (res.ok) {
            location.reload();
        } else {
            alert("Error al eliminar videojuego.");
        }
    })
    .catch(err => {
        console.error("Error:", err);
        alert("Error al eliminar videojuego.");
    });
}

function cargarVideojuegos(generoId = "") {
    const url = generoId ? `/videojuegos?genero_id=${generoId}` : "/videojuegos";
    console.log("Cargando videojuegos con URL:", url);

    fetch(url)
        .then(res => {
            if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);
            return res.json();
        })
        .then(data => {
            console.log("Datos recibidos:", data);
            const tbody = document.getElementById("videojuego-tbody");
            if (!tbody) return;
            tbody.innerHTML = "";

            if (data.length === 0) {
                const row = document.createElement("tr");
                row.innerHTML = `<td colspan="8" class="no-data">No hay videojuegos registrados.</td>`;
                tbody.appendChild(row);
            } else {
                data.forEach(videojuego => {
                    const row = document.createElement("tr");
                    row.innerHTML = `
                        <td>${videojuego.id}</td>
                        <td>${videojuego.titulo}</td>
                        <td>${videojuego.genero_nombre}</td>
                        <td>${videojuego.desarrollador_nombre}</td>
                        <td>${videojuego.editor_nombre}</td>
                        <td>${new Date(videojuego.fecha_lanzamiento).toLocaleDateString()}</td>
                        <td>$${videojuego.precio.toFixed(2)}</td>
                        <td>
                            <button onclick="verVideojuego(${videojuego.id})">Ver</button>
                            <button onclick="eliminarVideojuego(${videojuego.id})">Eliminar</button>
                        </td>
                    `;
                    tbody.appendChild(row);
                });
            }
        })
        .catch(err => {
            console.error("Error al cargar videojuegos:", err);
            const tbody = document.getElementById("videojuego-tbody");
            if (tbody) {
                tbody.innerHTML = `<tr><td colspan="8" style="color:red;">Error cargando videojuegos</td></tr>`;
            }
        });
}

document.addEventListener("DOMContentLoaded", () => {
    console.log("DOM cargado");

    // Vista previa de imagen cuando cambia la URL
    const inputImagenUrl = document.getElementById("imagen_url");
    const previewImagen = document.getElementById("previewImagen");

    if (inputImagenUrl && previewImagen) {
        inputImagenUrl.addEventListener("input", () => {
            const url = inputImagenUrl.value.trim();
            if (url) {
                previewImagen.src = url;
                previewImagen.style.display = "block";
            } else {
                previewImagen.src = "";
                previewImagen.style.display = "none";
            }
        });
    }

    // Cargar selects solo si existen
    const generoSelect = document.getElementById("genero_id");
    if (generoSelect) {
        fetch("/generos")
            .then(res => res.json())
            .then(data => {
                data.forEach(genero => {
                    const option = document.createElement("option");
                    option.value = genero.ID;
                    option.textContent = genero.Nombre;
                    generoSelect.appendChild(option);
                });
            });
    }

    const desarrolladorSelect = document.getElementById("desarrollador_id");
    if (desarrolladorSelect) {
        fetch("/Desarrollador")
            .then(res => res.json())
            .then(data => {
                data.forEach(dev => {
                    const option = document.createElement("option");
                    option.value = dev.ID;
                    option.textContent = dev.Nombre;
                    desarrolladorSelect.appendChild(option);
                });
            });
    }

    const editorSelect = document.getElementById("editor_id");
    if (editorSelect) {
        fetch("/editores")
            .then(res => res.json())
            .then(data => {
                data.forEach(editor => {
                    const option = document.createElement("option");
                    option.value = editor.ID;
                    option.textContent = editor.Nombre;
                    editorSelect.appendChild(option);
                });
            });
    }

    // Manejar mensajes y formulario creación o edición
    const form = document.getElementById("crearVideojuegoForm");
    if (!form) {
        console.log("No existe formulario crearVideojuegoForm, se asume lista");
    }

    const mensajeBox = document.getElementById("mensajeBox");
    const errorDiv = document.getElementById("mensajeError");
    const successDiv = document.getElementById("mensajeExito");
    const cerrarBtn = document.getElementById("cerrarMensaje");

    cerrarBtn?.addEventListener("click", () => {
        mensajeBox.style.display = "none";
        errorDiv.style.display = "none";
        successDiv.style.display = "none";
    });

    // Si hay id en URL, cargar videojuego para editar
    const params = new URLSearchParams(window.location.search);
    const id = params.get("id");
    if (id) {
        fetch(`/videojuego?id=${id}`)
            .then(res => res.json())
            .then(data => {
                document.getElementById("id").value = data.id;
                document.getElementById("titulo").value = data.titulo;
                document.getElementById("descripcion").value = data.descripcion;
                document.getElementById("genero_id").value = data.genero_id;
                document.getElementById("desarrollador_id").value = data.desarrollador_id;
                document.getElementById("editor_id").value = data.editor_id;
                document.getElementById("fecha_lanzamiento").value = data.fecha_lanzamiento.split("T")[0];
                document.getElementById("precio").value = data.precio;
                document.getElementById("imagen_url").value = data.imagen_url;
                document.getElementById("created_at").value = new Date(data.created_at).toLocaleString();
                document.getElementById("updated_at").value = new Date(data.updated_at).toLocaleString();
                previewImagen.src = data.imagen_url;
                previewImagen.style.display = data.imagen_url ? "block" : "none";
            });
    }

    if (form) {
        form.addEventListener("submit", function (event) {
            event.preventDefault();

            errorDiv.style.display = "none";
            successDiv.style.display = "none";
            mensajeBox.style.display = "none";

            const formData = new FormData(form);

            // Determinar URL según si es crear o editar
            const url = id ? "/videojuegos/editar" : "/videojuegos/crear";

            fetch(url, {
                method: "POST",
                body: formData
            })
            .then(res => res.json().then(data => ({ status: res.status, body: data })))
            .then(({ status, body }) => {
                if (status >= 400) {
                    errorDiv.textContent = body.error || "Error desconocido";
                    errorDiv.style.display = "block";
                    successDiv.style.display = "none";
                    mensajeBox.style.display = "block";
                } else {
                    successDiv.textContent = body.mensaje || (id ? "Actualizado correctamente" : "Guardado exitosamente");
                    successDiv.style.display = "block";
                    errorDiv.style.display = "none";
                    mensajeBox.style.display = "block";

                    if (!id) form.reset();
                    else window.location.href = "/templates/videojuego_lista.html";
                }
            })
            .catch(() => {
                errorDiv.textContent = "Error al conectar con el servidor";
                errorDiv.style.display = "block";
                successDiv.style.display = "none";
                mensajeBox.style.display = "block";
            });
        });
    }

    const filtroGenero = document.getElementById("filtroGenero");
    console.log("Filtro género encontrado:", filtroGenero);

    if (filtroGenero) {
        filtroGenero.innerHTML = "";

        const optionTodos = document.createElement("option");
        optionTodos.value = "";
        optionTodos.textContent = "Todos";
        filtroGenero.appendChild(optionTodos);

        fetch("/generos")
            .then(res => res.json())
            .then(generos => {
                console.log("Géneros recibidos:", generos);
                filtroGenero.innerHTML = "";  // Limpia antes de añadir
                filtroGenero.appendChild(optionTodos);  // Añadimos opción "Todos" de nuevo

                generos.forEach(genero => {
                    const option = document.createElement("option");
                    option.value = genero.ID;
                    option.textContent = genero.Nombre;
                    filtroGenero.appendChild(option);
                });

                console.log("Opciones cargadas en filtro:", filtroGenero.options.length);

                cargarVideojuegos(filtroGenero.value);
            })
            .catch(err => console.error("Error al cargar géneros:", err));

        filtroGenero.addEventListener("change", () => {
            console.log("Filtro cambiado a:", filtroGenero.value);
            cargarVideojuegos(filtroGenero.value);
        });
    } else {
        console.warn("No se encontró el filtroGenero");
        cargarVideojuegos();
    }
});
