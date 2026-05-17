import { useMemo, useState } from "react";
import {
  buildCreateKeyRequest,
  buildDestroyKeyRequest,
  buildGetKeyRequest,
  toHex,
} from "../services/kmip";

import {sendKMIPRequest } from "../services/api";
import type { KMIPResponse } from "../services/api";

type Operation = "create" | "get" | "destroy";

export function OperationsPage() {
  const [operation, setOperation] = useState<Operation>("create");
  const [keyId, setKeyId] = useState("");
  const [response, setResponse] = useState<KMIPResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [showRawBytes, setShowRawBytes] = useState(false);

  const requestBytes = useMemo(() => {
    if (operation === "create") return buildCreateKeyRequest();
    if (!keyId.trim()) return new Uint8Array();
    if (operation === "get") return buildGetKeyRequest(keyId.trim());
    return buildDestroyKeyRequest(keyId.trim());
  }, [operation, keyId]);

  async function sendOperation() {
    if (operation !== "create" && !keyId.trim()) {
      setResponse({
        status: 0,
        ok: false,
        data: {
          error: "missing_key_id",
          message: "Key ID is required for Get and Destroy operations",
        },
      });
      return;
    }

    setLoading(true);

    try {
      const result = await sendKMIPRequest(requestBytes);
      setResponse(result);

      if (
        operation === "create" &&
        result.ok &&
        typeof result.data === "object" &&
        result.data !== null &&
        "key_id" in result.data
      ) {
        setKeyId(String(result.data.key_id));
      }
    } finally {
      setLoading(false);
    }
  }

  function resetForm() {
    setOperation("create");
    setKeyId("");
    setResponse(null);
    setShowRawBytes(false);
  }

  return (
    <>
      <section className="page-header">
        <h2>Operations</h2>
        <p>Execute KMIP operations to manage cryptographic keys</p>
      </section>

      <section className="operations-grid">
        <div className="panel">
          <h3>Operation Builder</h3>
          <p>Generate and send TTLV encoded requests to the KMIP endpoint</p>

          <div className="form-stack">
            <label>
              Operation
              <select value={operation} onChange={(e) => setOperation(e.target.value as Operation)}>
                <option value="create">Create Key</option>
                <option value="get">Get Key</option>
                <option value="destroy">Destroy Key</option>
              </select>
            </label>

            <label>
              Key ID
              <input
                value={keyId}
                onChange={(e) => setKeyId(e.target.value)}
                placeholder="key-xxxxxxxx"
                disabled={operation === "create"}
              />
              {operation === "create" && <small>Auto-generated on key creation</small>}
            </label>

            <label>
              Object Type
              <select disabled={operation !== "create"}>
                <option>Symmetric Key</option>
              </select>
            </label>
          </div>

          <div className="button-row">
            <button className="primary wide" onClick={sendOperation} disabled={loading}>
              {loading ? "Sending..." : "➤ Send Operation"}
            </button>
            <button onClick={resetForm}>↻ Reset</button>
          </div>

          <div className="endpoint-line">POST /kmip · application/octet-stream</div>
        </div>

        <div className="panel">
          <h3>KMIP Response</h3>

          {!response ? (
            <div className="empty-state">
              <div>ⓘ</div>
              <p>No response yet. Send an operation to see results.</p>
            </div>
          ) : (
            <pre className="code-box">{JSON.stringify(response, null, 2)}</pre>
          )}
        </div>

        <div className="panel ttlv-panel">
          <h3>TTLV Inspector</h3>
          <p>Decode the binary request structure</p>

          <div className="ttlv-row">
            <div className="tag-icon">{"</>"}</div>
            <div>
              <strong>TagOperation</strong>
              <span>0x42005C</span>
            </div>
            <small>Enumeration · {operation}_key</small>
          </div>

          {operation === "create" ? (
            <div className="ttlv-row">
              <div className="tag-icon">{"</>"}</div>
              <div>
                <strong>ObjectType</strong>
                <span>0x420057</span>
              </div>
              <small>Enumeration · symmetric_key</small>
            </div>
          ) : (
            <div className="ttlv-row">
              <div className="tag-icon">{"</>"}</div>
              <div>
                <strong>UniqueIdentifier</strong>
                <span>0x420094</span>
              </div>
              <small>TextString · key_id</small>
            </div>
          )}

          <button className="wide" onClick={() => setShowRawBytes((value) => !value)}>
            Show Raw Bytes
          </button>

          {showRawBytes && <pre className="code-box small">{toHex(requestBytes) || "No bytes"}</pre>}
        </div>
      </section>
    </>
  );
}