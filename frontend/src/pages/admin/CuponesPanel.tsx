import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { api } from "../../api";

interface Cupon {
  ID: number;
  Codigo: string;
  TipoDescuento: string;
  ValorDescuento: number;
  FechaVencimiento: string;
  LimiteUsos: number;
  UsosActuales: number;
  Estado: string;
  EventoID: number;
}

export default function CuponesPanel() {
  const navigate = useNavigate();
  const [cupones, setCupones] = useState<Cupon[]>([]);
  const [error, setError] = useState("");
  const [mensaje, setMensaje] = useState("");
  const [mostrarForm, setMostrarForm] = useState(false);
  const [form, setForm] = useState({
    codigo: "",
    tipo_descuento: "PORCENTAJE",
    valor_descuento: "",
    fecha_vencimiento: "",
    limite_usos: "",
    evento_id: "",
  });

  const cargar = async () => {
    try {
      const data = await api.get("/admin/cupones");
      setCupones(data || []);
    } catch (err) {
      setError((err as Error).message);
    }
  };

  useEffect(() => { cargar(); }, []);

  const handleCrear = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!form.evento_id || parseInt(form.evento_id) <= 0) {
      setError("El ID de evento es obligatorio (el modelo requiere un evento asociado)");
      return;
    }

    try {
      await api.post("/admin/cupones", {
        codigo: form.codigo,
        tipo_descuento: form.tipo_descuento,
        valor_descuento: parseFloat(form.valor_descuento),
        fecha_vencimiento: new Date(form.fecha_vencimiento).toISOString(),
        limite_usos: parseInt(form.limite_usos),
        evento_id: parseInt(form.evento_id),
      });
      setMensaje("Cupón creado con éxito");
      setMostrarForm(false);
      setForm({ codigo: "", tipo_descuento: "PORCENTAJE", valor_descuento: "", fecha_vencimiento: "", limite_usos: "", evento_id: "" });
      cargar();
    } catch (err) {
      setError((err as Error).message);
    }
  };

  const handleDesactivar = async (id: number) => {
    if (!confirm("¿Desactivar este cupón?")) return;
    try {
      await api.put(`/admin/cupones/${id}/desactivar`);
      setMensaje("Cupón desactivado");
      cargar();
    } catch (err) {
      setError((err as Error).message);
    }
  };

  return (
    <div className="container">
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <h1>🏷️ Gestión de Cupones</h1>
        <div style={{ display: "flex", gap: "0.5rem" }}>
          <button onClick={() => setMostrarForm(!mostrarForm)}>
            {mostrarForm ? "Cancelar" : "+ Nuevo Cupón"}
          </button>
          <button onClick={() => navigate("/admin")} className="btn-secundario">← Panel</button>
        </div>
      </div>

      {mensaje && <p className="exito">{mensaje}</p>}
      {error && <p className="error">{error}</p>}

      {mostrarForm && (
        <div className="card" style={{ marginBottom: "1.5rem", padding: "1.5rem" }}>
          <h3>Nuevo Cupón</h3>
          <form onSubmit={handleCrear} style={{ display: "flex", flexDirection: "column", gap: "0.75rem" }}>
            <input placeholder="Código (ej: VERANO20)" value={form.codigo}
              onChange={(e) => setForm({ ...form, codigo: e.target.value })} required />
            <select value={form.tipo_descuento}
              onChange={(e) => setForm({ ...form, tipo_descuento: e.target.value })}>
              <option value="PORCENTAJE">Porcentaje (%)</option>
              <option value="MONTO_FIJO">Monto fijo ($)</option>
            </select>
            <input type="number" placeholder="Valor del descuento" value={form.valor_descuento}
              onChange={(e) => setForm({ ...form, valor_descuento: e.target.value })} required />
            <input type="datetime-local" value={form.fecha_vencimiento}
              onChange={(e) => setForm({ ...form, fecha_vencimiento: e.target.value })} required />
            <input type="number" placeholder="Límite de usos" value={form.limite_usos}
              onChange={(e) => setForm({ ...form, limite_usos: e.target.value })} required />
            <input type="number" placeholder="ID del evento (requerido)" value={form.evento_id}
              onChange={(e) => setForm({ ...form, evento_id: e.target.value })} required min="1" />
            <p style={{ fontSize: "0.8rem", color: "#666" }}>
              💡 El modelo requiere que cada cupón esté asociado a un evento. Ingresá el ID del evento al que aplica.
            </p>
            <button type="submit">Crear cupón</button>
          </form>
        </div>
      )}

      <table className="tabla-admin">
        <thead>
          <tr>
            <th>Código</th><th>Descuento</th><th>Vencimiento</th>
            <th>Usos</th><th>Estado</th><th>Evento</th><th>Acciones</th>
          </tr>
        </thead>
        <tbody>
          {cupones.map((c) => (
            <tr key={c.ID}>
              <td><strong>{c.Codigo}</strong></td>
              <td>{c.TipoDescuento === "PORCENTAJE" ? `${c.ValorDescuento}%` : `$${c.ValorDescuento}`}</td>
              <td>{new Date(c.FechaVencimiento).toLocaleDateString("es-AR")}</td>
              <td>{c.UsosActuales}/{c.LimiteUsos}</td>
              <td><span className={"badge " + c.Estado?.toLowerCase()}>{c.Estado}</span></td>
              <td>Evento #{c.EventoID}</td>
              <td>
                {c.Estado === "ACTIVO" && (
                  <button onClick={() => handleDesactivar(c.ID)} className="btn-peligro">Desactivar</button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}