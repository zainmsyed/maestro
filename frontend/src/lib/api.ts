const API_BASE = import.meta.env.VITE_API_BASE_URL ?? '/api';

type JsonBody = Record<string, unknown> | unknown[];
type RequestOptions = Omit<RequestInit, 'body'> & {
  body?: BodyInit | JsonBody;
};

export class ApiError extends Error {
  constructor(
    message: string,
    readonly status: number,
    readonly details?: unknown,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

function buildHeaders(body: RequestOptions['body'], headers?: HeadersInit): Headers {
  const next = new Headers(headers);
  if (body !== undefined && !(body instanceof FormData) && !next.has('Content-Type')) {
    next.set('Content-Type', 'application/json');
  }
  if (!next.has('Accept')) {
    next.set('Accept', 'application/json');
  }
  return next;
}

async function parseResponse(response: Response): Promise<unknown> {
  if (response.status === 204) return null;
  const contentType = response.headers.get('content-type') ?? '';
  if (contentType.includes('application/json')) return response.json();
  return response.text();
}

export async function apiFetch<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const url = `${API_BASE}${path.startsWith('/') ? path : `/${path}`}`;
  const { body, headers, ...init } = options;
  const response = await fetch(url, {
    ...init,
    headers: buildHeaders(body, headers),
    body: body === undefined || body instanceof FormData ? body : JSON.stringify(body),
  });

  const data = await parseResponse(response);
  if (!response.ok) {
    const message =
      typeof data === 'object' && data !== null && 'error' in data
        ? String((data as { error: unknown }).error)
        : `Request failed with status ${response.status}`;
    throw new ApiError(message, response.status, data);
  }

  return data as T;
}

export const api = {
  get: <T>(path: string, options?: RequestOptions) => apiFetch<T>(path, { ...options, method: 'GET' }),
  post: <T>(path: string, body?: RequestOptions['body'], options?: RequestOptions) =>
    apiFetch<T>(path, { ...options, method: 'POST', body }),
  patch: <T>(path: string, body?: RequestOptions['body'], options?: RequestOptions) =>
    apiFetch<T>(path, { ...options, method: 'PATCH', body }),
  delete: <T>(path: string, options?: RequestOptions) => apiFetch<T>(path, { ...options, method: 'DELETE' }),
};
