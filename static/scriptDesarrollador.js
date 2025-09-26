document.addEventListener("DOMContentLoaded", () => {
    fetch("/Desarrollador")
        .then(res => res.json())
        .then(data => {
            const tbody = document.querySelector("tbody");
            tbody.innerHTML = ""; // Limpiar contenido actual

            if (data.length === 0) {
                const row = document.createElement("tr");
                row.innerHTML = `<td colspan="3" class="no-data">No hay Desarrolladores registrados.</td>`;
                tbody.appendChild(row);
            } else {
                data.forEach(editores => {
                    const row = document.createElement("tr");
row.innerHTML = `
    <td>${editores.ID}</td>
    <td>${editores.Nombre}</td>
    <td>
        <button onclick="eliminarEditores(${editores.ID})">Eliminar</button>
    </td>
`;
                    tbody.appendChild(row);
                });
            }
        })
        .catch(err => {
            console.error("Error al cargar Desarrolladores:", err);
            const tbody = document.querySelector("tbody");
            tbody.innerHTML = `<tr><td colspan="3" style="color:red;">Error cargando Desarrolladores</td></tr>`;
        });
});

function eliminarDesarrollador(id) {
    if (!confirm("Al eliminar este desarrollador, también se eliminarán todos los videojuegos asociados. ¿Estás seguro de continuar?")) return;

    fetch('/Desarrollador/eliminar', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({ id })
    })
    .then(res => {
        if (res.ok) {
            // Vuelve a cargar la lista
            cargarDesarrollador();
        } else {
            alert("Error al eliminar Desarrollador.");
        }
    })
    .catch(err => {
        console.error("Error:", err);
        alert("Error al eliminar Desarrollador");
    });
}

function cargarDesarrollador() {
    fetch("/Desarrollador")
        .then(res => res.json())
        .then(data => {
            const tbody = document.querySelector("tbody");
            tbody.innerHTML = "";

            if (data.length === 0) {
                const row = document.createElement("tr");
                row.innerHTML = `<td colspan="3" class="no-data">No hay Desarrolladores registrados.</td>`;
                tbody.appendChild(row);
            } else {
                data.forEach(Desarrollador => {
                    const row = document.createElement("tr");
                    row.innerHTML = `
                        <td>${Desarrollador.ID}</td>
                        <td>${Desarrollador.Nombre}</td>
                        <td>
                            <button onclick="eliminarDesarrollador(${Desarrollador.ID})">Eliminar</button>
                        </td>
                    `;
                    tbody.appendChild(row);
                });
            }
        })
        .catch(err => {
            console.error("Error al cargar Desarrollador:", err);
            const tbody = document.querySelector("tbody");
            tbody.innerHTML = `<tr><td colspan="3" style="color:red;">Error cargando Desarrollador</td></tr>`;
        });
}
document.addEventListener("DOMContentLoaded", cargarDesarrollador);