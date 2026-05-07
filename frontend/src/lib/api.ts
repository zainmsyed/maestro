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

async function apiFetch<T>(path: string, options: RequestOptions = {}): Promise<T> {
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

export type SprintDraft = {
  name: string;
  start_date: string;
  length_days: number;
  source: 'imported' | 'generated';
};

export type ManualSprintConfig = {
  sprint_length_days: number;
  first_start_date: string;
  count: number;
};

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
  existing_skipped: number;
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

export type StoryRecord = {
  id: string;
  feature_id: string;
  title: string;
  description: string;
  status: string;
  owner: string;
  sprint: string;
  story_points: number | null;
  original_end_date: string | null;
  committed_end_date: string | null;
  actual_end_date: string | null;
  date_source: string;
  created_at: string;
  updated_at: string;
};

export type FeatureRecord = {
  id: string;
  epic_id: string | null;
  title: string;
  description: string;
  status: string;
  owner: string;
  sprint: string;
  story_points: number | null;
  original_end_date: string | null;
  committed_end_date: string | null;
  actual_end_date: string | null;
  date_source: string;
  stories: StoryRecord[];
  created_at: string;
  updated_at: string;
};

export type EpicRecord = {
  id: string;
  title: string;
  description: string;
  status: string;
  owner: string;
  sprint_start: string;
  sprint_end: string;
  original_end_date: string | null;
  committed_end_date: string | null;
  actual_end_date: string | null;
  is_synthetic: boolean;
  features: FeatureRecord[];
  created_at: string;
  updated_at: string;
};

export type AuditRecord = {
  id: number;
  entity_type: string;
  entity_id: string;
  changed_by: string;
  old_date: string | null;
  new_date: string | null;
  delta_days: number;
  reason: string | null;
  changed_at: string;
  date_source?: string;
};

type DatePatchResponse = {
  id: string;
  original_end_date: string | null;
  committed_end_date: string | null;
  date_source?: string;
};

function typeEndpoint(type: WorkItemType): string {
  const endpoints: Record<WorkItemType, string> = {
    epic: 'epics',
    feature: 'features',
    story: 'stories',
  };
  return endpoints[type];
}

export function formatApiError(error: unknown, fallback = 'Request failed.'): string {
  if (error instanceof ApiError) return error.message;
  if (error instanceof Error) return error.message;
  return fallback;
}

export const api = {
  get: <T>(path: string, options?: RequestOptions) => apiFetch<T>(path, { ...options, method: 'GET' }),
  post: <T>(path: string, body?: RequestOptions['body'], options?: RequestOptions) =>
    apiFetch<T>(path, { ...options, method: 'POST', body }),
  patch: <T>(path: string, body?: RequestOptions['body'], options?: RequestOptions) =>
    apiFetch<T>(path, { ...options, method: 'PATCH', body }),
  delete: <T>(path: string, options?: RequestOptions) => apiFetch<T>(path, { ...options, method: 'DELETE' }),
  importCsv: (file: File) => {
    const form = new FormData();
    form.append('file', file);
    return apiFetch<ImportReport>('/import', { method: 'POST', body: form });
  },
  patchDate: (type: WorkItemType, id: string, committedEndDate: string) =>
    apiFetch<DatePatchResponse>(`/${typeEndpoint(type)}/${id}/date`, {
      method: 'PATCH',
      body: {
        committed_end_date: committedEndDate,
        changed_by: 'pm',
        reason: 'Assigned during list view editing',
      },
    }),
  patchFeatureEpic: (id: string, epicID: string) =>
    apiFetch<FeatureRecord>(`/features/${id}/epic`, {
      method: 'PATCH',
      body: { epic_id: epicID },
    }),
};
