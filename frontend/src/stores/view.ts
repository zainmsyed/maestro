import { writable } from 'svelte/store';

export type PrimaryView = 'roadmap' | 'list' | 'health' | 'settings';
export type RoadmapMode = 'gantt' | 'list';

export type ViewState = {
  primary: PrimaryView;
  roadmapMode: RoadmapMode;
};

export const view = writable<ViewState>({
  primary: 'roadmap',
  roadmapMode: 'gantt',
});
