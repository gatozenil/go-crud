document.addEventListener("DOMContentLoaded", () => {
  cargarGeneros();

  const form = document.getElementById('formCrearGenero');
  form.addEventListener('submit', (e) => {
    e.preventDefault();

    const nombre = document.getElementById('nombre').value.trim();
    if (!nombre) {
      mostrarMensajeError("El nombre es obligatorio");
      return;
    }

    fetch('/genero/crear', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: new URLSearchParams({ nombre })
    })
    .then(res => res.json())
    .then(data => {
      if (data.error) {
        mostrarMensajeError(data.error);
      } else {
        mostrarMensajeExito(data.mensaje || "Género creado exitosamente");
        form.reset();
        cargarGeneros();
      }
    })
    .catch(() => mostrarMensajeError("Error al crear género"));
  });

  document.getElementById('cerrarMensaje').addEventListener('click', () => {
    document.getElementById('mensajeBox').style.display = "none";
  });
});

// Función para cargar géneros y mostrar en la tabla
function cargarGeneros() {
  fetch("/generos")
    .then(res => res.json())
    .then(data => {
      const tbody = document.querySelector("tbody");
      if (!tbody) return;  

      tbody.innerHTML = "";

      if (!Array.isArray(data) || data.length === 0) {
        const row = document.createElement("tr");
        row.innerHTML = `<td colspan="3" class="no-data">No hay géneros registrados.</td>`;
        tbody.appendChild(row);
      } else {
        data.forEach(genero => {
          const row = document.createElement("tr");
          row.innerHTML = `
            <td>${genero.ID}</td>
            <td>${genero.Nombre}</td>
            <td>
              <button onclick="eliminarGenero(${genero.ID})">Eliminar</button>
            </td>
          `;
          tbody.appendChild(row);
        });
      }
    })
    .catch(err => {
      console.error("Error al cargar géneros:", err);
      const tbody = document.querySelector("tbody");
      if (tbody) tbody.innerHTML = `<tr><td colspan="3" style="color:red;">Error cargando géneros</td></tr>`;
    });
}

// Funciones para mostrar mensajes
function mostrarMensajeError(text) {
  console.log("Mostrar error:", text);
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
  console.log("Mostrar éxito:", text);
  const box = document.getElementById('mensajeBox');
  const errorSpan = document.getElementById('mensajeError');
  const exitoSpan = document.getElementById('mensajeExito');

  exitoSpan.textContent = text;
  exitoSpan.style.display = "block";    

  errorSpan.textContent = "";
  errorSpan.style.display = "none";    

  box.style.display = "block";          
}
// Función para eliminar un género
function eliminarGenero(id) {
  if (!confirm("Al eliminar este genero, también se eliminarán todos los videojuegos asociados. ¿Estás seguro de continuar")) return;

  fetch('/genero/eliminar', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    body: new URLSearchParams({ id })
  })
  .then(res => res.json())
  .then(data => {
    if (!data.error) {
      cargarGeneros();  
    }
    else {
      console.error("Error al eliminar género:", data.error);
    }
  })
  .catch(err => {
    console.error("Error en fetch:", err);
  });
}
