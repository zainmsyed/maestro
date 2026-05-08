<script lang="ts">
  import type { ManualSprintConfig, SprintDraft } from '../lib/api';

  export let detectedSprints: string[] = [];
  export let drafts: SprintDraft[] = [];
  export let manualMode = false;
  export let manualConfig: ManualSprintConfig = {
    sprint_length_days: 14,
    first_start_date: new Date().toISOString().slice(0, 10),
    count: 6,
  };

  function importedDrafts(): SprintDraft[] {
    return detectedSprints.map((name) => ({
      name,
      start_date: '',
      length_days: manualConfig.sprint_length_days,
      source: 'imported',
    }));
  }

  function addDays(date: string, days: number): string {
    if (!date) return '';
    const parsed = new Date(`${date}T00:00:00`);
    if (Number.isNaN(parsed.getTime())) return '';
    parsed.setDate(parsed.getDate() + days);
    return parsed.toISOString().slice(0, 10);
  }

  function generatedDrafts(): SprintDraft[] {
    const count = Math.max(1, Math.min(52, Number(manualConfig.count) || 1));
    const length = Math.max(1, Math.min(60, Number(manualConfig.sprint_length_days) || 14));
    return Array.from({ length: count }, (_, index) => ({
      name: `Sprint ${index + 1}`,
      start_date: addDays(manualConfig.first_start_date, index * length),
      length_days: length,
      source: 'generated' as const,
    }));
  }

  function setManualMode(event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    manualMode = input.checked;
    drafts = manualMode ? generatedDrafts() : importedDrafts();
  }

  function updateConfig(field: keyof ManualSprintConfig, event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    const nextValue = field === 'first_start_date' ? input.value : Number(input.value);
    manualConfig = { ...manualConfig, [field]: nextValue };
    if (manualMode) drafts = generatedDrafts();
  }

  function updateDraftStart(index: number, event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    drafts = drafts.map((draft, draftIndex) =>
      draftIndex === index ? { ...draft, start_date: input.value } : draft,
    );
  }

  function updateDraftLength(index: number, event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    const length = Math.max(1, Math.min(60, Number(input.value) || 1));
    drafts = drafts.map((draft, draftIndex) =>
      draftIndex === index ? { ...draft, length_days: length } : draft,
    );
  }

  $: if (!manualMode && detectedSprints.length === 0) {
    manualMode = true;
  }

  $: if (manualMode && drafts.length === 0) {
    drafts = generatedDrafts();
  }
</script>

