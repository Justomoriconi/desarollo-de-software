import { BrowserRouter, Routes, Route } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";
import Navbar from "./components/Navbar";
import Home from "./pages/Home";
import EventoDetalle from "./pages/EventoDetalle";
import Login from "./pages/Login";
import Register from "./pages/Register";
import MisEntradas from "./pages/MisEntradas";

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Navbar />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/eventos/:id" element={<EventoDetalle />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/mis-entradas" element={<MisEntradas />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}