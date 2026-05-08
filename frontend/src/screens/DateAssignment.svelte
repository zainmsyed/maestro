<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte';
  import DateAssignmentRow from '../components/DateAssignmentRow.svelte';
  import { api, type DateAssignmentCandidate, type ImportReport } from '../lib/api';

  type DetailWithSprint = {
    id: string;
    sprint?: string;
    sprint_end?: string;
  };

  type AssignmentStatus = 'idle' | 'saving' | 'success' | 'error';

  type AssignmentRow = {
    candidate: DateAssignmentCandidate;
    sprintName: string;
    selected: boolean;
    dateValue: string;
    snapToSprint: boolean;
    status: AssignmentStatus;
    error: string | null;
  };

  export let report: ImportReport;

  const dispatch = createEventDispatcher<{
    done: { assigned: number; skipped: number };
    skip: void;
  }>();

  let rows: AssignmentRow[] = [];
  let loading = true;
  let assigning = false;
  let globalError: string | null = null;
  let successMessage: string | null = null;

  $: candidateCount = report?.date_assignment_candidates?.length ?? 0;
  $: selectedCount = rows.filter((row) => row.selected).length;
  $: datedCount = rows.filter((row) => row.dateValue).length;

  function validateDate(value: string): string | null {
    if (!value) return 'Choose a target date.';
    if (!/^\d{4}-\d{2}-\d{2}$/.test(value)) return 'Use YYYY-MM-DD.';
    return null;
  }

  async function loadSprintNames(): Promise<Record<string, string>> {
    const sprintNames: Record<string, string> = {};
    const [stories, features, epics] = await Promise.all([
      api.get<DetailWithSprint[]>('/stories').catch(() => []),
      api.get<DetailWithSprint[]>('/features').catch(() => []),
      api.get<DetailWithSprint[]>('/epics').catch(() => []),
    ]);

    for (const item of stories) sprintNames[item.id] = item.sprint ?? '';
    for (const item of features) sprintNames[item.id] = item.sprint ?? '';
    for (const item of epics) sprintNames[item.id] = item.sprint_end ?? '';
    return sprintNames;
  }

  async function initializeRows() {
    loading = true;
    globalError = null;
    try {
      if (candidateCount === 0) {
        rows = [];
        dispatch('done', { assigned: 0, skipped: 0 });
        return;
      }

      const sprintNames = await loadSprintNames();
      rows = report.date_assignment_candidates.map((candidate) => ({
        candidate,
        sprintName: sprintNames[candidate.id] ?? '',
        selected: false,
        dateValue: '',
        snapToSprint: true,
        status: 'idle',
        error: null,
      }));
    } catch (err) {
      globalError = err instanceof Error ? err.message : 'Could not load date assignment rows.';
    } finally {
      loading = false;
    }
  }

  async function assignRows(targetRows: AssignmentRow[], continueAfter = false) {
    successMessage = null;
    globalError = null;
    if (targetRows.length === 0) {
      globalError = 'Select at least one row or enter a date before assigning.';
      return;
    }

    let valid = true;
    rows = rows.map((row) => {
      if (!targetRows.includes(row)) return row;
      const error = validateDate(row.dateValue);
      if (error) valid = false;
      return { ...row, error };
    });
    if (!valid) return;

    assigning = true;
    let assigned = 0;
    let failed = 0;

    for (const row of targetRows) {
      rows = rows.map((current) =>
        current === row ? { ...current, status: 'saving', error: null } : current
      );

      try {
        await api.patchDate(row.candidate.work_item_type, row.candidate.id, row.dateValue);
        assigned += 1;
        rows = rows.map((current) =>
          current === row ? { ...current, status: 'success', selected: false } : current
        );
      } catch (err) {
        failed += 1;
        const message = err instanceof Error ? err.message : 'Failed to assign date.';
        rows = rows.map((current) =>
          current === row ? { ...current, status: 'error', error: message } : current
        );
      }
    }

    assigning = false;
    successMessage = failed === 0 ? `Assigned ${assigned} date${assigned === 1 ? '' : 's'}.` : `Assigned ${assigned}; ${failed} failed.`;

    if (continueAfter && failed === 0) {
      dispatch('done', { assigned, skipped: rows.length - assigned });
    }
  }

  function assignSelected() {
    assignRows(rows.filter((row) => row.selected));
  }

  function assignDatedAndContinue() {
    assignRows(rows.filter((row) => row.dateValue), true);
  }

  function skipForLater() {
    dispatch('skip');
  }

  onMount(initializeRows);
