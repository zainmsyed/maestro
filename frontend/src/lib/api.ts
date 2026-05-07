const API_BASE = import.meta.env.VITE_API_BASE_URL ?? '/api';

type JsonBody = Record<string, unknown> | unknown[];
type RequestOptions = Omit<RequestInit, 'body'> & {
  body?: BodyInit | JsonBody;
};

class ApiError extends Error {
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

export type WorkItemType = 'epic' | 'feature' | 'story';

export type DateAssignmentCandidate = {
  row_number: number;
  work_item_type: WorkItemType;
  id: string;
  title: string;
  assigned_owner: string;
};

export type ImportReport = {
  epic_count: number;
  feature_count: number;
  story_count: number;
  sprints_detected: string[];
  missing_dates_count: number;
  missing_sprint_count: number;
  orphaned_features: number;
  orphaned_stories: number;
  skipped_rows: number;
  detected_date_format: string;
  date_assignment_candidates: DateAssignmentCandidate[];
  ambiguous_dates: unknown[];
  warnings: string[];
  synthetic_story_ids: string[];
};

type DatePatchResponse = {
  id: string;
  original_end_date: string | null;
  committed_end_date: string | null;
  date_source?: string;
};

export function typeEndpoint(type: WorkItemType): string {
  const endpoints: Record<WorkItemType, string> = {
    epic: 'epics',
    feature: 'features',
    story: 'stories',
  };
  return endpoints[type];
}

export const api = {
  get: <T>(path: string, options?: RequestOptions) => apiFetch<T>(path, { ...options, method: 'GET' }),
  post: <T>(path: string, body?: RequestOptions['body'], options?: RequestOptions) =>
    apiFetch<T>(path, { ...options, method: 'POST', body }),
  patch: <T>(path: string, body?: RequestOptions['body'], options?: RequestOptions) =>
    apiFetch<T>(path, { ...options, method: 'PATCH', body }),
  delete: <T>(path: string, options?: RequestOptions) => apiFetch<T>(path, { ...options, method: 'DELETE' }),
  patchDate: (type: WorkItemType, id: string, committedEndDate: string) =>
    apiFetch<DatePatchResponse>(`/${typeEndpoint(type)}/${id}/date`, {
      method: 'PATCH',
      body: {
        committed_end_date: committedEndDate,
        changed_by: 'pm',
        reason: 'Assigned during post-import date assignment',
      },
    }),
};
