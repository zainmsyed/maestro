<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import DropZone from '../components/DropZone.svelte';
  import SprintPreview from '../components/SprintPreview.svelte';
  import StepBar from '../components/StepBar.svelte';
  import DateAssignment from './DateAssignment.svelte';
  import {
    api,
    formatApiError,
    type ImportReport,
    type ManualSprintConfig,
    type SprintDraft,
  } from '../lib/api';
  import { project } from '../stores/project';
  import { view } from '../stores/view';

  type PrimaryStep = 'import' | 'sprints' | 'confirm';
  type FlowStep = PrimaryStep | 'date-assignment';

  const dispatch = createEventDispatcher<{ complete: void }>();

  const defaultManualConfig = (): ManualSprintConfig => ({
    sprint_length_days: 14,
    first_start_date: new Date().toISOString().slice(0, 10),
    count: 6,
  });

  let step: FlowStep = 'import';
  let report: ImportReport | null = null;
  let selectedFileName = '';
  let importing = false;
  let importError: string | null = null;
  let assignedCount = 0;
  let skippedDateCount = 0;
  let sprintDrafts: SprintDraft[] = [];
  let manualSprintMode = false;
  let manualSprintConfig: ManualSprintConfig = defaultManualConfig();

  $: activeStep = step === 'date-assignment' ? 'confirm' : step;
  $: previewCounts = report
    ? [
        { label: 'Epics', value: report.epic_count },
        { label: 'Features', value: report.feature_count },
        { label: 'Stories', value: report.story_count },
        { label: 'Sprints', value: report.sprints_detected.length },
        { label: 'Missing dates', value: report.missing_dates_count },
        { label: 'Orphans', value: report.orphaned_features + report.orphaned_stories },
      ]
    : [];

  function importedSprintDrafts(nextReport: ImportReport, config: ManualSprintConfig): SprintDraft[] {
    return nextReport.sprints_detected.map((name) => ({
      name,
      start_date: '',
      length_days: config.sprint_length_days,
      source: 'imported' as const,
    }));
  }

  function normalizeReport(raw: ImportReport): ImportReport {
    return {
      ...raw,
      sprints_detected: raw.sprints_detected ?? [],
      date_assignment_candidates: raw.date_assignment_candidates ?? [],
      ambiguous_dates: raw.ambiguous_dates ?? [],
      warnings: raw.warnings ?? [],
      synthetic_story_ids: raw.synthetic_story_ids ?? [],
    };
  }

  function initSprintDrafts(nextReport: ImportReport, config: ManualSprintConfig): SprintDraft[] {
    if (nextReport.sprints_detected.length === 0) return [];
    return importedSprintDrafts(nextReport, config);
  }

  function unsupportedFileMessage(file: File): string | null {
    const isCsv = file.name.toLowerCase().endsWith('.csv') || file.type === 'text/csv' || file.type === 'application/vnd.ms-excel';
    return isCsv ? null : 'Unsupported format. Import an Azure DevOps CSV export.';
  }

  function describeImportError(error: unknown): string {
    const message = formatApiError(error, 'Could not import the CSV.');
    if (message.includes('404') || message.includes('Method Not Allowed')) {
      return 'The import service did not accept the upload. Verify that POST /api/import is available.';
    }
    return message;
  }

  async function importFile(event: CustomEvent<File>) {
    const file = event.detail;
    const formatError = unsupportedFileMessage(file);
    if (formatError) {
      importError = formatError;
      return;
    }

    importing = true;
    importError = null;
    selectedFileName = file.name;
    report = null;
    sprintDrafts = [];
    assignedCount = 0;
    skippedDateCount = 0;

    try {
      const rawReport = await api.importCsv(file);
      report = normalizeReport(rawReport);
      manualSprintMode = report.sprints_detected.length === 0;
      manualSprintConfig = {
        ...defaultManualConfig(),
        count: Math.max(report.sprints_detected.length, 6),
      };
      sprintDrafts = initSprintDrafts(report, manualSprintConfig);
      step = 'sprints';
    } catch (error) {
      importError = describeImportError(error);
      step = 'import';
    } finally {
      importing = false;
    }
  }

  function handleInvalidFile(event: CustomEvent<{ message: string }>) {
    importError = event.detail.message;
  }

  function backToImport() {
    step = 'import';
  }

  function continueFromSprints() {
    if (!report) return;
    if (report.missing_dates_count > 0 && report.date_assignment_candidates.length > 0) {
      step = 'date-assignment';
      return;
    }
    step = 'confirm';
  }

  function continueToConfirm(event: CustomEvent<{ assigned: number; skipped: number }>) {
    assignedCount = event.detail.assigned;
    skippedDateCount = event.detail.skipped;
    step = 'confirm';
  }

  function skipDateAssignment() {
    assignedCount = 0;
    skippedDateCount = report?.date_assignment_candidates.length ?? 0;
    step = 'confirm';
  }

  function buildRoadmap() {
    if (!report) return;
    project.set({
      name: `${report.epic_count + report.feature_count + report.story_count} item roadmap`,
      status: 'ready',
      lastSync: new Date().toISOString(),
      sourceFileName: selectedFileName || null,
      importReport: report,
      sprintDrafts,
      manualSprintConfig: manualSprintMode ? manualSprintConfig : null,
    });
    view.update((current) => ({ ...current, primary: 'list', roadmapMode: 'list' }));
    dispatch('complete');
  }
