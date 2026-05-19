import { useMemo, useState } from "react";
import {
  buildActivateKeyRequest,
  buildCreateKeyRequest,
  buildDestroyKeyRequest,
  buildGetAttributesRequest,
  buildGetKeyRequest,
  buildLocateKeysRequest,
  buildRevokeKeyRequest,
  toHex,
} from "../services/kmip";

import { sendKMIPRequest } from "../services/api";
import type { KMIPResponse } from "../services/api";

type Operation =
  | "create"
  | "get"
  | "destroy"
  | "activate"
  | "revoke"
  | "locate"
  | "get_attributes";

export function OperationsPage() {
  const [operation, setOperation] = useState<Operation>("create");
  const [keyId, setKeyId] = useState("");
  const [response, setResponse] = useState<KMIPResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [showRawBytes, setShowRawBytes] = useState(false);

  const requestBytes = useMemo(() => {
    if (operation === "create") return buildCreateKeyRequest();
    if (operation === "locate") return buildLocateKeysRequest();
    if (!keyId.trim()) return new Uint8Array();
    if (operation === "get") return buildGetKeyRequest(keyId.trim());
    if (operation === "activate") return buildActivateKeyRequest(keyId.trim());
    if (operation === "revoke") return buildRevokeKeyRequest(keyId.trim());
    if (operation === "get_attributes") return buildGetAttributesRequest(keyId.trim());
    return buildDestroyKeyRequest(keyId.trim());
  }, [operation, keyId]);

  async function sendOperation() {
    if (operation !== "create" && operation !== "locate" && !keyId.trim()) {
      setResponse({
        status: 0,
        ok: false,
        data: {
          error: "missing_key_id",
          message: "Key ID is required for this operation",
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
                <option value="activate">Activate Key</option>
                <option value="revoke">Revoke Key</option>
                <option value="destroy">Destroy Key</option>
                <option value="locate">Locate Keys</option>
                <option value="get_attributes">Get Attributes</option>
              </select>
            </label>

            <label>
              Key ID
              <input
                value={keyId}
                onChange={(e) => setKeyId(e.target.value)}
                placeholder="key-xxxxxxxx"
                disabled={operation === "create" || operation === "locate"}
              />
              {operation === "create" && <small>Auto-generated on key creation</small>}
              {operation === "locate" && <small>No key ID required for Locate</small>}
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

        <div className="operations-stack">
          <div className="panel response-panel">
            <div className="response-header">
              <div>
                <h3>KMIP Response</h3>
                <p>Decoded HTTP response from the KMIP endpoint</p>
              </div>

              {response && (
                <span className={`response-status ${response.ok ? "ok" : "error"}`}>
                  {response.status || "local"} · {response.ok ? "OK" : "Error"}
                </span>
              )}
            </div>

            {!response ? (
              <div className="empty-state">
                <div>ⓘ</div>
                <p>No response yet. Send an operation to see results.</p>
              </div>
            ) : (
              <SyntaxHighlightedJSON value={response} />
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
            ) : operation !== "locate" ? (
              <div className="ttlv-row">
                <div className="tag-icon">{"</>"}</div>
                <div>
                  <strong>UniqueIdentifier</strong>
                  <span>0x420094</span>
                </div>
                <small>TextString · key_id</small>
              </div>
            ) : null}

            <button className="wide" onClick={() => setShowRawBytes((value) => !value)}>
              Show Raw Bytes
            </button>

            {showRawBytes && <pre className="code-box small">{toHex(requestBytes) || "No bytes"}</pre>}
          </div>
        </div>
      </section>
    </>
  );
}

function SyntaxHighlightedJSON({ value }: { value: unknown }) {
  const lines = JSON.stringify(value, null, "\t").split("\n");

  return (
    <pre className="code-box json-view" aria-label="KMIP response JSON">
      {lines.map((line, index) => (
        <span className="json-line" key={`${index}-${line}`}>
          <span className="json-line-number">{index + 1}</span>
          <span className="json-line-content">{highlightJSONLine(line)}</span>
        </span>
      ))}
    </pre>
  );
}

function highlightJSONLine(line: string) {
  const parts = line.split(/("(?:\\.|[^"\\])*"(?=\s*:)|"(?:\\.|[^"\\])*"|true|false|null|-?\d+(?:\.\d+)?)/g);

  return parts.map((part, index) => {
    if (!part) return null;

    let className = "json-punctuation";
    if (/^"(?:\\.|[^"\\])*"$/.test(part)) {
      className = part.endsWith('"') && parts[index + 1]?.startsWith(":") ? "json-key" : "json-string";
    } else if (/^(true|false)$/.test(part)) {
      className = "json-boolean";
    } else if (part === "null") {
      className = "json-null";
    } else if (/^-?\d+(?:\.\d+)?$/.test(part)) {
      className = "json-number";
    }

    return (
      <span className={className} key={`${index}-${part}`}>
        {part}
      </span>
    );
  });
}
