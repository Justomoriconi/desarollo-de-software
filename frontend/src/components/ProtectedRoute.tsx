import { Navigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

interface Props {
  children: React.ReactNode;
  rol: "CLIENTE" | "ADMIN";
}

export default function ProtectedRoute({ children, rol }: Props) {
  const { isLoggedIn, userRol } = useAuth();

  if (!isLoggedIn) {
    return <Navigate to="/login" replace />;
  }

  if (rol === "ADMIN" && userRol !== "ADMIN") {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
}