</script>

<section class="onboarding" aria-label="Roadmap onboarding">
  <StepBar active={activeStep} />

  {#if step === 'import'}
    <div class="intro-grid">
      <section class="copy-card">
        <p class="eyebrow">Maestro setup</p>
        <h1>Build the roadmap from your Azure DevOps export</h1>
        <p>
          Start with a Tree of Work Items CSV. Maestro will import the hierarchy, summarize data quality,
          and ask you to confirm sprint cadence before entering the app.
        </p>
      </section>

      <DropZone disabled={importing} error={importError} on:file={importFile} on:invalid={handleInvalidFile} />
    </div>

    {#if importing}
      <div class="state-card" aria-live="polite">Importing {selectedFileName}…</div>
    {/if}
  {:else if step === 'sprints' && report}
    <div class="stack">
      <section class="preview-card" aria-label="Import preview">
        <div class="preview-title">
          <p class="eyebrow">Import preview</p>
          <h2>{selectedFileName || 'CSV import'} loaded</h2>
        </div>
        <div class="count-grid">
          {#each previewCounts as count}
            <div class="count-card">
              <span>{count.value}</span>
              <small>{count.label}</small>
            </div>
          {/each}
        </div>

        {#if report.warnings.length > 0 || report.skipped_rows > 0 || report.ambiguous_dates.length > 0}
          <div class="notice-list" aria-label="Import warnings">
            {#if report.skipped_rows > 0}
              <p>{report.skipped_rows} malformed row{report.skipped_rows === 1 ? '' : 's'} skipped during import.</p>
            {/if}
            {#if report.ambiguous_dates.length > 0}
              <p>{report.ambiguous_dates.length} ambiguous date{report.ambiguous_dates.length === 1 ? '' : 's'} detected.</p>
            {/if}
            {#each report.warnings.slice(0, 4) as warning}
              <p>{warning}</p>
            {/each}
          </div>
        {/if}
      </section>

      <SprintPreview
        detectedSprints={report.sprints_detected}
        bind:drafts={sprintDrafts}
        bind:manualMode={manualSprintMode}
        bind:manualConfig={manualSprintConfig}
      />

      <div class="actions">
        <button class="secondary" type="button" on:click={backToImport}>Back</button>
        <button class="primary" type="button" on:click={continueFromSprints}>Continue</button>
      </div>
    </div>
  {:else if step === 'date-assignment' && report}
    <div class="stack">
      <div class="state-card">
        <p class="eyebrow">Date assignment</p>
        <p>
          {report.missing_dates_count} imported item{report.missing_dates_count === 1 ? '' : 's'} need target dates before confirmation.
        </p>
      </div>
      <DateAssignment {report} on:done={continueToConfirm} on:skip={skipDateAssignment} />
    </div>
  {:else if step === 'confirm' && report}
    <section class="confirm-card">
      <p class="eyebrow">Ready</p>
      <h2>Confirm and build roadmap</h2>
      <p class="summary">
        Maestro will enter List view with the imported hierarchy, reviewed sprint cadence, and latest import report stored in project state.
      </p>

      <div class="summary-grid">
        <div><span>{report.epic_count}</span><small>Epics</small></div>
        <div><span>{report.feature_count}</span><small>Features</small></div>
        <div><span>{report.story_count}</span><small>Stories</small></div>
        <div><span>{sprintDrafts.length}</span><small>Sprints</small></div>
        <div><span>{assignedCount}</span><small>Dates assigned</small></div>
        <div><span>{skippedDateCount || report.missing_dates_count - assignedCount}</span><small>Dates deferred</small></div>
      </div>

      <div class="confirm-detail">
        <div>
          <span>Source file</span>
          <strong>{selectedFileName || 'CSV import'}</strong>
        </div>
        <div>
          <span>Sprint cadence</span>
          <strong>{manualSprintMode ? `${manualSprintConfig.count} generated sprints` : `${report.sprints_detected.length} imported sprints`}</strong>
        </div>
        <div>
          <span>Data quality</span>
          <strong>{report.orphaned_features + report.orphaned_stories} orphans · {report.skipped_rows} skipped rows</strong>
        </div>
      </div>

      <div class="actions">
        <button class="secondary" type="button" on:click={() => (step = 'sprints')}>Back</button>
        <button class="primary" type="button" on:click={buildRoadmap}>Build roadmap</button>
      </div>
    </section>
  {/if}
</section>

<style>
  .onboarding {
    width: min(1180px, 100%);
    margin: 0 auto;
    padding: 32px;
  }

  .intro-grid {
    display: grid;
    grid-template-columns: minmax(280px, 0.62fr) minmax(340px, 1fr);
    gap: 24px;
    align-items: stretch;
  }

  .copy-card,
  .preview-card,
  .confirm-card,
  .state-card {
    border: 1px solid var(--border);
    border-radius: 10px;
    background: var(--surface-glass);
    padding: 24px;
  }

  .copy-card {
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  .eyebrow {
    margin: 0 0 10px;
    color: var(--text3);
    font-size: 13px;
    letter-spacing: 0.12em;
    text-transform: uppercase;
  }

  h1,
  h2 {
    margin: 0;
    color: var(--text);
    font-family: var(--font-display);
    font-weight: 300;
    letter-spacing: -0.05em;
    line-height: 1.05;
  }

  h1 {
    font-size: clamp(34px, 5vw, 58px);
  }

  h2 {
    font-size: 30px;
  }

  .copy-card p:not(.eyebrow),
  .summary,
  .state-card p:not(.eyebrow) {
    margin: 16px 0 0;
    color: var(--text2);
    font-size: 15px;
  }

  .stack {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .preview-card {
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    gap: 20px;
    align-items: start;
  }

  .preview-title h2 {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .count-grid,
  .summary-grid {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 12px;
  }

  .count-card,
  .summary-grid div {
    padding: 15px;
    border: 1px solid var(--border);
    border-radius: 8px;
    background: var(--surface-panel);
  }

  .count-card span,
  .summary-grid span {
    display: block;
    color: var(--text);
    font-family: var(--font-display);
    font-size: 30px;
    font-weight: 300;
  }

  .count-card small,
  .summary-grid small {
    display: block;
    margin-top: 4px;
    color: var(--text3);
    font-size: 12px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .notice-list {
    grid-column: 1 / -1;
    display: grid;
    gap: 8px;
    padding: 14px 16px;
    border: 1px solid color-mix(in srgb, var(--amber) 32%, transparent);
    border-radius: 8px;
    background: color-mix(in srgb, var(--amber) 6%, transparent);
  }

  .notice-list p {
    margin: 0;
    color: var(--text2);
    font-size: 13px;
  }

  .summary-grid {
    margin: 24px 0;
  }

  .confirm-detail {
    display: grid;
    gap: 10px;
    margin-bottom: 24px;
    padding: 16px;
    border: 1px solid var(--border);
    border-radius: 8px;
    background: var(--bg2);
  }

  .confirm-detail div {
    display: flex;
    justify-content: space-between;
    gap: 18px;
    border-bottom: 1px solid var(--border);
    padding-bottom: 10px;
  }

  .confirm-detail div:last-child {
    border-bottom: 0;
    padding-bottom: 0;
  }

  .confirm-detail span {
    color: var(--text3);
    font-size: 12px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .confirm-detail strong {
    color: var(--text2);
    font-size: 13px;
    font-weight: 400;
    text-align: right;
  }

  .actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }

  button {
    padding: 10px 16px;
    border-radius: 6px;
    background: transparent;
    cursor: pointer;
    font-family: var(--font-mono);
    font-size: 13px;
    letter-spacing: 0.04em;
  }

  button.primary {
    border: 1px solid var(--accent2);
    color: var(--accent);
  }

  button.secondary {
    border: 1px solid var(--border2);
    color: var(--text2);
  }

  button:hover {
    background: var(--accent-a8);
  }

  @media (max-width: 980px) {
    .intro-grid,
    .preview-card {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 720px) {
    .count-grid,
    .summary-grid {
      grid-template-columns: 1fr 1fr;
    }

    .confirm-detail div,
    .actions {
      flex-direction: column;
      align-items: stretch;
    }

    .confirm-detail strong {
      text-align: left;
    }
  }
</style>
