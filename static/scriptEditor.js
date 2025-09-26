document.addEventListener("DOMContentLoaded", () => {
  cargarEditores();

  const form = document.getElementById('formCrearEditor');
  if(form) {
    form.addEventListener('submit', (e) => {
      e.preventDefault();

      const nombre = document.getElementById('nombre').value.trim();
      if (!nombre) {
        mostrarMensajeError("El nombre es obligatorio");
        return;
      }

      fetch('/editores/crear', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: new URLSearchParams({ nombre })
      })
      .then(res => res.json())
      .then(data => {
        if (data.error) {
          mostrarMensajeError(data.error);
        } else {
          mostrarMensajeExito(data.mensaje || "Editor creado exitosamente");
          form.reset();
          cargarEditores();
        }
      })
      .catch(() => mostrarMensajeError("Error al crear editor"));
    });
  }

  const cerrarBtn = document.getElementById('cerrarMensaje');
  if (cerrarBtn) {
    cerrarBtn.addEventListener('click', () => {
      document.getElementById('mensajeBox').style.display = "none";
    });
  }
});

function cargarEditores() {
  fetch("/editores")
    .then(res => res.json())
    .then(data => {
      const tbody = document.querySelector("tbody");
      if (!tbody) return;

      tbody.innerHTML = "";

      if (!Array.isArray(data) || data.length === 0) {
        const row = document.createElement("tr");
        row.innerHTML = `<td colspan="3" class="no-data">No hay editores registrados.</td>`;
        tbody.appendChild(row);
      } else {
        data.forEach(editor => {
          const row = document.createElement("tr");
          row.innerHTML = `
            <td>${editor.ID}</td>
            <td>${editor.Nombre}</td>
            <td>
              <button onclick="eliminarEditor(${editor.ID})">Eliminar</button>
            </td>
          `;
          tbody.appendChild(row);
        });
      }
    })
    .catch(err => {
      console.error("Error al cargar editores:", err);
      const tbody = document.querySelector("tbody");
      if (tbody) tbody.innerHTML = `<tr><td colspan="3" style="color:red;">Error cargando editores</td></tr>`;
    });
}
function eliminarEditor(id) {
  if (!confirm("Al eliminar este editor, también se eliminarán todos los videojuegos asociados. ¿Estás seguro de continuar?")) return;

  fetch('/editores/eliminar', {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: new URLSearchParams({ id })
  })
  .then(res => res.json())
  .then(data => {
    if (!data.error) {
      // 👇 Elimina o comenta esta línea si no quieres mensaje de éxito:
      // mostrarMensajeExito(data.mensaje || "Editor eliminado exitosamente");

      // Solo recarga la lista sin mostrar mensaje
      cargarEditores();
    } else {
      mostrarMensajeError(data.error);
    }
  })
  .catch(err => {
    console.error("Error:", err);
    mostrarMensajeError("Error al eliminar editor.");
  });
}


function mostrarMensajeError(text) {
  const box = document.getElementById('mensajeBox');
  const errorSpan = document.getElementById('mensajeError');
  const exitoSpan = document.getElementById('mensajeExito');

  errorSpan.textContent = text;
  errorSpan.style.display = "block";

  exitoSpan.textContent = "";
  exitoSpan.style.display = "none";

  box.style.display = "block";
}
function mostrarMensajeExito(text) {
  const box = document.getElementById('mensajeBox');
  const errorSpan = document.getElementById('mensajeError');
  const exitoSpan = document.getElementById('mensajeExito');

  exitoSpan.textContent = text;
  exitoSpan.style.display = "block";

  errorSpan.textContent = "";
  errorSpan.style.display = "none";

  box.style.display = "block";
}
