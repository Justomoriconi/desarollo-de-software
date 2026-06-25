import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { api } from "../../api";

interface Comprador {
  usuario_id: number;
  nombre: string;
  email: string;
  tipo_entrada: string;
  estado: string;
}

interface TipoReporte {
  id: number;
  nombre: string;
  precio: number;
  stock_disponible: number;
  vendidas: number;
}

interface Reporte {
  evento_id: number;
  nombre: string;
  estado: string;
  total_tipos_entrada: number;
  total_vendidas: number;
  total_canceladas: number;
  total_activas: number;
  compradores: Comprador[];
  tipos_entrada: TipoReporte[];
}

export default function ReporteEvento() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [reporte, setReporte] = useState<Reporte | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    api.get(`/admin/eventos/${id}/reporte`)
      .then(setReporte)
      .catch((err) => setError((err as Error).message));
  }, [id]);

  if (error) return <div className="container"><p className="error">{error}</p></div>;
  if (!reporte) return <div className="container"><p>Cargando reporte...</p></div>;

  return (
    <div className="container">
      <button onClick={() => navigate("/admin")} className="btn-secundario" style={{ marginBottom: "1rem" }}>
        ← Volver al panel
      </button>
      <h1>Reporte: {reporte.nombre}</h1>
      <span className={"badge " + reporte.estado?.toLowerCase()}>{reporte.estado}</span>

      <div className="reporte-metricas">
        <div className="metrica">
          <span className="metrica-valor">{reporte.total_vendidas}</span>
          <span className="metrica-label">Total vendidas</span>
        </div>
        <div className="metrica">
          <span className="metrica-valor" style={{ color: "#1e8449" }}>{reporte.total_activas}</span>
          <span className="metrica-label">Activas</span>
        </div>
        <div className="metrica">
          <span className="metrica-valor" style={{ color: "#c0392b" }}>{reporte.total_canceladas}</span>
          <span className="metrica-label">Canceladas</span>
        </div>
      </div>

      <h2>Por tipo de entrada</h2>
      <table className="tabla-admin">
        <thead>
          <tr><th>Tipo</th><th>Precio</th><th>Vendidas</th><th>Stock disponible</th></tr>
        </thead>
        <tbody>
          {reporte.tipos_entrada?.map((t) => (
            <tr key={t.id}>
              <td>{t.nombre}</td>
              <td>${t.precio}</td>
              <td>{t.vendidas}</td>
              <td>{t.stock_disponible}</td>
            </tr>
          ))}
        </tbody>
      </table>

      <h2>Compradores</h2>
      {!reporte.compradores || reporte.compradores.length === 0 ? (
        <p>No hay compradores para este evento.</p>
      ) : (
        <table className="tabla-admin">
          <thead>
            <tr><th>Nombre</th><th>Email</th><th>Tipo entrada</th><th>Estado</th></tr>
          </thead>
          <tbody>
            {reporte.compradores.map((c, i) => (
              <tr key={i}>
                <td>{c.nombre}</td>
                <td>{c.email}</td>
                <td>{c.tipo_entrada}</td>
                <td><span className={"badge " + c.estado?.toLowerCase()}>{c.estado}</span></td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}