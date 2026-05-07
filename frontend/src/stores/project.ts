import { writable } from 'svelte/store';

export type ProjectStatus = 'not-imported' | 'ready' | 'syncing' | 'error';

export type ProjectState = {
  name: string;
  status: ProjectStatus;
  lastSync: string | null;
};

export const project = writable<ProjectState>({
  name: 'No roadmap loaded',
  status: 'not-imported',
  lastSync: null,
});
