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

interface ValidacionCupon {
  cupon_id: number;
  codigo: string;
  tipo_descuento: string;
  valor_descuento: number;
  precio_original: number;
  precio_final: number;
  descuento: number;
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

  // Estado del cupón
  const [codigoCupon, setCodigoCupon] = useState<Record<number, string>>({});
  const [validacionCupon, setValidacionCupon] = useState<Record<number, ValidacionCupon | null>>({});
  const [errorCupon, setErrorCupon] = useState<Record<number, string>>({});
  const [validando, setValidando] = useState<Record<number, boolean>>({});

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

  const handleValidarCupon = async (tipoId: number) => {
    const codigo = codigoCupon[tipoId];
    if (!codigo) return;

    try {
      setValidando({ ...validando, [tipoId]: true });
      setErrorCupon({ ...errorCupon, [tipoId]: "" });
      const data = await api.post("/cupones/validar", {
        codigo,
        tipo_entrada_id: tipoId,
      });
      setValidacionCupon({ ...validacionCupon, [tipoId]: data });
    } catch (err) {
      setValidacionCupon({ ...validacionCupon, [tipoId]: null });
      setErrorCupon({ ...errorCupon, [tipoId]: (err as Error).message });
    } finally {
      setValidando({ ...validando, [tipoId]: false });
    }
  };

  const handleComprar = async (tipoId: number) => {
    if (!isLoggedIn) {
      navigate("/login");
      return;
    }

    try {
      setComprando(true);
      setError("");
      setExito("");

      const body: { tipo_entrada_id: number; codigo_cupon?: string } = {
        tipo_entrada_id: tipoId,
      };

      const cupon = validacionCupon[tipoId];
      if (cupon) {
        body.codigo_cupon = cupon.codigo;
      }

      const data = await api.post("/tickets", body);
      const precioFinal = data.ticket?.PrecioPagado || data.ticket?.precio_pagado;
      setExito(
        `🎉 ¡Felicitaciones! Compra realizada con éxito. Ticket #${data.ticket?.ID || data.ticket?.id}` +
        (cupon ? ` — Pagaste $${precioFinal} (descuento aplicado: $${cupon.descuento})` : "")
      );

      // Limpiar cupón usado
      setValidacionCupon({ ...validacionCupon, [tipoId]: null });
      setCodigoCupon({ ...codigoCupon, [tipoId]: "" });

      // Refrescar stock
      const tipos = await api.get("/eventos/" + id + "/tipos-entrada");
      setTiposEntrada(tipos || []);
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

      <h2>Entradas disponibles</h2>
      {exito && <p className="exito">{exito}</p>}
      {error && <p className="error">{error}</p>}

      {tiposEntrada.length === 0 && (
        <p>No hay tipos de entrada disponibles para este evento.</p>
      )}

      <div className="tipos-entrada">
        {tiposEntrada.map((tipo) => {
          const validacion = validacionCupon[tipo.id];
          return (
            <div key={tipo.id} className="tipo-entrada">
              <div>
                <strong>{tipo.nombre}</strong>
                <p>
                  {validacion ? (
                    <>
                      <span style={{ textDecoration: "line-through", color: "#999" }}>
                        ${tipo.precio}
                      </span>
                      {" → "}
                      <span style={{ color: "#1e8449", fontWeight: "bold" }}>
                        ${validacion.precio_final}
                      </span>
                      {" "}
                      <span className="badge activo">-${validacion.descuento}</span>
                    </>
                  ) : (
                    `$${tipo.precio}`
                  )}
                  {" — Stock: "}{tipo.stock_disponible}
                </p>

                {/* Campo de cupón */}
                {isLoggedIn && evento.estado === "ACTIVO" && tipo.stock_disponible > 0 && (
                  <div className="cupon-row">
                    <input
                      type="text"
                      placeholder="Código de cupón (opcional)"
                      value={codigoCupon[tipo.id] || ""}
                      onChange={(e) =>
                        setCodigoCupon({ ...codigoCupon, [tipo.id]: e.target.value })
                      }
                    />
                    <button
                      onClick={() => handleValidarCupon(tipo.id)}
                      disabled={!codigoCupon[tipo.id] || validando[tipo.id]}
                      className="btn-secundario"
                    >
                      Aplicar
                    </button>
                    {errorCupon[tipo.id] && (
                      <span className="error-inline">{errorCupon[tipo.id]}</span>
                    )}
                    {validacion && (
                      <span className="exito-inline">✓ Cupón aplicado</span>
                    )}
                  </div>
                )}
              </div>

              <button
                onClick={() => handleComprar(tipo.id)}
                disabled={
                  comprando ||
                  tipo.stock_disponible <= 0 ||
                  evento.estado !== "ACTIVO"
                }
              >
                {tipo.stock_disponible <= 0 ? "Sin stock" : "Comprar"}
              </button>
            </div>
          );
        })}
      </div>
    </div>
  );
}