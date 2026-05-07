<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let disabled = false;
  export let error: string | null = null;

  const dispatch = createEventDispatcher<{
    file: File;
    invalid: { message: string };
  }>();

  let dragActive = false;
  let fileInput: HTMLInputElement;

  function isCsvFile(file: File): boolean {
    const name = file.name.toLowerCase();
    return name.endsWith('.csv') || file.type === 'text/csv' || file.type === 'application/vnd.ms-excel';
  }

  function chooseFile(file: File | null | undefined) {
    if (!file || disabled) return;
    if (!isCsvFile(file)) {
      dispatch('invalid', { message: 'Unsupported format. Import an Azure DevOps CSV export.' });
      return;
    }
    dispatch('file', file);
  }

  function openPicker() {
    if (!disabled) fileInput.click();
  }

  function handleInput(event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    chooseFile(input.files?.[0]);
    input.value = '';
  }

  function handleDragOver(event: DragEvent) {
    event.preventDefault();
    if (!disabled) dragActive = true;
  }

  function handleDragLeave(event: DragEvent) {
    event.preventDefault();
    dragActive = false;
  }

  function handleDrop(event: DragEvent) {
    event.preventDefault();
    dragActive = false;
    chooseFile(event.dataTransfer?.files?.[0]);
  }
</script>

<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
<div
  class="drop-zone"
  class:active={dragActive}
  class:error={!!error}
  class:disabled
  aria-disabled={disabled}
  aria-label="Upload Azure DevOps CSV"
  on:click={openPicker}
  on:dragover={handleDragOver}
  on:dragleave={handleDragLeave}
  on:drop={handleDrop}
>
  <input bind:this={fileInput} type="file" accept=".csv,text/csv" on:change={handleInput} />

  <div class="drop-icon" aria-hidden="true">CSV</div>
  <div>
    <p class="drop-title">Drop your Tree of Work Items CSV</p>
    <p class="drop-copy">or click to choose a file from your machine</p>
  </div>
  <button class="picker-button" type="button" {disabled} on:click|stopPropagation={openPicker}>Choose file</button>

  {#if error}
    <p class="drop-error" aria-live="polite">{error}</p>
  {/if}
</div>

<style>
  .drop-zone {
    position: relative;
    display: grid;
    place-items: center;
    gap: 18px;
    min-height: 310px;
    padding: 36px;
    border: 1px dashed rgba(139, 111, 62, 0.62);
    border-radius: 10px;
    background:
      radial-gradient(circle at center, rgba(200, 169, 110, 0.07), transparent 55%),
      linear-gradient(145deg, rgba(29, 29, 32, 0.94), rgba(22, 22, 24, 0.94));
    color: var(--text2);
    cursor: pointer;
    text-align: center;
    transition: border-color 0.15s ease, background 0.15s ease, transform 0.15s ease;
  }

  .drop-zone.active {
    border-color: var(--accent);
    background:
      radial-gradient(circle at center, rgba(200, 169, 110, 0.13), transparent 56%),
      linear-gradient(145deg, rgba(29, 29, 32, 0.96), rgba(22, 22, 24, 0.96));
    transform: translateY(-1px);
  }

  .drop-zone.error {
    border-color: rgba(196, 97, 74, 0.7);
  }

  .drop-zone.disabled {
    cursor: progress;
    opacity: 0.72;
  }

  input[type='file'] {
    display: none;
  }

  .drop-icon {
    display: grid;
    place-items: center;
    width: 72px;
    height: 72px;
    border: 1px solid var(--border2);
    border-radius: 14px;
    background: var(--bg2);
    color: var(--accent);
    font-size: 13px;
    letter-spacing: 0.14em;
  }

  .drop-title {
    margin: 0;
    color: var(--text);
    font-family: var(--font-display);
    font-size: clamp(28px, 4vw, 42px);
    font-weight: 300;
    letter-spacing: -0.05em;
  }

  .drop-copy {
    max-width: 480px;
    margin: 12px auto 0;
    color: var(--text2);
    font-size: 15px;
  }

  .picker-button {
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

  .picker-button:hover:not(:disabled) {
    background: rgba(200, 169, 110, 0.08);
  }

  .picker-button:disabled {
    cursor: progress;
    opacity: 0.65;
  }

  .drop-error {
    margin: 0;
    color: var(--coral);
    font-size: 13px;
  }
</style>
