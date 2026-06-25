import { BrowserRouter, Routes, Route } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";
import Navbar from "./components/Navbar";
import ProtectedRoute from "./components/ProtectedRoute";
import Home from "./pages/Home";
import EventoDetalle from "./pages/EventoDetalle";
import Login from "./pages/Login";
import Register from "./pages/Register";
import MisEntradas from "./pages/MisEntradas";
import AdminPanel from "./pages/admin/AdminPanel";
import EventoForm from "./pages/admin/EventoForm";
import ReporteEvento from "./pages/admin/ReporteEvento";
import CuponesPanel from "./pages/admin/CuponesPanel";

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Navbar />
        <Routes>
          {/* Rutas públicas */}
          <Route path="/" element={<Home />} />
          <Route path="/eventos/:id" element={<EventoDetalle />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />

          {/* Rutas cliente */}
          <Route path="/mis-entradas" element={
            <ProtectedRoute rol="CLIENTE">
              <MisEntradas />
            </ProtectedRoute>
          } />

          {/* Rutas admin */}
          <Route path="/admin" element={
            <ProtectedRoute rol="ADMIN">
              <AdminPanel />
            </ProtectedRoute>
          } />
          <Route path="/admin/eventos/nuevo" element={
            <ProtectedRoute rol="ADMIN">
              <EventoForm />
            </ProtectedRoute>
          } />
          <Route path="/admin/eventos/:id/editar" element={
            <ProtectedRoute rol="ADMIN">
              <EventoForm />
            </ProtectedRoute>
          } />
          <Route path="/admin/eventos/:id/reporte" element={
            <ProtectedRoute rol="ADMIN">
              <ReporteEvento />
            </ProtectedRoute>
          } />
          <Route path="/admin/cupones" element={
            <ProtectedRoute rol="ADMIN">
              <CuponesPanel />
            </ProtectedRoute>
          } />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}