</script>

<section class="date-assignment" aria-labelledby="date-assignment-title">
  <header class="screen-header">
    <div>
      <p class="eyebrow">Post-import cleanup</p>
      <h2 id="date-assignment-title">Assign missing target dates</h2>
      <p class="intro">
        {candidateCount} imported item{candidateCount === 1 ? '' : 's'} need target dates before the roadmap is locked.
      </p>
    </div>
    <div class="header-actions">
      <button class="btn-ghost" type="button" disabled={assigning} on:click={skipForLater}>
        Skip — I'll do it later
      </button>
      <div class="count-pill">{candidateCount} missing</div>
    </div>
  </header>

  {#if loading}
    <div class="state-card">Loading date assignment rows…</div>
  {:else if candidateCount === 0}
    <div class="state-card">All imported items already have target dates.</div>
  {:else}
    <div class="toolbar">
      <div class="toolbar-meta">{selectedCount} selected · {datedCount} dated</div>
      <button class="btn-secondary" type="button" disabled={assigning || selectedCount === 0} on:click={assignSelected}>
        Set all selected
      </button>
    </div>

    <div class="rows" aria-label="Missing target date items">
      {#each rows as row (row.candidate.id)}
        <DateAssignmentRow
          candidate={row.candidate}
          sprintName={row.sprintName}
          bind:selected={row.selected}
          bind:dateValue={row.dateValue}
          bind:snapToSprint={row.snapToSprint}
          status={row.status}
          error={row.error}
          disabled={assigning}
        />
      {/each}
    </div>

    {#if globalError}
      <div class="error-banner">{globalError}</div>
    {/if}

    {#if successMessage}
      <div class="success-banner">{successMessage}</div>
    {/if}

    <footer class="actions">
      <button class="btn-ghost" type="button" disabled={assigning} on:click={skipForLater}>
        Skip — I'll do it later
      </button>
      <button class="btn-primary" type="button" disabled={assigning || datedCount === 0} on:click={assignDatedAndContinue}>
        {assigning ? 'Assigning…' : 'Set dates now'}
      </button>
    </footer>
  {/if}
</section>

<style>
  .date-assignment {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .screen-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 24px;
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

  .intro {
    max-width: 620px;
    margin: 12px 0 0;
    color: var(--text2);
    font-size: 15px;
  }

  .header-actions {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-shrink: 0;
  }

  .count-pill {
    flex-shrink: 0;
    padding: 6px 10px;
    border: 1px solid var(--accent2);
    border-radius: 999px;
    color: var(--accent);
    font-size: 13px;
  }

  .toolbar,
  .actions {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
  }

  .toolbar-meta {
    color: var(--text3);
    font-size: 13px;
  }

  .rows {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .state-card,
  .error-banner,
  .success-banner {
    padding: 16px 18px;
    border: 1px solid var(--border);
    border-radius: 8px;
    background: var(--bg2);
    color: var(--text2);
  }

  .error-banner {
    border-color: color-mix(in srgb, var(--coral) 45%, transparent);
    color: var(--coral);
  }

  .success-banner {
    border-color: color-mix(in srgb, var(--green) 45%, transparent);
    color: var(--green-light);
  }

  button {
    border-radius: 6px;
    cursor: pointer;
    font-family: var(--font-mono);
    font-size: 13px;
    letter-spacing: 0.04em;
    transition: background 0.15s ease, border-color 0.15s ease, color 0.15s ease;
  }

  button:disabled {
    cursor: not-allowed;
    opacity: 0.45;
  }

  .btn-primary,
  .btn-secondary,
  .btn-ghost {
    padding: 10px 16px;
    border: 1px solid var(--border);
    background: transparent;
    color: var(--text2);
  }

  .btn-primary,
  .btn-secondary {
    border-color: var(--accent2);
    color: var(--accent);
  }

  .btn-primary:hover:not(:disabled),
  .btn-secondary:hover:not(:disabled) {
    background: var(--accent-a8);
  }

  .btn-ghost:hover:not(:disabled) {
    border-color: var(--border2);
    color: var(--text);
  }
</style>
