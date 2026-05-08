<script lang="ts">
  import type { GroupBy } from '../lib/listView';

  export let value: GroupBy = 'epic';

  const options: Array<{ value: GroupBy; label: string }> = [
    { value: 'epic', label: 'Epic' },
    { value: 'sprint', label: 'Sprint' },
    { value: 'owner', label: 'Owner' },
  ];
</script>

<label class="group-select">
  <span>Group by</span>
  <select bind:value>
    {#each options as option}
      <option value={option.value}>{option.label}</option>
    {/each}
  </select>
</label>

<style>
  .group-select {
    display: flex;
    flex-direction: column;
    gap: 6px;
    color: var(--text3);
    font-size: 12px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  select {
    height: 36px;
    min-width: 140px;
    padding: 0 10px;
    border: 1px solid var(--border);
    border-radius: 6px;
    background: var(--bg2);
    color: var(--text);
    font-family: var(--font-mono);
    font-size: 14px;
  }

  select:focus {
    outline: none;
    border-color: var(--accent2);
  }
</style>
