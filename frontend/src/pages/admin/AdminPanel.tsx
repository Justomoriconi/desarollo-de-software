import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { api } from "../../api";

interface Evento {
  id: number;
  nombre: string;
  fecha: string;
  lugar: string;
  estado: string;
}

export default function AdminPanel() {
  const [eventos, setEventos] = useState<Evento[]>([]);
  const [error, setError] = useState("");
  const [mensaje, setMensaje] = useState("");
  const navigate = useNavigate();

  const cargarEventos = async () => {
    try {
      const data = await api.get("/eventos");
      setEventos(data || []);
    } catch (err) {
      setError((err as Error).message);
    }
  };

  useEffect(() => { cargarEventos(); }, []);

  const handleCancelar = async (id: number) => {
    if (!confirm("¿Cancelar este evento? Esta acción no se puede deshacer.")) return;
    try {
      setError(""); setMensaje("");
      await api.put(`/admin/eventos/${id}/cancelar`);
      setMensaje("Evento cancelado con éxito");
      cargarEventos();
    } catch (err) {
      setError((err as Error).message);
    }
  };

  return (
    <div className="container">
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <h1>⚙️ Panel de Administración</h1>
        <div style={{ display: "flex", gap: "0.5rem" }}>
          <button onClick={() => navigate("/admin/eventos/nuevo")}>+ Nuevo Evento</button>
          <button onClick={() => navigate("/admin/cupones")} className="btn-secundario">🏷️ Cupones</button>
        </div>
      </div>

      {mensaje && <p className="exito">{mensaje}</p>}
      {error && <p className="error">{error}</p>}

      <table className="tabla-admin">
        <thead>
          <tr>
            <th>ID</th><th>Nombre</th><th>Fecha</th><th>Lugar</th><th>Estado</th><th>Acciones</th>
          </tr>
        </thead>
        <tbody>
          {eventos.map((e) => (
            <tr key={e.id}>
              <td>{e.id}</td>
              <td>{e.nombre}</td>
              <td>{new Date(e.fecha).toLocaleDateString("es-AR")}</td>
              <td>{e.lugar}</td>
              <td><span className={"badge " + e.estado?.toLowerCase()}>{e.estado}</span></td>
              <td>
                <div style={{ display: "flex", gap: "0.4rem" }}>
                  <button onClick={() => navigate(`/admin/eventos/${e.id}/editar`)} className="btn-secundario">Editar</button>
                  <button onClick={() => navigate(`/admin/eventos/${e.id}/reporte`)} className="btn-secundario">Reporte</button>
                  {e.estado === "ACTIVO" && (
                    <button onClick={() => handleCancelar(e.id)} className="btn-peligro">Cancelar</button>
                  )}
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}