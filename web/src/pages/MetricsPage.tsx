import { useEffect, useState } from "react";
import { emptyMetrics, getMetrics } from "../services/api";
import type { Metrics } from "../services/api";
import { StatCard } from "../layout/StatCard";

export function MetricsPage() {
  const [metrics, setMetrics] = useState<Metrics>(emptyMetrics);

  async function loadMetrics() {
    try {
      const data = await getMetrics();
      setMetrics(data);
    } catch {
      // ignore for now
    }
  }

  useEffect(() => {
    loadMetrics();

    const interval = window.setInterval(loadMetrics, 2000);

    return () => window.clearInterval(interval);
  }, []);

  const totalResults =
    metrics.success_total + metrics.not_found_total + metrics.http_errors_total;

  const successRate = totalResults
    ? Math.round((metrics.success_total / totalResults) * 1000) / 10
    : 0;

  const errorRate = totalResults
    ? Math.round((metrics.http_errors_total / totalResults) * 1000) / 10
    : 0;

  const notFoundRate = totalResults
    ? Math.round((metrics.not_found_total / totalResults) * 1000) / 10
    : 0;

  const maxOperationValue = Math.max(
    metrics.create_key_total,
    metrics.get_key_total,
    metrics.destroy_key_total,
    1
  );

  return (
    <>
      <section className="page-header">
        <h2>Metrics</h2>
        <p>Monitor server performance and operation statistics</p>
      </section>

      <section className="panel">
        <h3>Metrics Overview</h3>
        <p>GET /metrics</p>

        <div className="stats-grid metrics">
          <StatCard
            label="HTTP Requests"
            value={metrics.http_requests_total}
            icon="↔"
            tone="blue"
          />

          <StatCard
            label="HTTP Errors"
            value={metrics.http_errors_total}
            icon="△"
            tone="red"
          />

          <StatCard
            label="Create Key"
            value={metrics.create_key_total}
            icon="+"
            tone="green"
          />

          <StatCard
            label="Get Key"
            value={metrics.get_key_total}
            icon="◎"
            tone="blue"
          />

          <StatCard
            label="Destroy Key"
            value={metrics.destroy_key_total}
            icon="▥"
            tone="orange"
          />

          <StatCard
            label="Success Total"
            value={metrics.success_total}
            icon="✓"
            tone="green"
          />

          <StatCard
            label="Not Found"
            value={metrics.not_found_total}
            icon="×"
            tone="orange"
          />
        </div>
      </section>

      <section className="metrics-layout">
        <div className="panel">
          <h3>Operations Distribution</h3>

          <div className="bar-chart">
            <Bar
              label="Create Key"
              value={metrics.create_key_total}
              max={maxOperationValue}
              tone="green"
            />

            <Bar
              label="Get Key"
              value={metrics.get_key_total}
              max={maxOperationValue}
              tone="blue"
            />

            <Bar
              label="Destroy Key"
              value={metrics.destroy_key_total}
              max={maxOperationValue}
              tone="orange"
            />
          </div>
        </div>

        <div className="panel">
          <h3>Results Distribution</h3>

          <div className="result-list">
            <ResultLine
              label="Success"
              value={metrics.success_total}
              tone="green"
            />

            <ResultLine
              label="Not Found"
              value={metrics.not_found_total}
              tone="orange"
            />

            <ResultLine
              label="Errors"
              value={metrics.http_errors_total}
              tone="red"
            />
          </div>
        </div>
      </section>

      <section className="panel">
        <h3>Server Health</h3>

        <div className="health-grid">
          <HealthMetric
            label="Success Rate"
            value={successRate}
            tone="green"
          />

          <HealthMetric
            label="Error Rate"
            value={errorRate}
            tone="red"
          />

          <HealthMetric
            label="Not Found Rate"
            value={notFoundRate}
            tone="orange"
          />
        </div>
      </section>
    </>
  );
}

function Bar({
  label,
  value,
  max,
  tone,
}: {
  label: string;
  value: number;
  max: number;
  tone: string;
}) {
  const height = max > 0 ? Math.max(6, Math.min(180, (value / max) * 180)) : 6;

  return (
    <div className="bar-item">
      <div className="bar-track">
        <div className={`bar ${tone}`} style={{ height }}></div>
      </div>
      <span>{label}</span>
    </div>
  );
}

function ResultLine({
  label,
  value,
  tone,
}: {
  label: string;
  value: number;
  tone: string;
}) {
  return (
    <div className="result-line">
      <span>
        <span className={`legend-dot ${tone}`}></span>
        {label}
      </span>
      <strong>{value}</strong>
    </div>
  );
}

function HealthMetric({
  label,
  value,
  tone,
}: {
  label: string;
  value: number;
  tone: string;
}) {
  return (
    <div>
      <span className="health-label">{label}</span>
      <strong className={`health-value ${tone}`}>{value}%</strong>

      <div className="progress">
        <div className={tone} style={{ width: `${value}%` }}></div>
      </div>
    </div>
  );
}