import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { emptyMetrics, getKeys, getMetrics} from "../services/api";
import type { KMSKey, Metrics } from "../services/api";
import { StatCard } from "../layout/StatCard";

export function DashboardPage() {
  const [metrics, setMetrics] = useState<Metrics>(emptyMetrics);
  const [keys, setKeys] = useState<KMSKey[]>([]);

  async function loadData() {
    try {
      const [metricsData, keysData] = await Promise.all([getMetrics(), getKeys()]);
      setMetrics(metricsData);
      setKeys(keysData);
    } catch {
      // backend route may not exist yet
    }
  }

  useEffect(() => {
    loadData();
  }, []);

  const activeKeys = keys.filter((key) => key.status === "active").length;

  return (
    <>
      <section className="page-header">
        <h2>Welcome to KMIP Lab</h2>
        <p>Browser-based KMIP testing console for cryptographic key management operations</p>
      </section>

      <section className="stats-grid four">
        <StatCard label="Active Keys" value={activeKeys} icon="⚿" tone="blue" />
        <StatCard label="Total Requests" value={metrics.http_requests_total} icon="⌁" tone="green" />
        <StatCard label="Total Keys" value={keys.length} icon="⬡" tone="purple" />
        <StatCard label="Successful Ops" value={metrics.success_total} icon="↗" tone="blue" />
      </section>

      <section className="feature-grid">
        <Link to="/operations" className="feature-card blue">
          <div className="feature-icon">⚿</div>
          <h3>Operations</h3>
          <p>Create, retrieve, and destroy cryptographic keys using KMIP operations.</p>
        </Link>

        <Link to="/keys" className="feature-card green">
          <div className="feature-icon">⬡</div>
          <h3>Keys</h3>
          <p>View all created keys, their status, and lifecycle information.</p>
        </Link>

        <Link to="/metrics" className="feature-card purple">
          <div className="feature-icon">↗</div>
          <h3>Metrics</h3>
          <p>Monitor operation statistics and server metrics.</p>
        </Link>

        <Link to="/audit" className="feature-card orange">
          <div className="feature-icon">⌁</div>
          <h3>Audit Logs</h3>
          <p>Review operation history when the backend audit route is enabled.</p>
        </Link>
      </section>

      <section className="info-card">
        <div className="feature-icon">⬡</div>
        <div>
          <h3>About KMIP Lab</h3>
          <p>
            KMIP Lab is a browser-based testing console for KMIP-like operations.
            It sends TTLV encoded binary requests to your Go backend and displays
            the response, keys, and metrics in a clean interface.
          </p>
        </div>
      </section>
    </>
  );
}