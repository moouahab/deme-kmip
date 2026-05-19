import { useEffect, useState } from "react";
import { getAuditEvents } from "../services/api";
import type { AuditEvent } from "../services/api";
import { StatCard } from "../layout/StatCard";

export function AuditPage() {
  const [events, setEvents] = useState<AuditEvent[]>([]);
  const [error, setError] = useState("");

  async function loadAudit() {
    try {
      setError("");
      const data = await getAuditEvents();
      setEvents(data);
    } catch {
      setError("GET /audit is not available yet. Add the backend route first.");
    }
  }

  useEffect(() => {
    loadAudit();
  }, []);

  const success = events.filter((event) => event.result === "success").length;
  const notFound = events.filter((event) => event.result === "not_found").length;
  const errors = events.filter((event) => event.result === "error").length;
  const recentEvents = [...events].reverse();

  return (
    <>
      <section className="page-header with-actions">
        <div>
          <h2>Audit Logs</h2>
          <p>Complete history of all KMIP operations</p>
        </div>

        <div className="button-row clean">
          <button onClick={loadAudit}>↻ Refresh</button>
          <button className="primary">↓ Export JSON</button>
        </div>
      </section>

      <section className="stats-grid three">
        <StatCard label="Successful Operations" value={success} tone="green" />
        <StatCard label="Not Found" value={notFound} tone="orange" />
        <StatCard label="Errors" value={errors} tone="red" />
      </section>

      <section className="panel">
        <h3>Recent Audit Events</h3>
        <p>Operation history and results</p>

        {error && <div className="warning-box">{error}</div>}

        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Time</th>
                <th>Operation</th>
                <th>Key ID</th>
                <th>Status</th>
                <th>Result</th>
                <th>Error</th>
              </tr>
            </thead>

            <tbody>
              {events.length === 0 ? (
                <tr>
                  <td colSpan={6}>No audit events available.</td>
                </tr>
              ) : (
                recentEvents.map((event, index) => (
                  <tr key={index}>
                    <td>{new Date(event.time).toLocaleTimeString()}</td>
                    <td>
                      <span className="op-badge">{event.operation}</span>
                    </td>
                    <td>{event.key_id || "-"}</td>
                    <td>{event.status || "-"}</td>
                    <td>
                      <span className={`result-badge ${event.result}`}>
                        {event.result}
                      </span>
                    </td>
                    <td>{event.error || "-"}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </section>
    </>
  );
}
