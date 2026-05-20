import { Outlet } from "react-router-dom";
import { useState } from "react";
import { Sidebar } from "./Sidebar";
import { Topbar } from "./Topbar";

export function DashboardLayout() {
  const [menuOpen, setMenuOpen] = useState(false);

  return (
    <div className={`app-shell ${menuOpen ? "menu-open" : ""}`}>
      <button
        className="sidebar-backdrop"
        aria-label="Close navigation menu"
        onClick={() => setMenuOpen(false)}
      />

      <Sidebar onNavigate={() => setMenuOpen(false)} />

      <main className="main">
        <Topbar onMenuClick={() => setMenuOpen(true)} />
        <div className="page">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
