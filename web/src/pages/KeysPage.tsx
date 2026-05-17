import { useEffect, useMemo, useState } from "react";
import { getKeys} from "../services/api";
import type { KMSKey } from "../services/api";

function objectTypeName(value: number) {
  switch (value) {
    case 2:
      return "Symmetric Key";
    case 3:
      return "Public Key";
    case 4:
      return "Private Key";
    case 7:
      return "Secret Data";
    default:
      return `Object Type ${value}`;
  }
}

function formatDate(value: string) {
  if (!value) return "-";
  return new Date(value).toLocaleString();
}

export function KeysPage() {
  const [keys, setKeys] = useState<KMSKey[]>([]);
  const [query, setQuery] = useState("");
  const [state, setState] = useState("all");
  const [error, setError] = useState("");

  async function loadKeys() {
    try {
      setError("");
      const data = await getKeys();
      setKeys(data);
    } catch {
      setError("GET /keys is not available yet. Add the backend route first.");
    }
  }

  useEffect(() => {
    loadKeys();
  }, []);

  const filteredKeys = useMemo(() => {
    return keys.filter((key) => {
      const matchesQuery = key.id.toLowerCase().includes(query.toLowerCase());
      const matchesState = state === "all" || key.status === state;

      return matchesQuery && matchesState;
    });
  }, [keys, query, state]);

  const activeCount = keys.filter((key) => key.status === "active").length;
  const destroyedCount = keys.filter((key) => key.status === "destroyed").length;

  return (
    <>
      <section className="page-header">
        <h2>Keys</h2>
        <p>View and manage all cryptographic keys</p>
      </section>

      <section className="panel">
        <div className="toolbar">
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search by key ID..."
          />

          <select value={state} onChange={(e) => setState(e.target.value)}>
            <option value="all">All States</option>
            <option value="active">Active</option>
            <option value="destroyed">Destroyed</option>
            <option value="revoked">Revoked</option>
            <option value="pre_active">Pre-Active</option>
          </select>

          <button onClick={loadKeys}>↻ Refresh</button>
        </div>

        {error && <div className="warning-box">{error}</div>}

        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Key ID</th>
                <th>Object Type</th>
                <th>State</th>
                <th>Created</th>
                <th>Updated</th>
              </tr>
            </thead>

            <tbody>
              {filteredKeys.length === 0 ? (
                <tr>
                  <td colSpan={5}>No keys found.</td>
                </tr>
              ) : (
                filteredKeys.map((key) => (
                  <tr key={key.id}>
                    <td>
                      <span className="key-cell">⚿ {key.id}</span>
                    </td>
                    <td>{objectTypeName(key.object_type)}</td>
                    <td>
                      <span className={`state-badge ${key.status}`}>{key.status}</span>
                    </td>
                    <td>{formatDate(key.created_at)}</td>
                    <td>{formatDate(key.updated_at)}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        <div className="table-footer">
          <span>Showing {filteredKeys.length} of {keys.length} keys</span>
          <span>
            <span className="green-dot"></span> {activeCount} Active ·{" "}
            <span className="red-dot"></span> {destroyedCount} Destroyed
          </span>
        </div>
      </section>
    </>
  );
}