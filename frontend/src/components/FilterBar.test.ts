import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import FilterBar from './FilterBar.svelte';

describe('FilterBar', () => {
  it('emits the newly selected filter value and renders active chips', async () => {
    const onFilterChange = vi.fn();

    const { rerender } = render(FilterBar, {
      filters: { epic: '', owner: '', sprint: '', status: '', type: '' },
      groupBy: 'epic',
      epicOptions: ['Alpha Epic'],
      ownerOptions: ['Alice', 'Bob'],
      sprintOptions: ['Sprint 1'],
      statusOptions: ['Active'],
      typeOptions: ['epic', 'feature', 'story'],
      onFilterChange,
      onExport: vi.fn(),
    });

    await fireEvent.change(screen.getByLabelText(/epic/i), { target: { value: 'Alpha Epic' } });

    expect(onFilterChange).toHaveBeenCalledWith({
      epic: 'Alpha Epic',
      owner: '',
      sprint: '',
      status: '',
      type: '',
    });

    await rerender({
      filters: { epic: 'Alpha Epic', owner: '', sprint: '', status: '', type: '' },
      groupBy: 'epic',
      epicOptions: ['Alpha Epic'],
      ownerOptions: ['Alice', 'Bob'],
      sprintOptions: ['Sprint 1'],
      statusOptions: ['Active'],
      typeOptions: ['epic', 'feature', 'story'],
      onFilterChange,
      onExport: vi.fn(),
    });

    expect(screen.getByRole('button', { name: /epic:\s*alpha epic/i })).toBeInTheDocument();
  });

  it('emits type filter changes and renders the type chip', async () => {
    const onFilterChange = vi.fn();

    const { rerender } = render(FilterBar, {
      filters: { epic: '', owner: '', sprint: '', status: '', type: '' },
      groupBy: 'epic',
      epicOptions: [],
      ownerOptions: [],
      sprintOptions: [],
      statusOptions: [],
      typeOptions: ['epic', 'feature', 'story'],
      onFilterChange,
      onExport: vi.fn(),
    });

    await fireEvent.change(screen.getByLabelText(/type/i), { target: { value: 'story' } });

    expect(onFilterChange).toHaveBeenCalledWith({
      epic: '',
      owner: '',
      sprint: '',
      status: '',
      type: 'story',
    });

    await rerender({
      filters: { epic: '', owner: '', sprint: '', status: '', type: 'story' },
      groupBy: 'epic',
      epicOptions: [],
      ownerOptions: [],
      sprintOptions: [],
      statusOptions: [],
      typeOptions: ['epic', 'feature', 'story'],
      onFilterChange,
      onExport: vi.fn(),
    });

    expect(screen.getByRole('button', { name: /type:\s*story/i })).toBeInTheDocument();

    // Click the chip to clear
    await fireEvent.click(screen.getByRole('button', { name: /type:\s*story/i }));

    expect(onFilterChange).toHaveBeenCalledWith({
      epic: '',
      owner: '',
      sprint: '',
      status: '',
      type: '',
    });
  });
});
