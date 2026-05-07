<script lang="ts">
  import { onMount } from 'svelte';
  import { project } from './stores/project';
  import { view, type PrimaryView, type RoadmapMode } from './stores/view';
  import Onboarding from './screens/Onboarding.svelte';

  let onboardingOpen = false;

  function preventDrop(event: DragEvent) {
    event.preventDefault();
  }

  onMount(() => {
    window.addEventListener('dragover', preventDrop);
    window.addEventListener('drop', preventDrop);
    return () => {
      window.removeEventListener('dragover', preventDrop);
      window.removeEventListener('drop', preventDrop);
    };
  });

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

  function openOnboarding() {
    onboardingOpen = true;
  }

  function closeOnboarding() {
    onboardingOpen = false;
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
        <button class="primary-action" type="button" on:click={openOnboarding}>Import CSV</button>
      </div>
    </header>

    <main class="content">
      {#if onboardingOpen || $project.status === 'not-imported'}
        <Onboarding on:complete={closeOnboarding} />
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
