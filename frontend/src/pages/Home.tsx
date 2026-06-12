import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { api } from "../api";

interface Evento {
  id: number;
  nombre: string;
  descripcion: string;
  fecha: string;
  lugar: string;
  estado: string;
}

export default function Home() {
  const [eventos, setEventos] = useState<Evento[]>([]);
  const [busqueda, setBusqueda] = useState("");
  const [estado, setEstado] = useState("");
  const [error, setError] = useState("");

  const cargarEventos = async () => {
    try {
      setError("");
      const params = new URLSearchParams();
      if (busqueda) params.append("nombre", busqueda);
      if (estado) params.append("estado", estado);

      const query = params.toString() ? "?" + params.toString() : "";
      const data = await api.get("/eventos" + query);
      setEventos(data || []);
    } catch (err) {
      setError((err as Error).message);
    }
  };

  useEffect(() => {
    cargarEventos();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleBuscar = (e: React.FormEvent) => {
    e.preventDefault();
    cargarEventos();
  };

  return (
    <div className="container">
      <h1>Catálogo de Eventos</h1>

      <form onSubmit={handleBuscar} className="filtros">
        <input
          type="text"
          placeholder="Buscar por nombre..."
          value={busqueda}
          onChange={(e) => setBusqueda(e.target.value)}
        />
        <select value={estado} onChange={(e) => setEstado(e.target.value)}>
          <option value="">Todos los estados</option>
          <option value="ACTIVO">Activos</option>
          <option value="CANCELADO">Cancelados</option>
        </select>
        <button type="submit">Buscar</button>
      </form>

      {error && <p className="error">{error}</p>}

      <div className="grilla">
        {eventos.length === 0 && !error && <p>No se encontraron eventos.</p>}
        {eventos.map((evento) => (
          <Link to={"/eventos/" + evento.id} key={evento.id} className="card">
            <h3>{evento.nombre}</h3>
            <p>📅 {new Date(evento.fecha).toLocaleString("es-AR")}</p>
            <p>📍 {evento.lugar}</p>
            <span className={"badge " + evento.estado.toLowerCase()}>
              {evento.estado}
            </span>
          </Link>
        ))}
      </div>
    </div>
  );
}