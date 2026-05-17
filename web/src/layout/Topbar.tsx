export function Topbar() {
  function refreshPage() {
    window.location.reload();
  }

  return (
    <header className="topbar">
      <div>
        <h1>KMS Operation Console</h1>
        <p>Browser-based KMIP-like operation tester using TTLV encoded requests</p>
      </div>

      <div className="topbar-actions">
        <span className="badge warning">● Local Dev</span>
        <span className="badge">localhost:8080</span>
        <button onClick={refreshPage}>↻ Refresh</button>
      </div>
    </header>
  );
}