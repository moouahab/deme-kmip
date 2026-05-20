type TopbarProps = {
  onMenuClick?: () => void;
};

export function Topbar({ onMenuClick }: TopbarProps) {
  function refreshPage() {
    window.location.reload();
  }

  return (
    <header className="topbar">
      <div className="topbar-title-row">
        <button className="menu-button" onClick={onMenuClick} aria-label="Open navigation menu">
          ☰
        </button>

        <div>
          <h1>KMS Operation Console</h1>
          <p>Browser-based KMIP-like operation tester using TTLV encoded requests</p>
        </div>
      </div>

      <div className="topbar-actions">
        <span className="badge warning">● Local Dev</span>
        <span className="badge">localhost:8080</span>
        <button onClick={refreshPage}>↻ Refresh</button>
      </div>
    </header>
  );
}
