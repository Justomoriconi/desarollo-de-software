import { useEffect, useState } from "react";
import { api } from "../api";

interface Ticket {
  id: number;
  evento_nombre: string;
  tipo_entrada: string;
  precio_pagado: number;
  fecha_compra: string;
  estado: string;
}

export default function MisEntradas() {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [error, setError] = useState("");
  const [mensaje, setMensaje] = useState("");
  const [transferiendo, setTransferiendo] = useState<number | null>(null);
  const [emailDestino, setEmailDestino] = useState("");

  const cargarTickets = async () => {
    try {
      setError("");
      const data = await api.get("/tickets");
      setTickets(data || []);
    } catch (err) {
      setError((err as Error).message);
    }
  };

  useEffect(() => {
    cargarTickets();
  }, []);

  const handleCancelar = async (id: number) => {
    if (!confirm("¿Seguro que querés cancelar esta entrada?")) return;

    try {
      setError("");
      setMensaje("");
      await api.put("/tickets/" + id + "/cancelar");
      setMensaje("Entrada cancelada con éxito");
      cargarTickets();
    } catch (err) {
      setError((err as Error).message);
    }
  };

  const handleTransferir = async (id: number) => {
    try {
      setError("");
      setMensaje("");
      await api.put("/tickets/" + id + "/transferir", {
        email_destino: emailDestino,
      });
      setMensaje("Entrada transferida con éxito a " + emailDestino);
      setTransferiendo(null);
      setEmailDestino("");
      cargarTickets();
    } catch (err) {
      setError((err as Error).message);
    }
  };

  return (
    <div className="container">
      <h1>Mis Entradas</h1>

      {mensaje && <p className="exito">{mensaje}</p>}
      {error && <p className="error">{error}</p>}

      {tickets.length === 0 && <p>Todavía no compraste ninguna entrada.</p>}

      <div className="lista-tickets">
        {tickets.map((ticket) => (
          <div key={ticket.id} className="ticket">
            <div>
              <strong>{ticket.evento_nombre}</strong>
              <p>{ticket.tipo_entrada} — ${ticket.precio_pagado}</p>
              <p>
                Comprada el {new Date(ticket.fecha_compra).toLocaleString("es-AR")}
              </p>
              <span className={"badge " + ticket.estado.toLowerCase()}>
                {ticket.estado}
              </span>
            </div>

            {ticket.estado === "ACTIVO" && (
              <div className="acciones">
                <button onClick={() => handleCancelar(ticket.id)} className="btn-peligro">
                  Cancelar
                </button>
                <button onClick={() => setTransferiendo(ticket.id)}>
                  Transferir
                </button>
              </div>
            )}

            {transferiendo === ticket.id && (
              <div className="transferir">
                <input
                  type="email"
                  placeholder="Email del destinatario"
                  value={emailDestino}
                  onChange={(e) => setEmailDestino(e.target.value)}
                />
                <button onClick={() => handleTransferir(ticket.id)}>
                  Confirmar
                </button>
                <button onClick={() => setTransferiendo(null)} className="btn-secundario">
                  Cerrar
                </button>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}