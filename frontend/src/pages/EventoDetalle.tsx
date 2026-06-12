import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { api } from "../api";
import { useAuth } from "../context/AuthContext";

interface Evento {
  id: number;
  nombre: string;
  descripcion: string;
  fecha: string;
  lugar: string;
  estado: string;
}

interface TipoEntrada {
  id: number;
  nombre: string;
  precio: number;
  stock_disponible: number;
}

export default function EventoDetalle() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { isLoggedIn } = useAuth();

  const [evento, setEvento] = useState<Evento | null>(null);
  const [tiposEntrada, setTiposEntrada] = useState<TipoEntrada[]>([]);
  const [error, setError] = useState("");
  const [exito, setExito] = useState("");
  const [comprando, setComprando] = useState(false);

  useEffect(() => {
    const cargar = async () => {
      try {
        const data = await api.get("/eventos/" + id);
        setEvento(data);

        const tipos = await api.get("/eventos/" + id + "/tipos-entrada");
        setTiposEntrada(tipos || []);
      } catch (err) {
        setError((err as Error).message);
      }
    };
    cargar();
  }, [id]);

  const handleComprar = async (tipoEntradaId: number) => {
    if (!isLoggedIn) {
      navigate("/login");
      return;
    }

    try {
      setComprando(true);
      setError("");
      setExito("");
      const data = await api.post("/tickets", {
        tipo_entrada_id: tipoEntradaId,
      });
      setExito("🎉 ¡Felicitaciones! Compra realizada con éxito. Ticket #" + data.ticket_id);
    } catch (err) {
      setError("❌ " + (err as Error).message);
    } finally {
      setComprando(false);
    }
  };

  if (!evento) {
    return (
      <div className="container">
        {error ? <p className="error">{error}</p> : <p>Cargando...</p>}
      </div>
    );
  }

  return (
    <div className="container">
      <h1>{evento.nombre}</h1>
      <p className="descripcion">{evento.descripcion}</p>
      <p>📅 {new Date(evento.fecha).toLocaleString("es-AR")}</p>
      <p>📍 {evento.lugar}</p>
      <span className={"badge " + evento.estado.toLowerCase()}>
        {evento.estado}
      </span>

      <h2>Entradas</h2>
      {exito && <p className="exito">{exito}</p>}
      {error && <p className="error">{error}</p>}

      {tiposEntrada.length === 0 && (
        <p>No hay tipos de entrada disponibles para este evento.</p>
      )}

      <div className="tipos-entrada">
        {tiposEntrada.map((tipo) => (
          <div key={tipo.id} className="tipo-entrada">
            <div>
              <strong>{tipo.nombre}</strong>
              <p>${tipo.precio} — Stock: {tipo.stock_disponible}</p>
            </div>
            <button
              onClick={() => handleComprar(tipo.id)}
              disabled={comprando || tipo.stock_disponible <= 0 || evento.estado !== "ACTIVO"}
            >
              {tipo.stock_disponible <= 0 ? "Sin stock" : "Comprar"}
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}