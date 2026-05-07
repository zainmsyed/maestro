import { writable } from 'svelte/store';
import type { ImportReport, ManualSprintConfig, SprintDraft } from '../lib/api';

export type ProjectStatus = 'not-imported' | 'ready' | 'syncing' | 'error';

export type ProjectState = {
  name: string;
  status: ProjectStatus;
  lastSync: string | null;
  sourceFileName: string | null;
  importReport: ImportReport | null;
  sprintDrafts: SprintDraft[];
  manualSprintConfig: ManualSprintConfig | null;
};

export const project = writable<ProjectState>({
  name: 'No roadmap loaded',
  status: 'not-imported',
  lastSync: null,
  sourceFileName: null,
  importReport: null,
  sprintDrafts: [],
  manualSprintConfig: null,
});
