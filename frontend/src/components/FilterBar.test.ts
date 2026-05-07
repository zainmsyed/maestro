import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import FilterBar from './FilterBar.svelte';

describe('FilterBar', () => {
  it('emits the newly selected filter value and renders active chips', async () => {
    const onFilterChange = vi.fn();

    const { rerender } = render(FilterBar, {
      filters: { epic: '', owner: '', sprint: '', status: '' },
      groupBy: 'epic',
      epicOptions: ['Alpha Epic'],
      ownerOptions: ['Alice', 'Bob'],
      sprintOptions: ['Sprint 1'],
      statusOptions: ['Active'],
      onFilterChange,
      onExport: vi.fn(),
    });

    await fireEvent.change(screen.getByLabelText(/epic/i), { target: { value: 'Alpha Epic' } });

    expect(onFilterChange).toHaveBeenCalledWith({
      epic: 'Alpha Epic',
      owner: '',
      sprint: '',
      status: '',
    });

    await rerender({
      filters: { epic: 'Alpha Epic', owner: '', sprint: '', status: '' },
      groupBy: 'epic',
      epicOptions: ['Alpha Epic'],
      ownerOptions: ['Alice', 'Bob'],
      sprintOptions: ['Sprint 1'],
      statusOptions: ['Active'],
      onFilterChange,
      onExport: vi.fn(),
    });

    expect(screen.getByRole('button', { name: /epic:\s*alpha epic/i })).toBeInTheDocument();
  });
});
