import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function Navbar() {
  const { isLoggedIn, userRol, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/");
  };

  return (
    <nav className="navbar">
      <Link to="/" className="navbar-brand">🎟️ Ticketera</Link>
      <div className="navbar-links">
        <Link to="/">Eventos</Link>
        {isLoggedIn ? (
          <>
            {userRol === "ADMIN" && (
              <Link to="/admin" className="link-admin">⚙️ Admin</Link>
            )}
            <Link to="/mis-entradas">Mis Entradas</Link>
            <button onClick={handleLogout} className="btn-link">
              Cerrar sesión
            </button>
          </>
        ) : (
          <>
            <Link to="/login">Iniciar sesión</Link>
            <Link to="/register">Registrarse</Link>
          </>
        )}
      </div>
    </nav>
  );
}