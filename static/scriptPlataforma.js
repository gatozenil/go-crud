document.addEventListener("DOMContentLoaded", () => {
    fetch("/Plataforma")
        .then(res => res.json())
        .then(data => {
            const tbody = document.querySelector("tbody");
            tbody.innerHTML = ""; // Limpiar contenido actual

            if (data.length === 0) {
                const row = document.createElement("tr");
                row.innerHTML = `<td colspan="3" class="no-data">No hay Plataformas registradas.</td>`;
                tbody.appendChild(row);
            } else {
                data.forEach(genero => {
                    const row = document.createElement("tr");
row.innerHTML = `
    <td>${Plataforma.ID}</td>
    <td>${Plataforma.Nombre}</td>
    <td>
        <button onclick="eliminarPlataforma(${Plataforma.ID})">Eliminar</button>
    </td>
`;
                    tbody.appendChild(row);
                });
            }
        })
        .catch(err => {
            console.error("Error al cargar Plataformas:", err);
            const tbody = document.querySelector("tbody");
            tbody.innerHTML = `<tr><td colspan="3" style="color:red;">Error cargando Plataformas</td></tr>`;
        });
});

function eliminarPlataforma(id) {
    if (!confirm("Al eliminar esta platafoma, también se eliminarán todos los videojuegos asociados. ¿Estás seguro de continuar?")) return;

    fetch('/Plataforma/eliminar', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({ id })
    })
    .then(res => {
        if (res.ok) {
            // Vuelve a cargar la lista
            cargarPlataforma();
        } else {
            alert("Error al eliminar Plataforma.");
        }
    })
    .catch(err => {
        console.error("Error:", err);
        alert("Error al eliminar Plataforma.");
    });
}

function cargarPlataforma() {
    fetch("/Plataforma")
        .then(res => res.json())
        .then(data => {
            const tbody = document.querySelector("tbody");
            tbody.innerHTML = "";

            if (data.length === 0) {
                const row = document.createElement("tr");
                row.innerHTML = `<td colspan="3" class="no-data">No hay Plataformas registrados.</td>`;
                tbody.appendChild(row);
            } else {
                data.forEach(Plataforma => {
                    const row = document.createElement("tr");
                    row.innerHTML = `
                        <td>${Plataforma.ID}</td>
                        <td>${Plataforma.Nombre}</td>
                        <td>
                            <button onclick="eliminarPlataforma(${Plataforma.ID})">Eliminar</button>
                        </td>
                    `;
                    tbody.appendChild(row);
                });
            }
        })
        .catch(err => {
            console.error("Error al cargar Plataformas:", err);
            const tbody = document.querySelector("tbody");
            tbody.innerHTML = `<tr><td colspan="3" style="color:red;">Error cargando Plataformas</td></tr>`;
        });
}

document.addEventListener("DOMContentLoaded", cargarPlataforma);