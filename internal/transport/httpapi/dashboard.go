package httpapi

import (
	"net/http"
)

func HandleDashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>KMIP Lab Dashboard</title>
	<style>
		body {
			margin: 0;
			min-height: 100vh;
			display: grid;
			place-items: center;
			background: #070b16;
			color: #f8fafc;
			font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
		}
		.card {
			width: min(560px, calc(100% - 40px));
			border: 1px solid #223047;
			background: #0f1628;
			border-radius: 18px;
			padding: 32px;
			box-shadow: 0 18px 60px rgba(0, 0, 0, 0.22);
		}
		h1 {
			margin: 0 0 12px;
			font-size: 32px;
		}
		p {
			color: #94a3b8;
			line-height: 1.6;
		}
		a {
			display: inline-block;
			margin-top: 18px;
			background: #2563eb;
			color: white;
			text-decoration: none;
			padding: 12px 16px;
			border-radius: 10px;
			font-weight: 700;
		}
		code {
			color: #38bdf8;
		}
	</style>
</head>
<body>
	<div class="card">
		<h1>KMIP Lab Dashboard</h1>
		<p>
			The Go backend is running. The full React dashboard is available through the Vite frontend.
		</p>
		<p>
			Backend endpoints:
			<br><code>POST /kmip</code>
			<br><code>GET /metrics</code>
			<br><code>GET /keys</code>
			<br><code>GET /audit</code>
			<br><code>GET /health</code>
		</p>
		<a href="http://localhost:5173">Open React Dashboard</a>
	</div>
</body>
</html>`))
	}
}