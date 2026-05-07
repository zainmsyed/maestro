<script lang="ts">
  import { project } from './stores/project';
  import { view, type PrimaryView, type RoadmapMode } from './stores/view';
  import DateAssignment from './screens/DateAssignment.svelte';
  import type { ImportReport } from './lib/api';

  // TODO(story-007): remove /date_screen temporary test mount once onboarding flow is wired
  const showDateScreen = typeof window !== 'undefined' && window.location.pathname === '/date_screen';

  const mockReport: ImportReport = {
    epic_count: 3,
    feature_count: 12,
    story_count: 28,
    sprints_detected: ['Sprint 1', 'Sprint 2'],
    missing_dates_count: 3,
    missing_sprint_count: 0,
    orphaned_features: 0,
    orphaned_stories: 0,
    skipped_rows: 0,
    detected_date_format: 'YYYY-MM-DD',
    date_assignment_candidates: [
      { row_number: 2, work_item_type: 'story', id: 'S-501', title: 'Build login flow', assigned_owner: 'Alice' },
      { row_number: 5, work_item_type: 'feature', id: 'F-203', title: 'User dashboard', assigned_owner: 'Bob' },
      { row_number: 8, work_item_type: 'epic', id: 'E-102', title: 'Platform v2', assigned_owner: 'Charlie' },
    ],
    ambiguous_dates: [],
    warnings: [],
    synthetic_story_ids: [],
  };

  type NavItem = { id: PrimaryView; label: string; icon: string };

  const navItems: NavItem[] = [
    { id: 'roadmap', label: 'Roadmap', icon: 'grid' },
    { id: 'list',    label: 'List',    icon: 'list' },
    { id: 'health',  label: 'Health',  icon: 'pulse' },
    { id: 'settings',label: 'Settings',icon: 'gear' },
  ];

  const tabs: Array<{ id: RoadmapMode; label: string }> = [
    { id: 'gantt', label: 'Gantt' },
    { id: 'list',  label: 'List' },
  ];

  function setPrimary(primary: PrimaryView) {
    view.update((c) => ({ ...c, primary }));
  }

  function setRoadmapMode(roadmapMode: RoadmapMode) {
    view.update((c) => ({ ...c, roadmapMode }));
  }

  const icons: Record<string, string> = {
    grid:  '<rect x="3"  y="3"  width="7" height="7" rx="1.5"/><rect x="14" y="3"  width="7" height="7" rx="1.5"/><rect x="3"  y="14" width="7" height="7" rx="1.5"/><rect x="14" y="14" width="7" height="7" rx="1.5"/>',
    list:  '<path d="M4 7h16M4 12h16M4 17h16"/>',
    pulse: '<path d="M22 12h-4l-3 9L9 3l-3 9H2"/>',
    gear:  '<circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>',
  };
</script>

<div class="app-shell">
  <aside class="sidebar" aria-label="Primary navigation">
    <button class="logo" aria-label="Maestro home" on:click={() => setPrimary('roadmap')}>M</button>

    {#each navItems as item}
      <button
        class="nav-button"
        class:active={$view.primary === item.id}
        type="button"
        title={item.label}
        aria-label={item.label}
        on:click={() => setPrimary(item.id)}
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          {@html icons[item.icon]}
        </svg>
      </button>
    {/each}

    <div class="sidebar-spacer"></div>
  </aside>

  <section class="shell-main">
    <header class="topbar">
      <div class="topbar-title">Maestro</div>
      <div class="topbar-sep">/</div>
      <div class="topbar-project">{$project.name}</div>

      <div class="topbar-right">
        <div class="view-tabs" aria-label="Roadmap view mode">
          {#each tabs as tab}
            <button
              class="view-tab"
              class:active={$view.roadmapMode === tab.id}
              type="button"
              on:click={() => setRoadmapMode(tab.id)}
            >
              {tab.label}
            </button>
          {/each}
        </div>
        <button class="primary-action" type="button">Import CSV</button>
      </div>
    </header>

    <main class="content">
      {#if showDateScreen}
        <DateAssignment report={mockReport} on:done={() => alert('done')} on:skip={() => alert('skipped')} />
      {:else if $project.status === 'not-imported'}
        <div class="empty-state">
          <div class="empty-card">
            <p class="empty-eyebrow">Getting started</p>
            <h1 class="empty-title">No roadmap loaded</h1>
            <p class="empty-body">
              Import an Azure DevOps Tree of Work Items CSV to build your first roadmap.
            </p>
            <button class="empty-cta" type="button">Import CSV</button>
          </div>
        </div>
      {:else}
        <div class="content-inner">
          <section class="hero-card">
            <p class="eyebrow">Roadmap</p>
            <h1 class="hero-title">{$project.name}</h1>
            <p class="hero-copy">
              The roadmap canvas will appear here once the relevant screen stories are implemented.
            </p>
          </section>

          <aside class="panel-stack" aria-label="Status panels">
            <section class="panel-card">
              <p class="panel-label">Project</p>
              <p class="panel-text">{$project.name}</p>
            </section>
            <section class="panel-card">
              <p class="panel-label">Sync</p>
              <p class="panel-text">{$project.lastSync ?? 'No sync recorded'}</p>
            </section>
          </aside>
        </div>
      {/if}
    </main>

    <footer class="statusbar">
      <div class="status-item">
        <span class="status-dot" class:ready={$project.status === 'ready'} class:waiting={$project.status !== 'ready'}></span>
        <span>{$project.status}</span>
      </div>
      <div class="status-item">API /api</div>
      <div class="status-item">View {$view.primary}</div>
    </footer>
  </section>
</div>
