import { createContext, useContext, useState } from "react";
import type { ReactNode } from "react";
import { api } from "../api";

interface AuthContextType {
  token: string | null;
  isLoggedIn: boolean;
  userRol: string | null;
  login: (email: string, password: string) => Promise<void>;
  register: (nombre: string, email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | null>(null);

function getRolFromToken(token: string): string | null {
  try {
    const payload = JSON.parse(atob(token.split(".")[1]));
    return payload.rol || null;
  } catch {
    return null;
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(
    localStorage.getItem("token")
  );
  const [userRol, setUserRol] = useState<string | null>(() => {
    const t = localStorage.getItem("token");
    return t ? getRolFromToken(t) : null;
  });

  const login = async (email: string, password: string) => {
    const data = await api.post("/login", { email, password });
    localStorage.setItem("token", data.token);
    setToken(data.token);
    setUserRol(getRolFromToken(data.token));
  };

  const register = async (nombre: string, email: string, password: string) => {
    await api.post("/register", { nombre, email, password });
    await login(email, password);
  };

  const logout = () => {
    localStorage.removeItem("token");
    setToken(null);
    setUserRol(null);
  };

  return (
    <AuthContext.Provider
      value={{ token, isLoggedIn: !!token, userRol, login, register, logout }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth debe usarse dentro de AuthProvider");
  }
  return context;
}