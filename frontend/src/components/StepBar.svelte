<script lang="ts">
  type OnboardingStep = 'import' | 'sprints' | 'confirm';

  export let active: OnboardingStep = 'import';

  const steps: Array<{ id: OnboardingStep; label: string; description: string }> = [
    { id: 'import', label: 'Import', description: 'Load CSV' },
    { id: 'sprints', label: 'Sprints', description: 'Review cadence' },
    { id: 'confirm', label: 'Confirm', description: 'Build roadmap' },
  ];

  $: activeIndex = steps.findIndex((step) => step.id === active);
</script>

<nav class="step-bar" aria-label="Onboarding progress">
  {#each steps as step, index}
    <div
      class="step"
      class:active={step.id === active}
      class:complete={index < activeIndex}
      aria-current={step.id === active ? 'step' : undefined}
    >
      <span class="step-index">{index + 1}</span>
      <span class="step-copy">
        <span class="step-label">{step.label}</span>
        <span class="step-description">{step.description}</span>
      </span>
    </div>
    {#if index < steps.length - 1}
      <span class="step-line" class:complete={index < activeIndex}></span>
    {/if}
  {/each}
</nav>

<style>
  .step-bar {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 24px;
  }

  .step {
    display: inline-flex;
    align-items: center;
    gap: 10px;
    min-width: 0;
    color: var(--text3);
  }

  .step-index {
    display: grid;
    place-items: center;
    width: 30px;
    height: 30px;
    border: 1px solid var(--border);
    border-radius: 999px;
    background: var(--bg2);
    color: var(--text3);
    font-size: 13px;
  }

  .step.active .step-index,
  .step.complete .step-index {
    border-color: var(--accent2);
    color: var(--accent);
  }

  .step.active .step-index {
    background: var(--accent-a8);
  }

  .step-label,
  .step-description {
    display: block;
  }

  .step-label {
    color: var(--text2);
    font-size: 13px;
    letter-spacing: 0.1em;
    text-transform: uppercase;
  }

  .step.active .step-label,
  .step.complete .step-label {
    color: var(--text);
  }

  .step-description {
    margin-top: 2px;
    color: var(--text3);
    font-size: 12px;
  }

  .step-line {
    height: 1px;
    flex: 1;
    min-width: 30px;
    background: var(--border);
  }

  .step-line.complete {
    background: color-mix(in srgb, var(--accent2) 70%, transparent);
  }

  @media (max-width: 720px) {
    .step-bar {
      align-items: flex-start;
      flex-direction: column;
    }

    .step-line {
      display: none;
    }
  }
</style>
