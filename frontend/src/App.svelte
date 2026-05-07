<script lang="ts">
  import { project } from './stores/project';
  import { view, type PrimaryView, type RoadmapMode } from './stores/view';

  const navItems: Array<{ id: PrimaryView; label: string; icon: string }> = [
    { id: 'roadmap', label: 'Roadmap', icon: 'M12 4v16M4 8h16M4 16h16' },
    { id: 'list', label: 'List', icon: 'M6 7h12M6 12h12M6 17h12' },
    { id: 'health', label: 'Health', icon: 'M4 14l4-4 4 5 4-8 4 6' },
    { id: 'settings', label: 'Settings', icon: 'M12 8a4 4 0 100 8 4 4 0 000-8z' },
  ];

  const tabs: Array<{ id: RoadmapMode; label: string }> = [
    { id: 'gantt', label: 'Gantt' },
    { id: 'list', label: 'List' },
  ];

  function setPrimary(primary: PrimaryView) {
    view.update((current) => ({ ...current, primary }));
  }

  function setRoadmapMode(roadmapMode: RoadmapMode) {
    view.update((current) => ({ ...current, roadmapMode }));
  }
</script>

<div class="app-shell">
  <aside class="sidebar" aria-label="Primary navigation">
    <button class="logo" aria-label="Maestro home" on:click={() => setPrimary('roadmap')}>M</button>

    {#each navItems as item}
      <button
        class:active={$view.primary === item.id}
        class="nav-button"
        type="button"
        title={item.label}
        aria-label={item.label}
        on:click={() => setPrimary(item.id)}
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          <path d={item.icon} />
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
              class:active={$view.roadmapMode === tab.id}
              class="view-tab"
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
      <div class="content-inner">
        <section class="hero-card" aria-labelledby="shell-title">
          <p class="eyebrow">Roadmap command center</p>
          <h1 id="shell-title" class="hero-title">Plan work with quiet precision.</h1>
          <p class="hero-copy">
            The Maestro shell is ready. Upcoming stories will fill this canvas with onboarding,
            roadmap lists, gantt timelines, health signals, and settings.
          </p>

          <div class="placeholder-grid" aria-label="Roadmap placeholders">
            <div class="placeholder-card">
              <p class="placeholder-label">Current view</p>
              <div class="placeholder-value">{$view.primary}</div>
            </div>
            <div class="placeholder-card">
              <p class="placeholder-label">Mode</p>
              <div class="placeholder-value">{$view.roadmapMode}</div>
            </div>
            <div class="placeholder-card">
              <p class="placeholder-label">Status</p>
              <div class="placeholder-value">{$project.status}</div>
            </div>
          </div>
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
    </main>

    <footer class="statusbar">
      <div class="status-item">
        <span class:ready={$project.status === 'ready'} class:waiting={$project.status !== 'ready'} class="status-dot"></span>
        <span>{$project.status}</span>
      </div>
      <div class="status-item">API /api</div>
      <div class="status-item">View {$view.primary}</div>
    </footer>
  </section>
</div>
