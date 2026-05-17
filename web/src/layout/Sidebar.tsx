import { NavLink } from "react-router-dom";

const items = [
  { to: "/", label: "Dashboard", icon: "▦" },
  { to: "/operations", label: "Operations", icon: "›_" },
  { to: "/keys", label: "Keys", icon: "⚿" },
  { to: "/metrics", label: "Metrics", icon: "▥" },
  { to: "/audit", label: "Audit Logs", icon: "▤" },
];

export function Sidebar() {
  return (
    <aside className="sidebar">
      <div className="brand">
        <div className="brand-icon">⚿</div>
        <span>KMIP Lab</span>
      </div>

      <nav className="nav">
        {items.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) => `nav-item ${isActive ? "active" : ""}`}
          >
            <span>{item.icon}</span>
            {item.label}
          </NavLink>
        ))}
      </nav>

      <div className="api-status">
        <span className="status-dot"></span>
        API Online
      </div>
    </aside>
  );
}