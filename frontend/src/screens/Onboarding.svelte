<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte';
  import DateAssignment from './DateAssignment.svelte';
  import { api, type ImportReport } from '../lib/api';
  import { project } from '../stores/project';

  type Step = 'loading' | 'date-assignment' | 'confirm' | 'empty' | 'error';

  const dispatch = createEventDispatcher<{ complete: void }>();

  let step: Step = 'loading';
  let report: ImportReport | null = null;
  let error: string | null = null;
  let assignedCount = 0;

  async function loadImportReport() {
    step = 'loading';
    error = null;
    try {
      report = await api.get<ImportReport>('/import/report');
      if ((report.date_assignment_candidates?.length ?? 0) === 0 || report.missing_dates_count === 0) {
        step = 'confirm';
      } else {
        step = 'date-assignment';
      }
    } catch (err) {
      error = err instanceof Error ? err.message : 'Could not load import report.';
      step = 'error';
    }
  }

  function continueToConfirm(event: CustomEvent<{ assigned: number; skipped: number }>) {
    assignedCount = event.detail.assigned;
    step = 'confirm';
  }

  function skipToConfirm() {
    step = 'confirm';
  }

  function completeOnboarding() {
    const currentReport = report;
    if (currentReport) {
      project.update((current) => ({
        ...current,
        name: `${currentReport.epic_count} epics, ${currentReport.feature_count} features`,
        status: 'ready',
        lastSync: new Date().toISOString(),
      }));
    }
    dispatch('complete');
  }

  onMount(loadImportReport);
</script>

<section class="onboarding" aria-label="Roadmap onboarding">
  {#if step === 'loading'}
    <div class="state-card">Loading latest import report…</div>
  {:else if step === 'error'}
    <div class="state-card error">
      <p>{error}</p>
      <button type="button" on:click={loadImportReport}>Retry</button>
    </div>
  {:else if step === 'date-assignment' && report}
    <DateAssignment {report} on:done={continueToConfirm} on:skip={skipToConfirm} />
  {:else if step === 'confirm' && report}
    <div class="confirm-card">
      <p class="eyebrow">Ready</p>
      <h2>Roadmap import reviewed</h2>
      <p class="summary">
        {report.epic_count} epics, {report.feature_count} features, and {report.story_count} stories are ready.
      </p>
      <div class="summary-grid">
        <div><span>{report.missing_dates_count}</span><small>Missing dates</small></div>
        <div><span>{assignedCount}</span><small>Assigned now</small></div>
        <div><span>{report.sprints_detected.length}</span><small>Sprints</small></div>
      </div>
      <button class="btn-primary" type="button" on:click={completeOnboarding}>Continue</button>
    </div>
  {:else}
    <div class="state-card">No import report is available yet.</div>
  {/if}
</section>

<style>
  .onboarding {
    width: min(1180px, 100%);
    margin: 0 auto;
    padding: 32px;
  }

  .state-card,
  .confirm-card {
    padding: 24px;
    border: 1px solid var(--border);
    border-radius: 10px;
    background: linear-gradient(145deg, rgba(29, 29, 32, 0.94), rgba(22, 22, 24, 0.94));
    color: var(--text2);
  }

  .state-card.error {
    border-color: rgba(196, 97, 74, 0.45);
    color: var(--coral);
  }

  .eyebrow {
    margin: 0 0 8px;
    color: var(--text3);
    font-size: 13px;
    letter-spacing: 0.12em;
    text-transform: uppercase;
  }

  h2 {
    margin: 0;
    color: var(--text);
    font-family: var(--font-display);
    font-size: 30px;
    font-weight: 300;
    letter-spacing: -0.04em;
  }

  .summary {
    margin: 12px 0 0;
    color: var(--text2);
    font-size: 15px;
  }

  .summary-grid {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 12px;
    margin: 24px 0;
  }

  .summary-grid div {
    padding: 16px;
    border: 1px solid var(--border);
    border-radius: 8px;
    background: var(--bg2);
  }

  .summary-grid span {
    display: block;
    color: var(--text);
    font-family: var(--font-display);
    font-size: 30px;
    font-weight: 300;
  }

  .summary-grid small {
    display: block;
    margin-top: 4px;
    color: var(--text3);
    font-size: 13px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  button {
    padding: 10px 16px;
    border: 1px solid var(--accent2);
    border-radius: 6px;
    background: transparent;
    color: var(--accent);
    cursor: pointer;
    font-family: var(--font-mono);
    font-size: 13px;
    letter-spacing: 0.04em;
  }

  button:hover {
    background: rgba(200, 169, 110, 0.08);
  }
</style>
