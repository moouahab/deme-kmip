export type KMIPResponse = {
  status: number;
  ok: boolean;
  data: unknown;
};

export type Metrics = {
  http_requests_total: number;
  http_errors_total: number;
  create_key_total: number;
  get_key_total: number;
  destroy_key_total: number;
  success_total: number;
  not_found_total: number;
};

export type KMSKey = {
  id: string;
  created_at: string;
  object_type: number;
  status: string;
  updated_at: string;
};

export type AuditEvent = {
  time: string;
  operation: string;
  key_id?: string;
  status?: string;
  result: string;
  error?: string;
};

export const emptyMetrics: Metrics = {
  http_requests_total: 0,
  http_errors_total: 0,
  create_key_total: 0,
  get_key_total: 0,
  destroy_key_total: 0,
  success_total: 0,
  not_found_total: 0,
};

export async function sendKMIPRequest(body: Uint8Array): Promise<KMIPResponse> {
  const arrayBuffer = new ArrayBuffer(body.byteLength);
  const view = new Uint8Array(arrayBuffer);
  view.set(body);

  const response = await fetch("/kmip", {
    method: "POST",
    headers: {
      "Content-Type": "application/octet-stream",
    },
    body: arrayBuffer,
  });

  const data = await response.json();

  return {
    status: response.status,
    ok: response.ok,
    data,
  };
}

export async function getMetrics(): Promise<Metrics> {
  const response = await fetch("/metrics");

  if (!response.ok) {
    throw new Error("cannot fetch metrics");
  }

  return response.json();
}

export async function getKeys(): Promise<KMSKey[]> {
  const response = await fetch("/keys");

  if (!response.ok) {
    throw new Error("cannot fetch keys");
  }

  return response.json();
}

export async function getAuditEvents(): Promise<AuditEvent[]> {
  const response = await fetch("/audit");

  if (!response.ok) {
    throw new Error("cannot fetch audit events");
  }

  return response.json();
}