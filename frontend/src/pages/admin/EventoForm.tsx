import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { api } from "../../api";

interface TipoEntrada {
  id: number;
  nombre: string;
  precio: number;
  stock_disponible: number;
}

export default function EventoForm() {
  const { id } = useParams();
  const navigate = useNavigate();
  const esEdicion = !!id;

  const [form, setForm] = useState({
    nombre: "", descripcion: "", fecha: "", lugar: "", estado: "ACTIVO",
  });
  const [error, setError] = useState("");
  const [cargando, setCargando] = useState(false);
  const [eventoID, setEventoID] = useState<number | null>(null);

  // Tipos de entrada
  const [tipos, setTipos] = useState<TipoEntrada[]>([]);
  const [nuevoTipo, setNuevoTipo] = useState({ nombre: "", precio: "", stock_disponible: "" });
  const [errorTipo, setErrorTipo] = useState("");

  const cargarTipos = async (eid: number) => {
    try {
      const data = await api.get(`/eventos/${eid}/tipos-entrada`);
      setTipos(data || []);
    } catch {
      setTipos([]);
    }
  };

  useEffect(() => {
    if (esEdicion) {
      api.get(`/eventos/${id}`).then((data) => {
        setForm({
          nombre: data.nombre || "",
          descripcion: data.descripcion || "",
          fecha: data.fecha ? new Date(data.fecha).toISOString().slice(0, 16) : "",
          lugar: data.lugar || "",
          estado: data.estado || "ACTIVO",
        });
        setEventoID(Number(id));
        cargarTipos(Number(id));
      });
    }
  }, [id, esEdicion]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>
  ) => setForm({ ...form, [e.target.name]: e.target.value });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    if (!form.nombre.trim()) { setError("El nombre es obligatorio"); return; }
    if (!form.lugar.trim()) { setError("El lugar es obligatorio"); return; }
    if (!form.fecha) { setError("La fecha es obligatoria"); return; }

    try {
      setCargando(true);
      const body = { ...form, fecha: new Date(form.fecha).toISOString() };

      if (esEdicion) {
        await api.put(`/admin/eventos/${id}`, body);
        navigate("/admin");
      } else {
        const data = await api.post("/admin/eventos", body);
        // Quedarse en la página para agregar tipos de entrada
        const nuevoID = data.ID || data.id;
        setEventoID(nuevoID);
      }
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setCargando(false);
    }
  };

  const handleAgregarTipo = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorTipo("");

    if (!nuevoTipo.nombre.trim()) { setErrorTipo("El nombre es obligatorio"); return; }
    if (!nuevoTipo.precio || parseFloat(nuevoTipo.precio) <= 0) { setErrorTipo("El precio debe ser mayor a 0"); return; }
    if (!nuevoTipo.stock_disponible || parseInt(nuevoTipo.stock_disponible) < 0) { setErrorTipo("El stock no puede ser negativo"); return; }

    try {
      await api.post(`/admin/eventos/${eventoID}/tipos-entrada`, {
        nombre: nuevoTipo.nombre,
        precio: parseFloat(nuevoTipo.precio),
        stock_disponible: parseInt(nuevoTipo.stock_disponible),
      });
      setNuevoTipo({ nombre: "", precio: "", stock_disponible: "" });
      cargarTipos(eventoID!);
    } catch (err) {
      setErrorTipo((err as Error).message);
    }
  };

  const handleEliminarTipo = async (tipoId: number) => {
    if (!confirm("¿Eliminar este tipo de entrada?")) return;
    try {
      await api.delete(`/admin/eventos/${eventoID}/tipos-entrada/${tipoId}`);
      cargarTipos(eventoID!);
    } catch (err) {
      setErrorTipo((err as Error).message);
    }
  };

  return (
    <div className="container">
      <h1>{esEdicion ? "Editar Evento" : "Nuevo Evento"}</h1>

      {/* Formulario principal del evento */}
      {!eventoID || esEdicion ? (
        <>
          {error && <p className="error">{error}</p>}
          <div className="formulario">
            <form onSubmit={handleSubmit}>
              <label>Nombre *</label>
              <input name="nombre" value={form.nombre} onChange={handleChange}
                placeholder="Nombre del evento" required />
              <label>Descripción</label>
              <textarea name="descripcion" value={form.descripcion} onChange={handleChange}
                placeholder="Descripción" rows={3}
                style={{ padding: "0.6rem", border: "1px solid #ccc", borderRadius: "6px" }} />
              <label>Fecha y hora *</label>
              <input type="datetime-local" name="fecha" value={form.fecha}
                onChange={handleChange} required />
              <label>Lugar *</label>
              <input name="lugar" value={form.lugar} onChange={handleChange}
                placeholder="Lugar del evento" required />
              <label>Estado</label>
              <select name="estado" value={form.estado} onChange={handleChange}>
                <option value="ACTIVO">ACTIVO</option>
                <option value="CANCELADO">CANCELADO</option>
              </select>
              <div style={{ display: "flex", gap: "0.5rem", marginTop: "0.5rem" }}>
                <button type="submit" disabled={cargando}>
                  {cargando ? "Guardando..." : esEdicion ? "Guardar cambios" : "Crear evento"}
                </button>
                <button type="button" onClick={() => navigate("/admin")} className="btn-secundario">
                  Cancelar
                </button>
              </div>
            </form>
          </div>
        </>
      ) : null}

      {/* Sección de tipos de entrada — aparece después de crear o al editar */}
      {eventoID && (
        <div style={{ marginTop: "2rem" }}>
          <h2>Tipos de entrada</h2>
          <p style={{ color: "#666", marginBottom: "1rem" }}>
            Definí las categorías de entradas disponibles para este evento (ej: Campo, Platea, VIP).
          </p>

          {errorTipo && <p className="error">{errorTipo}</p>}

          {/* Tabla de tipos existentes */}
          {tipos.length > 0 && (
            <table className="tabla-admin" style={{ marginBottom: "1.5rem" }}>
              <thead>
                <tr>
                  <th>Nombre</th><th>Precio</th><th>Stock</th><th></th>
                </tr>
              </thead>
              <tbody>
                {tipos.map((t) => (
                  <tr key={t.id}>
                    <td>{t.nombre}</td>
                    <td>${t.precio}</td>
                    <td>{t.stock_disponible}</td>
                    <td>
                      <button onClick={() => handleEliminarTipo(t.id)} className="btn-peligro">
                        Eliminar
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}

          {/* Formulario para agregar tipo */}
          <div className="card" style={{ padding: "1.25rem" }}>
            <h3 style={{ marginBottom: "1rem" }}>Agregar tipo de entrada</h3>
            <form onSubmit={handleAgregarTipo}
              style={{ display: "flex", gap: "0.75rem", flexWrap: "wrap", alignItems: "flex-end" }}>
              <div style={{ display: "flex", flexDirection: "column", gap: "0.3rem", flex: 2 }}>
                <label style={{ fontSize: "0.85rem" }}>Nombre</label>
                <input placeholder="ej: Campo, Platea, VIP"
                  value={nuevoTipo.nombre}
                  onChange={(e) => setNuevoTipo({ ...nuevoTipo, nombre: e.target.value })} />
              </div>
              <div style={{ display: "flex", flexDirection: "column", gap: "0.3rem", flex: 1 }}>
                <label style={{ fontSize: "0.85rem" }}>Precio ($)</label>
                <input type="number" placeholder="5000"
                  value={nuevoTipo.precio}
                  onChange={(e) => setNuevoTipo({ ...nuevoTipo, precio: e.target.value })} />
              </div>
              <div style={{ display: "flex", flexDirection: "column", gap: "0.3rem", flex: 1 }}>
                <label style={{ fontSize: "0.85rem" }}>Stock</label>
                <input type="number" placeholder="100"
                  value={nuevoTipo.stock_disponible}
                  onChange={(e) => setNuevoTipo({ ...nuevoTipo, stock_disponible: e.target.value })} />
              </div>
              <button type="submit" style={{ alignSelf: "flex-end" }}>Agregar</button>
            </form>
          </div>

          <div style={{ marginTop: "1.5rem" }}>
            <button onClick={() => navigate("/admin")}>
              ✓ Terminar y volver al panel
            </button>
          </div>
        </div>
      )}
    </div>
  );
}