<section class="sprint-preview" aria-label="Sprint review">
  <div class="preview-header">
    <div>
      <p class="eyebrow">Sprint review</p>
      <h2>Confirm your roadmap cadence</h2>
      <p class="intro">
        Review imported sprint labels or generate a clean sequence when the CSV does not carry a usable cadence.
      </p>
    </div>

    <label class="manual-toggle">
      <input type="checkbox" checked={manualMode} on:change={setManualMode} />
      <span>Use manual sprint config</span>
    </label>
  </div>

  {#if manualMode}
    <div class="manual-config" aria-label="Manual sprint configuration">
      <label>
        <span>Sprint length</span>
        <input
          type="number"
          min="1"
          max="60"
          value={manualConfig.sprint_length_days}
          on:input={(event) => updateConfig('sprint_length_days', event)}
        />
      </label>
      <label>
        <span>First start date</span>
        <input
          type="date"
          value={manualConfig.first_start_date}
          on:input={(event) => updateConfig('first_start_date', event)}
        />
      </label>
      <label>
        <span>Sprint count</span>
        <input
          type="number"
          min="1"
          max="52"
          value={manualConfig.count}
          on:input={(event) => updateConfig('count', event)}
        />
      </label>
    </div>
  {:else if detectedSprints.length === 0}
    <div class="empty-note">No sprints were detected. Enable manual config to generate a cadence.</div>
  {/if}

  <div class="table-wrap">
    <table>
      <caption class="sr-only">Sprint cadence review: editable start dates and lengths</caption>
      <thead>
        <tr>
          <th scope="col">Sprint</th>
          <th scope="col">Source</th>
          <th scope="col">Start date</th>
          <th scope="col">Length</th>
        </tr>
      </thead>
      <tbody>
        {#each drafts as draft, index}
          <tr>
            <td>
              <span class="sprint-name">{draft.name}</span>
            </td>
            <td><span class="source-pill">{draft.source}</span></td>
            <td>
              <input
                type="date"
                aria-label={`Start date for ${draft.name}`}
                value={draft.start_date}
                on:input={(event) => updateDraftStart(index, event)}
              />
            </td>
            <td>
              <input
                class="length-input"
                type="number"
                min="1"
                max="60"
                aria-label={`Length in days for ${draft.name}`}
                value={draft.length_days}
                on:input={(event) => updateDraftLength(index, event)}
              />
              <span class="days-label">days</span>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</section>

<style>
  .sprint-preview {
    border: 1px solid var(--border);
    border-radius: 10px;
    background: var(--surface-glass);
    padding: 24px;
  }

  .preview-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 24px;
    margin-bottom: 20px;
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
    max-width: 650px;
    margin: 10px 0 0;
    color: var(--text2);
    font-size: 15px;
  }

  .manual-toggle {
    display: inline-flex;
    align-items: center;
    gap: 10px;
    padding: 9px 12px;
    border: 1px solid var(--border);
    border-radius: 6px;
    background: var(--bg2);
    color: var(--text2);
    cursor: pointer;
    font-size: 13px;
    white-space: nowrap;
  }

  .manual-toggle input {
    accent-color: var(--accent);
  }

  .manual-config {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 14px;
    margin-bottom: 18px;
    padding: 16px;
    border: 1px solid color-mix(in srgb, var(--accent2) 32%, transparent);
    border-radius: 8px;
    background: color-mix(in srgb, var(--accent) 4%, transparent);
  }

  label {
    display: flex;
    flex-direction: column;
    gap: 7px;
    color: var(--text3);
    font-size: 12px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  input[type='date'],
  input[type='number'] {
    height: 36px;
    padding: 0 10px;
    border: 1px solid var(--border);
    border-radius: 6px;
    background: var(--bg2);
    color: var(--text);
    font-family: var(--font-mono);
    font-size: 14px;
  }

  input:focus {
    outline: none;
    border-color: var(--accent2);
  }

  .empty-note {
    margin-bottom: 18px;
    padding: 14px 16px;
    border: 1px solid var(--border);
    border-radius: 8px;
    background: var(--bg2);
    color: var(--text2);
    font-size: 13px;
  }

  .table-wrap {
    overflow: auto;
    border: 1px solid var(--border);
    border-radius: 8px;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    min-width: 680px;
  }

  th,
  td {
    padding: 13px 14px;
    border-bottom: 1px solid var(--border);
    text-align: left;
  }

  th {
    background: var(--bg2);
    color: var(--text3);
    font-size: 12px;
    font-weight: 400;
    letter-spacing: 0.1em;
    text-transform: uppercase;
  }

  tr:last-child td {
    border-bottom: 0;
  }

  .sprint-name {
    color: var(--text);
  }

  .source-pill {
    display: inline-flex;
    padding: 3px 8px;
    border: 1px solid var(--border);
    border-radius: 999px;
    color: var(--text2);
    font-size: 12px;
    text-transform: capitalize;
  }

  .length-input {
    width: 78px;
  }

  .days-label {
    margin-left: 8px;
    color: var(--text3);
    font-size: 12px;
  }

  .sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    white-space: nowrap;
    border: 0;
  }

  @media (max-width: 840px) {
    .preview-header,
    .manual-toggle {
      align-items: flex-start;
    }

    .preview-header,
    .manual-config {
      grid-template-columns: 1fr;
      flex-direction: column;
    }
  }
</style>
