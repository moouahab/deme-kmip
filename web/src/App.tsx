import { Navigate, Route, Routes } from "react-router-dom";
import { DashboardLayout } from "./layout/DashboardLayout";
import { DashboardPage } from "./pages/DashboardPage";
import { OperationsPage } from "./pages/OperationsPage";
import { KeysPage } from "./pages/KeysPage";
import { MetricsPage } from "./pages/MetricsPage";
import { AuditPage } from "./pages/AuditPage";

export default function App() {
  return (
    <Routes>
      <Route element={<DashboardLayout />}>
        <Route path="/" element={<DashboardPage />} />
        <Route path="/operations" element={<OperationsPage />} />
        <Route path="/keys" element={<KeysPage />} />
        <Route path="/metrics" element={<MetricsPage />} />
        <Route path="/audit" element={<AuditPage />} />
      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}